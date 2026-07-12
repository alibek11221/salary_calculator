package delete_vacation

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"salary_calculator/internal/dto/delete_vacation"
	"salary_calculator/internal/pkg/http/response"
	deleteVacationUC "salary_calculator/internal/usecase/vacations/delete"
)

type Handler struct {
	usecase usecase
}

func NewHandler(usecase usecase) *Handler {
	return &Handler{usecase: usecase}
}

// ServeHTTP godoc
// @Summary      Удалить отпуск
// @Description  Удаляет период отпуска по id
// @Tags         vacations
// @Accept       json
// @Produce      json
// @Param        input body delete_vacation.In true "ID отпуска"
// @Success      200  {object}  delete_vacation.Out
// @Failure      400  {object}  map[string]string "error"
// @Failure      500  {object}  map[string]string "error"
// @Router       /vacations/ [delete]
func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req delete_vacation.In
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, fmt.Sprintf("invalid JSON: %s", err.Error()))
		return
	}

	out, err := h.usecase.Do(r.Context(), req)
	if err != nil {
		if errors.Is(err, deleteVacationUC.ErrInvalidID) {
			response.BadRequest(w, err.Error())
			return
		}

		response.InternalServerError(w, "internal error")
		return
	}

	response.Ok(w, out)
}
