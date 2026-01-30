package edit_salary_change_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	edit_salary_change_dto "salary_calculator/internal/dto/edit_salary_change"
	"salary_calculator/internal/dto/value_objects"
	"salary_calculator/internal/http/handlers/salary_change/edit_salary_change"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestHandler_ServeHTTP(t *testing.T) {
	type fields struct {
		u *Mockusecase
	}

	sd, _ := value_objects.NewSalaryDate("2025_01")

	tests := []struct {
		name           string
		setup          func(f fields)
		reqBody        interface{}
		wantStatus     int
		wantInResponse string
	}{
		{
			name: "success",
			reqBody: edit_salary_change_dto.In{
				ID:    "550e8400-e29b-41d4-a716-446655440000",
				Value: 120000,
				Date:  *sd,
			},
			setup: func(f fields) {
				f.u.EXPECT().Do(gomock.Any(), gomock.Any()).Return(&edit_salary_change_dto.Out{Ok: true}, nil)
			},
			wantStatus:     http.StatusOK,
			wantInResponse: `"ok":true`,
		},
		{
			name:       "invalid json",
			reqBody:    "invalid",
			setup:      func(f fields) {},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "internal error",
			reqBody: edit_salary_change_dto.In{
				ID:    "550e8400-e29b-41d4-a716-446655440000",
				Value: 120000,
				Date:  *sd,
			},
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

			h := edit_salary_change.NewHandler(f.u)

			var body []byte
			if s, ok := tt.reqBody.(string); ok {
				body = []byte(s)
			} else {
				body, _ = json.Marshal(tt.reqBody)
			}

			req := httptest.NewRequest(http.MethodPut, "/changes", bytes.NewReader(body))
			w := httptest.NewRecorder()

			h.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
			if tt.wantInResponse != "" {
				assert.Contains(t, w.Body.String(), tt.wantInResponse)
			}
		})
	}
}
