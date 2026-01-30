package response

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResponses(t *testing.T) {
	t.Run("Ok", func(t *testing.T) {
		w := httptest.NewRecorder()
		data := map[string]string{"foo": "bar"}
		err := Ok(w, data)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var resp map[string]string
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, data, resp)
	})

	t.Run("Created", func(t *testing.T) {
		w := httptest.NewRecorder()
		data := map[string]string{"foo": "bar"}
		err := Created(w, data)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("BadRequest", func(t *testing.T) {
		w := httptest.NewRecorder()
		err := BadRequest(w, "bad request")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		var resp map[string]string
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, "bad request", resp["error"])
	})

	t.Run("Conflict", func(t *testing.T) {
		w := httptest.NewRecorder()
		err := Conflict(w, "conflict")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusConflict, w.Code)
	})

	t.Run("InternalServerError", func(t *testing.T) {
		w := httptest.NewRecorder()
		err := InternalServerError(w, "internal error")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
