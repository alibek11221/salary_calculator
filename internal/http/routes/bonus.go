package routes

import (
	"salary_calculator/internal/app"
	"salary_calculator/internal/generated/dbstore"
	"salary_calculator/internal/http/handlers/bonus/add_bonus"
	"salary_calculator/internal/http/handlers/bonus/delete_bonus"
	"salary_calculator/internal/http/handlers/bonus/edit_bonus"
	"salary_calculator/internal/http/handlers/bonus/get_bonuses"
	addBonusUC "salary_calculator/internal/usecase/add_bonus"
	deleteBonusUC "salary_calculator/internal/usecase/delete_bonus"
	editBonusUC "salary_calculator/internal/usecase/edit_bonus"
	getBonusesUC "salary_calculator/internal/usecase/get_bonuses"

	"github.com/go-chi/chi/v5"
)

type BonusRoutesRegistrar struct {
	app *app.App
}

func NewBonusRoutesRegistrar(a *app.App) *BonusRoutesRegistrar {
	return &BonusRoutesRegistrar{app: a}
}

func (b *BonusRoutesRegistrar) Register(router chi.Router) {
	repo := dbstore.New(b.app.DB)

	router.Get("/", get_bonuses.New(getBonusesUC.New(repo)).ServeHTTP)
	router.Post("/", add_bonus.NewHandler(addBonusUC.New(repo)).ServeHTTP)
	router.Put("/", edit_bonus.NewHandler(editBonusUC.New(repo)).ServeHTTP)
	router.Delete("/", delete_bonus.NewHandler(deleteBonusUC.New(repo)).ServeHTTP)
}

func (b *BonusRoutesRegistrar) Name() string {
	return "bonuses"
}
