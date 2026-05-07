package payroll

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/google/uuid"
	"github.com/protone/erp/internal/domain/employee"
	"github.com/protone/erp/internal/domain/finance"
	"github.com/protone/erp/internal/domain/overtime"
	"github.com/protone/erp/internal/domain/payroll"
	"github.com/protone/erp/pkg/pdf"
)

type PayrollUseCase struct {
	payrollRepo   payroll.Repository
	employeeRepo  employee.Repository
	employeeQuery employee.QueryRepository
	financeRepo   finance.Repository
	overtimeRepo  overtime.Repository
	pdfGen        pdf.Generator
}

func NewPayrollUseCase(pr payroll.Repository, er employee.Repository, eq employee.QueryRepository, fr finance.Repository, or overtime.Repository, pg pdf.Generator) *PayrollUseCase {
	return &PayrollUseCase{
		payrollRepo:   pr,
		employeeRepo:  er,
		employeeQuery: eq,
		financeRepo:   fr,
		overtimeRepo:  or,
		pdfGen:        pg,
	}
}

func (uc *PayrollUseCase) GenerateMonthlyPayroll(ctx context.Context, tenantID uuid.UUID, month, year int) error {
	// 1. Check if period already exists
	period, err := uc.payrollRepo.GetPeriod(ctx, tenantID, month, year)
	if err != nil {
		return err
	}

	if period == nil {
		period = payroll.NewPayrollPeriod(tenantID, month, year)
		if err := uc.payrollRepo.CreatePeriod(ctx, period); err != nil {
			return err
		}
	}

	if period.Status != payroll.StatusDraft {
		return fmt.Errorf("payroll for this period is already %s", period.Status)
	}

	// 1.5 Clear existing slips for re-generation
	if err := uc.payrollRepo.DeleteSlipsByPeriod(ctx, period.ID); err != nil {
		return err
	}

	// 2. Get All Employees
	employees, _, err := uc.employeeQuery.List(ctx, employee.Filter{
		TenantID: tenantID,
		Page:     1,
		PerPage:  1000,
	})
	if err != nil {
		return err
	}

	// 3. Get All Salary Components
	components, err := uc.payrollRepo.GetComponents(ctx, tenantID)
	if err != nil {
		return err
	}

	// 4. Generate Slips
	now := time.Now()
	var totalPeriodAmount float64
	for _, emp := range employees {
		empID, _ := uuid.Parse(emp.ID)
		slip := &payroll.PayrollSlip{
			ID:              uuid.New(),
			TenantID:        tenantID,
			PayrollPeriodID: period.ID,
			EmployeeID:      empID,
			CreatedAt:       now,
			UpdatedAt:       now,
		}

		for _, comp := range components {
			detail := payroll.PayrollSlipDetail{
				SalaryComponentID: comp.ID,
				Type:              comp.Type,
				Amount:            comp.DefaultAmount,
			}
			slip.Details = append(slip.Details, detail)

			if comp.Type == payroll.TypeAllowance {
				if comp.Code == "BASIC" {
					slip.BasicSalary = comp.DefaultAmount
				} else {
					slip.TotalAllowance += comp.DefaultAmount
				}
			} else {
				slip.TotalDeduction += comp.DefaultAmount
			}
		}

		// Get Overtime for this employee in this period
		otHours, _ := uc.overtimeRepo.GetApprovedSumByEmployee(ctx, empID, month, year)
		otAmount := otHours * 50000 // Tarif 50rb per jam

		slip.OvertimeHours = otHours
		slip.OvertimeAmount = otAmount

		slip.CalculateNet()
		totalPeriodAmount += slip.NetSalary

		if err := uc.payrollRepo.SaveSlip(ctx, slip); err != nil {
			return err
		}
	}

	// 5. Update Period Total
	period.TotalAmount = totalPeriodAmount
	return uc.payrollRepo.UpdatePeriod(ctx, period)
}

func (uc *PayrollUseCase) GetMySlip(ctx context.Context, tenantID, userID uuid.UUID, month, year int) (*payroll.PayrollSlip, error) {
	emp, err := uc.employeeRepo.FindByUserID(ctx, tenantID, userID)
	if err != nil {
		return nil, err
	}
	if emp == nil {
		return nil, fmt.Errorf("employee not found")
	}

	period, err := uc.payrollRepo.GetPeriod(ctx, tenantID, month, year)
	if err != nil {
		return nil, err
	}
	if period == nil {
		return nil, fmt.Errorf("payroll period not found")
	}

	return uc.payrollRepo.GetSlipByEmployee(ctx, period.ID, emp.ID)
}

func (uc *PayrollUseCase) ApprovePayroll(ctx context.Context, tenantID uuid.UUID, month, year int) error {
	period, err := uc.payrollRepo.GetPeriod(ctx, tenantID, month, year)
	if err != nil || period == nil {
		return fmt.Errorf("payroll period not found")
	}

	if period.Status != payroll.StatusDraft {
		return fmt.Errorf("payroll is already %s", period.Status)
	}

	period.Status = payroll.StatusApproved
	return uc.payrollRepo.UpdatePeriod(ctx, period)
}

