package routes

import (
	"net/http"
	"time"

	"salary_calculator/internal/app"
	"salary_calculator/internal/http/handlers/add_s_change"
	"salary_calculator/internal/http/handlers/delete_s_change"
	"salary_calculator/internal/http/handlers/edit_s_change"
	"salary_calculator/internal/http/handlers/get_salary_report"
	"salary_calculator/internal/http/handlers/list_s_changes"
	"salary_calculator/internal/pkg/cache/file"
	"salary_calculator/internal/pkg/http/work_calendar"
	"salary_calculator/internal/services/calculator"
	"salary_calculator/internal/services/work_days"
	getSalaryReportUC "salary_calculator/internal/usecase/get_salary_report"
	addSalaryChangeUC "salary_calculator/internal/usecase/salary_change/add"
	deleteSalaryChangeUC "salary_calculator/internal/usecase/salary_change/delete"
	editSalaryChangeUC "salary_calculator/internal/usecase/salary_change/edit"
	listSalaryChangesUC "salary_calculator/internal/usecase/salary_change/list"

	"github.com/go-chi/chi/v5"
)

type SalaryRoutesRegistrar struct {
	app *app.App
}

func NewSalaryRoutesRegistrar(a *app.App) *SalaryRoutesRegistrar {
	return &SalaryRoutesRegistrar{app: a}
}

func (s *SalaryRoutesRegistrar) Register(router chi.Router) {
	httpClient := &http.Client{Timeout: 10 * time.Second}
	cache := file.New[string, work_calendar.WorkdayResponse](s.app.Config.Cache.Dir, s.app.Config.Cache.TTL, s.app.Logger)

	workDaysClient := work_calendar.New(httpClient, cache, s.app.Config.WorkCalendarApiToken, s.app.Logger)
	workDaysCalc := work_days.New()
	salaryCalc := calculator.New(s.app.Repo)

	router.Get("/report", get_salary_report.New(
		getSalaryReportUC.New(s.app.Repo, workDaysClient, workDaysCalc, salaryCalc),
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
