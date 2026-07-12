package add_vacation

import "salary_calculator/internal/pkg/types"

type In struct {
	From types.ISODate `json:"from"`
	To   types.ISODate `json:"to"`
}

type Out struct {
	Ok bool `json:"ok"`
}
