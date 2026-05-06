package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/protone/erp/config"
	"github.com/protone/erp/internal/delivery/http/handler"
	httpmiddleware "github.com/protone/erp/internal/delivery/http/middleware"
	"github.com/protone/erp/internal/infrastructure/persistence/postgres"
	attuc "github.com/protone/erp/internal/usecase/attendance"
	"github.com/protone/erp/internal/usecase/auth"
	empcommand "github.com/protone/erp/internal/usecase/employee/command"
	empquery "github.com/protone/erp/internal/usecase/employee/query"
	finuc "github.com/protone/erp/internal/usecase/finance"
	leaveuc "github.com/protone/erp/internal/usecase/leave"
	overuc "github.com/protone/erp/internal/usecase/overtime"
	payrolluc "github.com/protone/erp/internal/usecase/payroll"
	reimuc "github.com/protone/erp/internal/usecase/reimbursement"
	reportinguc "github.com/protone/erp/internal/usecase/reporting"
	pkgjwt "github.com/protone/erp/pkg/jwt"
	"github.com/protone/erp/pkg/pdf"
	"github.com/protone/erp/pkg/response"
)

// New membuat router dengan semua dependency injected di sini.
// Ini adalah composition root — satu-satunya tempat yang tahu semua dependensi.
func New(cfg *config.Config, pool *pgxpool.Pool) http.Handler {
	r := chi.NewRouter()

	// ─── Global Middleware ──────────────────────────────────────────────
	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)
	r.Use(chimiddleware.CleanPath)
	r.Use(corsMiddleware)

	// ─── Dependency Injection ───────────────────────────────────────────

	// JWT
	jwtManager := pkgjwt.NewManager(cfg.JWT)
	authMW := httpmiddleware.NewAuthMiddleware(jwtManager)

	// PDF Generator
	pdfGen := pdf.NewMarotoGenerator()

	// Repositories (adapters)
	empWriteRepo, empReadRepo := postgres.NewEmployeeRepository(pool)
	tenantRepo := postgres.NewTenantRepository(pool)
	userRepo := postgres.NewUserRepository(pool)
	attRepo, attQueryRepo := postgres.NewAttendanceRepository(pool)
	leaveRepo := postgres.NewLeaveRepository(pool)
	payrollRepo := postgres.NewPayrollRepository(pool)
	reimRepo := postgres.NewReimbursementRepository(pool)
	financeRepo := postgres.NewFinanceRepository(pool)
	overRepo := postgres.NewOvertimeRepository(pool)
	txManager := postgres.NewTxManager(pool)

	// Use Cases
	registerTenantUC := auth.NewRegisterTenant(tenantRepo, userRepo, txManager)
	loginUC := auth.NewLogin(userRepo, jwtManager)

	createEmpUC := empcommand.NewCreateEmployee(empWriteRepo)
	updateEmpUC := empcommand.NewUpdateEmployee(empWriteRepo)
	deleteEmpUC := empcommand.NewDeleteEmployee(empWriteRepo)
	listEmpUC := empquery.NewListEmployees(empReadRepo)
	getEmpUC := empquery.NewGetEmployee(empReadRepo)

	attUC := attuc.NewAttendanceUseCase(attRepo, empWriteRepo)
	leaveUC := leaveuc.NewLeaveUseCase(leaveRepo, empWriteRepo)
	payrollUC := payrolluc.NewPayrollUseCase(payrollRepo, empWriteRepo, empReadRepo, financeRepo, overRepo, pdfGen)
	reimUC := reimuc.NewReimbursementUseCase(reimRepo, empWriteRepo, financeRepo)
	overUC := overuc.NewOvertimeUseCase(overRepo, empWriteRepo)
	budgetUC := finuc.NewUseCase(financeRepo)
	reportingUC := reportinguc.NewReportingUseCase(payrollRepo, attQueryRepo)

	// Handlers
	empHandler := handler.NewEmployeeHandler(
		createEmpUC, updateEmpUC, deleteEmpUC,
		listEmpUC, getEmpUC,
	)
	authHandler := handler.NewAuthHandler(registerTenantUC, loginUC)
	attHandler := handler.NewAttendanceHandler(attUC)
	leaveHandler := handler.NewLeaveHandler(leaveUC)
	payrollHandler := handler.NewPayrollHandler(payrollUC)
	reimHandler := handler.NewReimbursementHandler(reimUC)
	overHandler := handler.NewOvertimeHandler(overUC)
	budgetHandler := handler.NewBudgetHandler(budgetUC)
	reportingHandler := handler.NewReportingHandler(reportingUC)

	// ─── Routes ─────────────────────────────────────────────────────────

	// Health check (public)
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		response.Success(w, "ProtoERP API is running", map[string]string{
			"version": "1.0.0",
			"status":  "healthy",
		})
	})

	// API v1
	r.Route("/api/v1", func(r chi.Router) {
		// Auth routes (public)
		r.Post("/auth/register-tenant", authHandler.RegisterTenant)
		r.Post("/auth/login", authHandler.Login)
		// r.Post("/auth/refresh", authHandler.Refresh)

		// Protected routes
		r.Group(func(r chi.Router) {
			r.Use(authMW.Authenticate)

			// Dashboard routes
			r.Get("/dashboard/stats", reportingHandler.GetDashboardStats)

			// Employee routes
			r.Route("/employees", func(r chi.Router) {
				r.Get("/", empHandler.List)
				r.Post("/", empHandler.Create)
				r.Get("/{id}", empHandler.Get)
				r.Put("/{id}", empHandler.Update)
				r.Delete("/{id}", empHandler.Delete)
			})

			// Department list (Quick helper)
			r.Get("/departments", func(w http.ResponseWriter, r *http.Request) {
				query := `SELECT id, name FROM departments`
				rows, _ := pool.Query(r.Context(), query)
				var depts []map[string]string
				for rows.Next() {
					var id, name string
					rows.Scan(&id, &name)
					depts = append(depts, map[string]string{"id": id, "name": name})
				}
				response.Success(w, "Department list fetched", depts)
			})

			// Attendance routes
			r.Route("/attendance", func(r chi.Router) {
				r.Post("/clock-in", attHandler.ClockIn)
				r.Post("/clock-out", attHandler.ClockOut)
			})

			// Leave routes
			r.Route("/leaves", func(r chi.Router) {
				r.Get("/types", leaveHandler.GetTypes)
				r.Post("/request", leaveHandler.RequestLeave)
				r.Get("/my-requests", leaveHandler.GetMyRequests)
				r.Patch("/approve", leaveHandler.Approve)
				r.Post("/init-balance", leaveHandler.GiveBalance)
			})

			// Payroll routes
			r.Route("/payroll", func(r chi.Router) {
				r.Post("/generate", payrollHandler.Generate)
				r.Get("/my-slip", payrollHandler.GetMySlip)
				r.Patch("/approve", payrollHandler.Approve)
				r.Post("/pay", payrollHandler.Pay)
				r.Get("/slips/{id}/download", payrollHandler.DownloadSlip)
			})

			// Reimbursement routes
			r.Route("/reimbursements", func(r chi.Router) {
				r.Post("/submit", reimHandler.Submit)
				r.Get("/my-claims", reimHandler.GetMyClaims)
				r.Patch("/approve", reimHandler.Approve)
			})

			// Overtime routes
			r.Route("/overtime", func(r chi.Router) {
				r.Post("/request", overHandler.Request)
				r.Get("/my-requests", overHandler.GetMyRequests)
				r.Patch("/approve", overHandler.Approve)
			})

			// Budget routes
			r.Route("/budgets", func(r chi.Router) {
				r.Post("/set", budgetHandler.SetBudget)
				r.Get("/check", func(w http.ResponseWriter, r *http.Request) {
					deptID := r.URL.Query().Get("department_id")
					month := r.URL.Query().Get("month")
					year := r.URL.Query().Get("year")
					query := `SELECT allocated_amount, spent_amount FROM department_budgets WHERE department_id = $1 AND month = $2 AND year = $3`
					var alloc, spent float64
					pool.QueryRow(r.Context(), query, deptID, month, year).Scan(&alloc, &spent)
					response.Success(w, "Budget status", map[string]float64{
						"allocated": alloc,
						"spent":     spent,
						"remaining": alloc - spent,
					})
				})
			})
		})
	})

	return r
}

// corsMiddleware menambahkan CORS headers (untuk development)
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
