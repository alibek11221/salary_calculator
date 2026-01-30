package get_salary_changes_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	get_salary_changes_dto "salary_calculator/internal/dto/get_salary_changes"
	"salary_calculator/internal/http/handlers/salary_change/get_salary_changes"

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
				f.u.EXPECT().Do(gomock.Any()).Return(&get_salary_changes_dto.Out{
					Changes: []get_salary_changes_dto.Change{
						{ID: "1", Value: 1000},
					},
				}, nil)
			},
			wantStatus:     http.StatusOK,
			wantInResponse: `"value":1000`,
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

			h := get_salary_changes.New(f.u)
			req := httptest.NewRequest(http.MethodGet, "/changes", nil)
			w := httptest.NewRecorder()

			h.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
			if tt.wantInResponse != "" {
				assert.Contains(t, w.Body.String(), tt.wantInResponse)
			}
		})
	}
}
