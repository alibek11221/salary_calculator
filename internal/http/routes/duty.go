package routes

import (
	"salary_calculator/internal/app"
	"salary_calculator/internal/http/handlers/add_duty"
	"salary_calculator/internal/http/handlers/delete_duty"
	"salary_calculator/internal/http/handlers/edit_duty"
	"salary_calculator/internal/http/handlers/list_duties"
	addDutyUC "salary_calculator/internal/usecase/duties/add"
	deleteDutyUC "salary_calculator/internal/usecase/duties/delete"
	editDutyUC "salary_calculator/internal/usecase/duties/edit"
	listDytiesUC "salary_calculator/internal/usecase/duties/list"

	"github.com/go-chi/chi/v5"
)

type DutyRoutesRegistrar struct {
	app *app.App
}

func NewDutyRoutesRegistrar(a *app.App) *DutyRoutesRegistrar {
	return &DutyRoutesRegistrar{app: a}
}

func (d *DutyRoutesRegistrar) Register(router chi.Router) {
	router.Get("/", list_duties.New(listDytiesUC.New(d.app.Repo)).ServeHTTP)
	router.Post("/", add_duty.NewHandler(addDutyUC.New(d.app.Repo)).ServeHTTP)
	router.Put("/", edit_duty.NewHandler(editDutyUC.New(d.app.Repo)).ServeHTTP)
	router.Delete("/", delete_duty.NewHandler(deleteDutyUC.New(d.app.Repo)).ServeHTTP)
}

func (d *DutyRoutesRegistrar) Name() string {
	return "duties"
}
