package routes

import (
	"salary_calculator/internal/app"
	"salary_calculator/internal/http/handlers/add_s_change"
	"salary_calculator/internal/http/handlers/delete_s_change"
	"salary_calculator/internal/http/handlers/edit_s_change"
	"salary_calculator/internal/http/handlers/get_salary_report"
	"salary_calculator/internal/http/handlers/list_s_changes"
	getSalaryReportUC "salary_calculator/internal/usecase/get_salary_report"
	addSalaryChangeUC "salary_calculator/internal/usecase/salary_change/add"
	deleteSalaryChangeUC "salary_calculator/internal/usecase/salary_change/delete"
	editSalaryChangeUC "salary_calculator/internal/usecase/salary_change/edit"
	listSalaryChangesUC "salary_calculator/internal/usecase/salary_change/list"

	"github.com/go-chi/chi/v5"
)

type SalaryRoutesRegistrar struct {
	app    *app.App
	shared *SharedServices
}

func NewSalaryRoutesRegistrar(a *app.App, shared *SharedServices) *SalaryRoutesRegistrar {
	return &SalaryRoutesRegistrar{app: a, shared: shared}
}

func (s *SalaryRoutesRegistrar) Register(router chi.Router) {
	router.Get("/report", get_salary_report.New(
		getSalaryReportUC.New(
			s.app.Repo,
			s.shared.WorkdaysParser,
			s.shared.WorkdaysCalculator,
			s.shared.SalaryCalculator,
			s.shared.VacationPay,
		),
	).ServeHTTP)

	router.Route("/changes", func(r chi.Router) {
		r.Get("/", list_s_changes.New(listSalaryChangesUC.New(s.app.Repo)).ServeHTTP)
		r.Post("/", add_s_change.NewHandler(addSalaryChangeUC.New(s.app.Repo)).ServeHTTP)
		r.Put("/", edit_s_change.NewHandler(editSalaryChangeUC.New(s.app.Repo)).ServeHTTP)
		r.Delete("/", delete_s_change.NewHandler(deleteSalaryChangeUC.New(s.app.Repo)).ServeHTTP)
	})
}

func (s *SalaryRoutesRegistrar) Name() string {
	return "salary"
}
