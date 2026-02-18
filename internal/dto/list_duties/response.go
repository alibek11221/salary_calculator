package list_duties

type DutyItem struct {
	Date       string `json:"date"`
	InWorkdays int32  `json:"in_workdays"`
	InHolidays int32  `json:"in_holidays"`
}

type Out struct {
	Duties []DutyItem `json:"duties"`
}
