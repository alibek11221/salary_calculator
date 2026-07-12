package routes

import (
	"salary_calculator/internal/app"
	"salary_calculator/internal/http/handlers/add_vacation"
	"salary_calculator/internal/http/handlers/delete_vacation"
	"salary_calculator/internal/http/handlers/edit_vacation"
	"salary_calculator/internal/http/handlers/estimate_vacation"
	"salary_calculator/internal/http/handlers/list_vacations"
	"salary_calculator/internal/pkg/http/work_calendar_parser"
	"salary_calculator/internal/services/calculator"
	"salary_calculator/internal/services/vacation_pay"
	"salary_calculator/internal/services/work_days"
	getSalaryReportUC "salary_calculator/internal/usecase/get_salary_report"
	addVacationUC "salary_calculator/internal/usecase/vacations/add"
	deleteVacationUC "salary_calculator/internal/usecase/vacations/delete"
	editVacationUC "salary_calculator/internal/usecase/vacations/edit"
	estimateVacationUC "salary_calculator/internal/usecase/vacations/estimate"
	listVacationsUC "salary_calculator/internal/usecase/vacations/list"

	"github.com/go-chi/chi/v5"
)

type VacationRoutesRegistrar struct {
	app *app.App
}

func NewVacationRoutesRegistrar(a *app.App) *VacationRoutesRegistrar {
	return &VacationRoutesRegistrar{app: a}
}

func (v *VacationRoutesRegistrar) Register(router chi.Router) {
	workDaysClient := work_calendar_parser.New(v.app.Config.WorkdaysConfig.Dir, v.app.Config.WorkdaysConfig.CacheCap, v.app.Logger)
	workDaysCalc := work_days.New()
	salaryCalc := calculator.New(v.app.Repo)
	vacationPaySvc := vacation_pay.New(v.app.Repo, workDaysClient)
	reportUC := getSalaryReportUC.New(v.app.Repo, workDaysClient, workDaysCalc, salaryCalc, vacationPaySvc)

	router.Get("/", list_vacations.New(listVacationsUC.New(v.app.Repo)).ServeHTTP)
	router.Post("/", add_vacation.NewHandler(addVacationUC.New(v.app.Repo)).ServeHTTP)
	router.Put("/", edit_vacation.NewHandler(editVacationUC.New(v.app.Repo)).ServeHTTP)
	router.Delete("/", delete_vacation.NewHandler(deleteVacationUC.New(v.app.Repo)).ServeHTTP)
	router.Get("/estimate", estimate_vacation.NewHandler(
		estimateVacationUC.New(v.app.Repo, vacationPaySvc, reportUC),
	).ServeHTTP)
}

func (v *VacationRoutesRegistrar) Name() string {
	return "vacations"
}
