package add_bonus_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	add_bonus_dto "salary_calculator/internal/dto/add_bonus"
	"salary_calculator/internal/dto/value_objects"
	"salary_calculator/internal/http/handlers/bonus/add_bonus"
	addBonusUC "salary_calculator/internal/usecase/add_bonus"

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
			reqBody: add_bonus_dto.In{
				Value:       5000,
				Date:        *sd,
				Coefficient: 1.0,
			},
			setup: func(f fields) {
				f.u.EXPECT().Do(gomock.Any(), gomock.Any()).Return(&add_bonus_dto.Out{Ok: true}, nil)
			},
			wantStatus:     http.StatusOK,
			wantInResponse: `"ok":true`,
		},
		{
			name: "duplicate error",
			reqBody: add_bonus_dto.In{
				Value:       5000,
				Date:        *sd,
				Coefficient: 1.0,
			},
			setup: func(f fields) {
				f.u.EXPECT().Do(gomock.Any(), gomock.Any()).Return(nil, addBonusUC.ErrDuplicateBonus)
			},
			wantStatus:     http.StatusBadRequest,
			wantInResponse: "bonus already exists for this date",
		},
		{
			name:       "invalid json",
			reqBody:    "invalid",
			setup:      func(f fields) {},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "internal error",
			reqBody: add_bonus_dto.In{
				Value:       5000,
				Date:        *sd,
				Coefficient: 1.0,
			},
			setup: func(f fields) {
				f.u.EXPECT().Do(gomock.Any(), gomock.Any()).Return(nil, errors.New("uc error"))
			},
			wantStatus:     http.StatusInternalServerError,
			wantInResponse: "internal error",
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

			h := add_bonus.NewHandler(f.u)

			var body []byte
			if s, ok := tt.reqBody.(string); ok {
				body = []byte(s)
			} else {
				body, _ = json.Marshal(tt.reqBody)
			}

			req := httptest.NewRequest(http.MethodPost, "/bonuses", bytes.NewReader(body))
			w := httptest.NewRecorder()

			h.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
			if tt.wantInResponse != "" {
				assert.Contains(t, w.Body.String(), tt.wantInResponse)
			}
		})
	}
}
