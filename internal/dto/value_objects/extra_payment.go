package value_objects

import (
	"salary_calculator/internal/pkg/utils"
)

type ExtraPaymentType string

const (
	Advance ExtraPaymentType = "advance"
	Salary  ExtraPaymentType = "salary"
)

type ExtraPayment struct {
	Name  string           `json:"name"`
	Value float64          `json:"value"`
	T     ExtraPaymentType `json:"type"`
}

type ExtraPaymentsCollection struct {
	payments []ExtraPayment
	total    float64
	t        ExtraPaymentType
}

func NewExtraPaymentsCollection(t ExtraPaymentType, payments ...ExtraPayment) *ExtraPaymentsCollection {
	total := 0.0

	for _, payment := range payments {
		total += payment.Value
	}

	total = utils.ToTwoDecimals(total)

	return &ExtraPaymentsCollection{
		payments: payments,
		total:    total,
		t:        t,
	}
}

func (e *ExtraPaymentsCollection) Push(payment ExtraPayment) {
	e.payments = append(e.payments, payment)
	e.total += payment.Value
}

func (e *ExtraPaymentsCollection) ToDto() ExtraPaymentsCollectionDto {
	return ExtraPaymentsCollectionDto{
		Payments: e.payments,
		Total:    e.total,
	}
}

func (e *ExtraPaymentsCollection) Type() ExtraPaymentType {
	return e.t
}

func (e *ExtraPaymentsCollection) Total() float64 {
	return e.total
}

type ExtraPaymentsCollectionDto struct {
	Payments []ExtraPayment `json:"payments"`
	Total    float64        `json:"total"`
}
