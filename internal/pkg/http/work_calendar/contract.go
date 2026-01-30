package work_calendar

import "net/http"

//go:generate mockgen -source=contract.go -destination mocks_test.go -package "${GOPACKAGE}_test"

type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type cache interface {
	Put(k string, v WorkdayResponse) error
	Get(k string) (value WorkdayResponse, ok bool)
}
