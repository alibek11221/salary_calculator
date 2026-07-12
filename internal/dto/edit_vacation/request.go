package edit_vacation

import "salary_calculator/internal/pkg/types"

type In struct {
	ID   string        `json:"id"`
	From types.ISODate `json:"from"`
	To   types.ISODate `json:"to"`
}

type Out struct {
	Ok bool `json:"ok"`
}
