package routes

import (
	"net/http"

	"salary_calculator/internal/app"
	"salary_calculator/internal/generated/dbstore"
	"salary_calculator/internal/http/handlers/report/get_salary_report"
	"salary_calculator/internal/http/handlers/salary_change/add_salary_change"
	"salary_calculator/internal/http/handlers/salary_change/delete_salary_change"
	"salary_calculator/internal/http/handlers/salary_change/edit_salary_change"
	"salary_calculator/internal/http/handlers/salary_change/get_salary_changes"
	"salary_calculator/internal/pkg/cache/file"
	"salary_calculator/internal/pkg/http/work_calendar"
	"salary_calculator/internal/services/calculator"
	"salary_calculator/internal/services/work_days"
	addSalaryChangeUC "salary_calculator/internal/usecase/add_salary_change"
	deleteSalaryChangeUC "salary_calculator/internal/usecase/delete_salary_change"
	editSalaryChangeUC "salary_calculator/internal/usecase/edit_salary_change"
	getSalaryChangesUC "salary_calculator/internal/usecase/get_salary_changes"
	getSalaryReportUC "salary_calculator/internal/usecase/get_salary_report"

	"github.com/go-chi/chi/v5"
)

type SalaryRoutesRegistrar struct {
	app *app.App
}

func NewSalaryRoutesRegistrar(a *app.App) *SalaryRoutesRegistrar {
	return &SalaryRoutesRegistrar{app: a}
}

func (s *SalaryRoutesRegistrar) Register(router chi.Router) {
	repo := dbstore.New(s.app.DB)
	httpClient := &http.Client{}
	cache := file.New[string, work_calendar.WorkdayResponse](s.app.Config.Cache.Dir, s.app.Config.Cache.TTL)

	workDaysClient := work_calendar.New(httpClient, cache, s.app.Config.WorkCalendarApiToken)
	workDaysCalc := work_days.New()
	salaryCalc := calculator.New()

	router.Get("/report", get_salary_report.New(
		getSalaryReportUC.New(repo, workDaysClient, workDaysCalc, salaryCalc),
	).ServeHTTP)

	router.Route("/changes", func(r chi.Router) {
		r.Get("/", get_salary_changes.New(getSalaryChangesUC.New(repo)).ServeHTTP)
		r.Post("/", add_salary_change.NewHandler(addSalaryChangeUC.New(repo)).ServeHTTP)
		r.Put("/", edit_salary_change.NewHandler(editSalaryChangeUC.New(repo)).ServeHTTP)
		r.Delete("/", delete_salary_change.NewHandler(deleteSalaryChangeUC.New(repo)).ServeHTTP)
	})
}

func (s *SalaryRoutesRegistrar) Name() string {
	return "salary"
}
