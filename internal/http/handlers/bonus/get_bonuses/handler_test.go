package get_bonuses_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	get_bonuses_dto "salary_calculator/internal/dto/get_bonuses"
	"salary_calculator/internal/http/handlers/bonus/get_bonuses"

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
		wantStatus     int
		wantInResponse string
	}{
		{
			name: "success",
			setup: func(f fields) {
				f.u.EXPECT().Do(gomock.Any()).Return(&get_bonuses_dto.Out{
					Bonuses: []get_bonuses_dto.Bonus{
						{ID: "1", Value: 5000},
					},
				}, nil)
			},
			wantStatus:     http.StatusOK,
			wantInResponse: `"value":5000`,
		},
		{
			name: "internal error",
			setup: func(f fields) {
				f.u.EXPECT().Do(gomock.Any()).Return(nil, errors.New("uc error"))
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

			h := get_bonuses.New(f.u)

			req := httptest.NewRequest(http.MethodGet, "/bonuses", nil)
			w := httptest.NewRecorder()

			h.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
			if tt.wantInResponse != "" {
				assert.Contains(t, w.Body.String(), tt.wantInResponse)
			}
		})
	}
}
