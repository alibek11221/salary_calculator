package value_objects

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
	totals   map[ExtraPaymentType]float64
}

func NewExtraPaymentsCollection(payments ...ExtraPayment) *ExtraPaymentsCollection {
	totals := make(map[ExtraPaymentType]float64, 2)
	totals = map[ExtraPaymentType]float64{
		Advance: 0,
		Salary:  0,
	}

	for _, payment := range payments {
		totals[payment.T] += payment.Value
	}

	return &ExtraPaymentsCollection{
		payments: payments,
		totals:   totals,
	}
}

func (e *ExtraPaymentsCollection) Push(payment ExtraPayment) {
	if e == nil {
		return
	}
	e.payments = append(e.payments, payment)
	if e.totals == nil {
		e.totals = make(map[ExtraPaymentType]float64)
	}
	e.totals[payment.T] += payment.Value
}

func (e *ExtraPaymentsCollection) ToDto() ExtraPaymentsCollectionDto {
	if e == nil {
		return ExtraPaymentsCollectionDto{}
	}
	total := 0.0

	for _, t := range e.totals {
		total += t
	}

	return ExtraPaymentsCollectionDto{
		Payments: e.payments,
		Total:    total,
	}
}

func (e *ExtraPaymentsCollection) Total() map[ExtraPaymentType]float64 {
	if e == nil {
		return nil
	}
	return e.totals
}

type ExtraPaymentsCollectionDto struct {
	Payments []ExtraPayment `json:"payments"`
	Total    float64        `json:"total"`
}