func (uc *PayrollUseCase) PayPayroll(ctx context.Context, tenantID uuid.UUID, month, year int) error {
	period, err := uc.payrollRepo.GetPeriod(ctx, tenantID, month, year)
	if err != nil || period == nil {
		return fmt.Errorf("payroll period not found")
	}

	if period.Status != payroll.StatusApproved {
		return fmt.Errorf("payroll must be approved before payment")
	}

	// 1. Mark as Paid
	period.Status = payroll.StatusPaid
	if err := uc.payrollRepo.UpdatePeriod(ctx, period); err != nil {
		return err
	}

	// 2. BUDGET & FINANCE INTEGRATION
	// Group totals by department first for efficiency
	deptTotals := make(map[uuid.UUID]float64)
	slips, err := uc.payrollRepo.GetSlipsByPeriod(ctx, period.ID)
	if err != nil {
		return fmt.Errorf("failed to fetch slips for payment: %w", err)
	}

	for _, slip := range slips {
		emp, err := uc.employeeRepo.FindByID(ctx, tenantID, slip.EmployeeID)
		if err == nil && emp != nil && emp.DepartmentID != uuid.Nil {
			deptTotals[emp.DepartmentID] += slip.NetSalary
		}
	}

	// Update budgets once per department
	for deptID, amount := range deptTotals {
		budget, err := uc.financeRepo.GetBudget(ctx, tenantID, deptID, month, year)
		if err != nil {
			continue
		}
		if budget != nil {
			if err := budget.Deduct(amount); err != nil {
				return fmt.Errorf("budget limit exceeded for department %s: %v", deptID, err)
			}
			if err := uc.financeRepo.UpdateBudget(ctx, budget); err != nil {
				return fmt.Errorf("failed to update budget for department %s: %w", deptID, err)
			}
		}
	}

	// 3. Create Journal Entry
	expenseAcc, err := uc.financeRepo.GetCOAByCode(ctx, tenantID, "5-1001") // Beban Gaji
	if err != nil || expenseAcc == nil {
		return fmt.Errorf("expense account (5-1001) not configured properly")
	}

	cashAcc, err := uc.financeRepo.GetCOAByCode(ctx, tenantID, "1-1001") // Kas/Bank
	if err != nil || cashAcc == nil {
		return fmt.Errorf("cash/bank account (1-1001) not configured properly")
	}

	journal := &finance.JournalEntry{
		ID:          uuid.New(),
		TenantID:    tenantID,
		JournalNo:   fmt.Sprintf("PYRL/%d/%d/%s", year, month, uuid.NewString()[:8]),
		Date:        time.Now(),
		Description: fmt.Sprintf("Payroll Payment - %02d/%d", month, year),
		Status:      finance.StatusPosted,
		SourceType:  "payroll",
		SourceID:    &period.ID,
	}

	// Debet: Beban Gaji
	journal.AddLine(expenseAcc.ID, fmt.Sprintf("Salaries Payment %02d/%d", month, year), period.TotalAmount, 0)
	// Kredit: Kas/Bank
	journal.AddLine(cashAcc.ID, fmt.Sprintf("Salaries Payment %02d/%d", month, year), 0, period.TotalAmount)

	if err := uc.financeRepo.CreateJournal(ctx, journal); err != nil {
		return fmt.Errorf("failed to create accounting journal: %w", err)
	}

	return nil
}

func (uc *PayrollUseCase) GetSlipPDF(ctx context.Context, tenantID, employeeID, periodID uuid.UUID) (io.Reader, error) {
	// 1. Get Slip Data
	slip, err := uc.payrollRepo.GetSlipByEmployee(ctx, periodID, employeeID)
	if err != nil {
		return nil, err
	}
	if slip == nil {
		return nil, fmt.Errorf("slip not found")
	}

	// 2. Get Period Data
	period, _ := uc.payrollRepo.GetPeriodByID(ctx, periodID)

	// 3. Get Employee Detail (for Dept/Position Name)
	emp, err := uc.employeeQuery.GetByID(ctx, tenantID, employeeID)
	if err != nil {
		return nil, err
	}

	// 4. Prepare PDF Data
	pdfData := pdf.PayrollSlipData{
		CompanyName:    "PROTONE ERP SYSTEM",
		Period:         fmt.Sprintf("%02d/%d", period.Month, period.Year),
		EmployeeName:   emp.FullName,
		EmployeeNIK:    emp.NIK,
		Department:     emp.DepartmentName,
		Position:       emp.PositionName,
		BasicSalary:    slip.BasicSalary,
		OvertimeAmount: slip.OvertimeAmount,
		NetSalary:      slip.NetSalary,
	}

	// 5. Fetch Salary Components to map their real names
	components, _ := uc.payrollRepo.GetComponents(ctx, tenantID)
	compMap := make(map[uuid.UUID]payroll.SalaryComponent)
	for _, c := range components {
		compMap[c.ID] = c
	}

	// 6. Group details by type
	for _, d := range slip.Details {
		comp, exists := compMap[d.SalaryComponentID]
		compName := "Unknown Component"
		if exists {
			compName = comp.Name
		}

		compData := pdf.ComponentData{
			Name:   compName,
			Amount: d.Amount,
		}

		if d.Type == payroll.TypeAllowance {
			// Skip BASIC salary from allowances list to avoid duplication
			if exists && comp.Code == "BASIC" {
				continue
			}
			pdfData.Allowances = append(pdfData.Allowances, compData)
		} else {
			pdfData.Deductions = append(pdfData.Deductions, compData)
		}
	}

	return uc.pdfGen.GeneratePayrollSlip(ctx, pdfData)
}

func (uc *PayrollUseCase) DownloadMySlipPDF(ctx context.Context, tenantID, userID, periodID uuid.UUID) (io.Reader, error) {
	emp, err := uc.employeeRepo.FindByUserID(ctx, tenantID, userID)
	if err != nil {
		return nil, fmt.Errorf("employee not found for user: %w", err)
	}
	return uc.GetSlipPDF(ctx, tenantID, emp.ID, periodID)
}
