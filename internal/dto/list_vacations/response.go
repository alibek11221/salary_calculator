package list_vacations

import "salary_calculator/internal/pkg/types"

type VacationItem struct {
	ID   string        `json:"id"`
	From types.ISODate `json:"from"`
	To   types.ISODate `json:"to"`
}

type Out struct {
	Vacations []VacationItem `json:"vacations"`
}
