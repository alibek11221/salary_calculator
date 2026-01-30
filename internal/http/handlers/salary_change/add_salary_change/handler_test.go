package add_salary_change_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"salary_calculator/internal/dto/add_salay_change"
	"salary_calculator/internal/dto/value_objects"
	"salary_calculator/internal/http/handlers/salary_change/add_salary_change"
	uc "salary_calculator/internal/usecase/add_salary_change"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestHandler_ServeHTTP(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := NewMockusecase(ctrl)
	handler := add_salary_change.NewHandler(mockUC)

	date, _ := value_objects.NewSalaryDate("2024_01")
	validIn := add_salay_change.In{
		Value: 1000,
		Date:  *date,
	}

	t.Run("Success", func(t *testing.T) {
		mockUC.EXPECT().Do(gomock.Any(), validIn).Return(&add_salay_change.Out{Ok: true}, nil)

		body, _ := json.Marshal(validIn)
		req := httptest.NewRequest(http.MethodPost, "/add", bytes.NewReader(body))
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp add_salay_change.Out
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.True(t, resp.Ok)
	})

	t.Run("Duplicate Error", func(t *testing.T) {
		mockUC.EXPECT().Do(gomock.Any(), validIn).Return(nil, uc.ErrDuplicateSalaryChange)

		body, _ := json.Marshal(validIn)
		req := httptest.NewRequest(http.MethodPost, "/add", bytes.NewReader(body))
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "already exists")
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/add", bytes.NewReader([]byte("{invalid")))
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Internal Error", func(t *testing.T) {
		mockUC.EXPECT().Do(gomock.Any(), validIn).Return(nil, errors.New("something went wrong"))

		body, _ := json.Marshal(validIn)
		req := httptest.NewRequest(http.MethodPost, "/add", bytes.NewReader(body))
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
