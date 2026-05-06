package reimbursement

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/protone/erp/internal/domain/employee"
	"github.com/protone/erp/internal/domain/finance"
	"github.com/protone/erp/internal/domain/reimbursement"
)

type UseCase interface {
	SubmitClaim(ctx context.Context, input SubmitClaimInput) error
	GetMyClaims(ctx context.Context, tenantID, userID uuid.UUID) ([]reimbursement.Reimbursement, error)
	ApproveClaim(ctx context.Context, tenantID, managerUserID, claimID uuid.UUID) error
}

type SubmitClaimInput struct {
	TenantID    uuid.UUID
	UserID      uuid.UUID
	Category    string
	Amount      float64
	Date        time.Time
	Description string
	ReceiptURL  *string
}

type ReimbursementUseCase struct {
	repo         reimbursement.Repository
	employeeRepo employee.Repository
	financeRepo  finance.Repository
}

func NewReimbursementUseCase(r reimbursement.Repository, er employee.Repository, fr finance.Repository) *ReimbursementUseCase {
	return &ReimbursementUseCase{
		repo:         r,
		employeeRepo: er,
		financeRepo:  fr,
	}
}

func (uc *ReimbursementUseCase) SubmitClaim(ctx context.Context, input SubmitClaimInput) error {
	emp, err := uc.employeeRepo.FindByUserID(ctx, input.TenantID, input.UserID)
	if err != nil {
		return err
	}
	if emp == nil {
		return fmt.Errorf("employee not found")
	}

	claim := reimbursement.NewReimbursement(input.TenantID, emp.ID, input.Category, input.Amount, input.Date, input.Description)
	if input.ReceiptURL != nil {
		claim.ReceiptURL = *input.ReceiptURL
	}

	return uc.repo.Save(ctx, claim)
}

func (uc *ReimbursementUseCase) GetMyClaims(ctx context.Context, tenantID, userID uuid.UUID) ([]reimbursement.Reimbursement, error) {
	emp, err := uc.employeeRepo.FindByUserID(ctx, tenantID, userID)
	if err != nil {
		return nil, err
	}
	if emp == nil {
		return nil, fmt.Errorf("employee not found")
	}

	return uc.repo.GetByEmployee(ctx, emp.ID)
}

func (uc *ReimbursementUseCase) ApproveClaim(ctx context.Context, tenantID, managerUserID, claimID uuid.UUID) error {
	claim, err := uc.repo.GetByID(ctx, claimID)
	if err != nil {
		return err
	}
	if claim == nil {
		return fmt.Errorf("claim not found")
	}

	// BUDGET INTEGRATION
	emp, _ := uc.employeeRepo.FindByID(ctx, tenantID, claim.EmployeeID)
	if emp != nil && emp.DepartmentID != uuid.Nil {
		now := time.Now()
		budget, _ := uc.financeRepo.GetBudget(ctx, tenantID, emp.DepartmentID, int(now.Month()), now.Year())
		if budget != nil {
			if err := budget.Deduct(claim.Amount); err != nil {
				return fmt.Errorf("budget insufficient: %v", err)
			}
			_ = uc.financeRepo.UpdateBudget(ctx, budget)
		}
	}

	claim.Approve(managerUserID)
	if err := uc.repo.Update(ctx, claim); err != nil {
		return err
	}

	// FINANCE INTEGRATION: Create Journal Entry
	expenseAcc, _ := uc.financeRepo.GetCOAByCode(ctx, tenantID, "5-1002") // Beban Reimbursement
	cashAcc, _ := uc.financeRepo.GetCOAByCode(ctx, tenantID, "1-1001")    // Kas/Bank

	if expenseAcc == nil || cashAcc == nil {
		return nil // Still success even if accounting fails, but ideally log this
	}

	journal := &finance.JournalEntry{
		ID:          uuid.New(),
		TenantID:    tenantID,
		JournalNo:   fmt.Sprintf("REIM/%s", uuid.NewString()[:8]),
		Date:        time.Now(),
		Description: fmt.Sprintf("Reimbursement: %s - %s", claim.Category, claim.Description),
		Status:      finance.StatusPosted,
		SourceType:  "reimbursement",
		SourceID:    &claim.ID,
	}

	journal.AddLine(expenseAcc.ID, claim.Description, claim.Amount, 0)
	journal.AddLine(cashAcc.ID, claim.Description, 0, claim.Amount)

	return uc.financeRepo.CreateJournal(ctx, journal)
}
