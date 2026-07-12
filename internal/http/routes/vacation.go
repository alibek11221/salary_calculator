package routes

import (
	"salary_calculator/internal/app"
	"salary_calculator/internal/http/handlers/add_vacation"
	"salary_calculator/internal/http/handlers/delete_vacation"
	"salary_calculator/internal/http/handlers/edit_vacation"
	"salary_calculator/internal/http/handlers/estimate_vacation"
	"salary_calculator/internal/http/handlers/list_vacations"
	getSalaryReportUC "salary_calculator/internal/usecase/get_salary_report"
	addVacationUC "salary_calculator/internal/usecase/vacations/add"
	deleteVacationUC "salary_calculator/internal/usecase/vacations/delete"
	editVacationUC "salary_calculator/internal/usecase/vacations/edit"
	estimateVacationUC "salary_calculator/internal/usecase/vacations/estimate"
	listVacationsUC "salary_calculator/internal/usecase/vacations/list"

	"github.com/go-chi/chi/v5"
)

type VacationRoutesRegistrar struct {
	app    *app.App
	shared *SharedServices
}

func NewVacationRoutesRegistrar(a *app.App, shared *SharedServices) *VacationRoutesRegistrar {
	return &VacationRoutesRegistrar{app: a, shared: shared}
}

func (v *VacationRoutesRegistrar) Register(router chi.Router) {
	reportUC := getSalaryReportUC.New(
		v.app.Repo,
		v.shared.WorkdaysParser,
		v.shared.WorkdaysCalculator,
		v.shared.SalaryCalculator,
		v.shared.VacationPay,
	)

	router.Get("/", list_vacations.New(listVacationsUC.New(v.app.Repo)).ServeHTTP)
	router.Post("/", add_vacation.NewHandler(addVacationUC.New(v.app.Repo)).ServeHTTP)
	router.Put("/", edit_vacation.NewHandler(editVacationUC.New(v.app.Repo)).ServeHTTP)
	router.Delete("/", delete_vacation.NewHandler(deleteVacationUC.New(v.app.Repo)).ServeHTTP)
	router.Get("/estimate", estimate_vacation.NewHandler(
		estimateVacationUC.New(v.app.Repo, v.shared.VacationPay, reportUC),
	).ServeHTTP)
}

func (v *VacationRoutesRegistrar) Name() string {
	return "vacations"
}
