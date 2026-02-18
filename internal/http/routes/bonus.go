package routes

import (
	"salary_calculator/internal/app"
	"salary_calculator/internal/http/handlers/add_bonus"
	"salary_calculator/internal/http/handlers/delete_bonus"
	"salary_calculator/internal/http/handlers/edit_bonus"
	"salary_calculator/internal/http/handlers/list_bonuses"
	addBonusUC "salary_calculator/internal/usecase/bonuses/add"
	deleteBonusUC "salary_calculator/internal/usecase/bonuses/delete"
	editBonusUC "salary_calculator/internal/usecase/bonuses/edit"
	getBonusesUC "salary_calculator/internal/usecase/bonuses/list"

	"github.com/go-chi/chi/v5"
)

type BonusRoutesRegistrar struct {
	app *app.App
}

func NewBonusRoutesRegistrar(a *app.App) *BonusRoutesRegistrar {
	return &BonusRoutesRegistrar{app: a}
}

func (b *BonusRoutesRegistrar) Register(router chi.Router) {
	router.Get("/", list_bonuses.New(getBonusesUC.New(b.app.Repo)).ServeHTTP)
	router.Post("/", add_bonus.NewHandler(addBonusUC.New(b.app.Repo)).ServeHTTP)
	router.Put("/", edit_bonus.NewHandler(editBonusUC.New(b.app.Repo)).ServeHTTP)
	router.Delete("/", delete_bonus.NewHandler(deleteBonusUC.New(b.app.Repo)).ServeHTTP)
}

func (b *BonusRoutesRegistrar) Name() string {
	return "bonuses"
}
