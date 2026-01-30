package get_salary_report_test

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	get_salary_report_dto "salary_calculator/internal/dto/get_salary_report"
	"salary_calculator/internal/http/handlers/report/get_salary_report"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestHandler_ServeHTTP(t *testing.T) {
	type fields struct {
		u *Mockusecase
	}
	tests := []struct {
		name           string
		setup          func(f fields)
		query          string
		wantStatus     int
		wantInResponse string
	}{
		{
			name:  "success",
			query: "?year=2025&month=1",
			setup: func(f fields) {
				f.u.EXPECT().Do(gomock.Any(), get_salary_report_dto.In{Year: 2025, Month: 1}).
					Return(&get_salary_report_dto.Out{BaseSalary: 100000}, nil)
			},
			wantStatus:     http.StatusOK,
			wantInResponse: `"base_salary":100000`,
		},
		{
			name:       "missing parameters",
			query:      "",
			setup:      func(f fields) {},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "invalid year",
			query:      "?year=abc&month=1",
			setup:      func(f fields) {},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "invalid month",
			query:      "?year=2025&month=abc",
			setup:      func(f fields) {},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:  "internal error",
			query: "?year=2025&month=1",
			setup: func(f fields) {
				f.u.EXPECT().Do(gomock.Any(), gomock.Any()).Return(nil, errors.New("uc error"))
			},
			wantStatus:     http.StatusInternalServerError,
			wantInResponse: "uc error",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{
				u: NewMockusecase(ctrl),
			}
			tt.setup(f)

			h := get_salary_report.New(f.u)

			req := httptest.NewRequest(http.MethodGet, "/report"+tt.query, nil)
			w := httptest.NewRecorder()

			h.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code, fmt.Sprintf("test %s failed", tt.name))
			if tt.wantInResponse != "" {
				assert.Contains(t, w.Body.String(), tt.wantInResponse)
			}
		})
	}
}
