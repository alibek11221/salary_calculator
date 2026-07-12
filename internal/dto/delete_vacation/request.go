package delete_vacation

type In struct {
	ID string `json:"id"`
}

type Out struct {
	Ok bool `json:"ok"`
}
