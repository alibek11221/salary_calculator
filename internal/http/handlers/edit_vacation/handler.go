package edit_vacation

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"salary_calculator/internal/dto/edit_vacation"
	"salary_calculator/internal/pkg/http/response"
	"salary_calculator/internal/services/vacation_pay"
	editVacationUC "salary_calculator/internal/usecase/vacations/edit"
)

type Handler struct {
	usecase usecase
}

func NewHandler(usecase usecase) *Handler {
	return &Handler{usecase: usecase}
}

// ServeHTTP godoc
// @Summary      Изменить отпуск
// @Description  Обновляет период отпуска по id
// @Tags         vacations
// @Accept       json
// @Produce      json
// @Param        input body edit_vacation.In true "Отпуск (id + даты YYYY-MM-DD)"
// @Success      200  {object}  edit_vacation.Out
// @Failure      400  {object}  map[string]string "error"
// @Failure      409  {object}  map[string]string "error"
// @Failure      500  {object}  map[string]string "error"
// @Router       /vacations/ [put]
func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req edit_vacation.In
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, fmt.Sprintf("invalid JSON: %s", err.Error()))
		return
	}

	out, err := h.usecase.Do(r.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, vacation_pay.ErrInvalidPeriod), errors.Is(err, editVacationUC.ErrInvalidID):
			response.BadRequest(w, err.Error())
		case errors.Is(err, editVacationUC.ErrOverlappingVacation):
			response.Conflict(w, "vacation overlaps with existing one")
		default:
			response.InternalServerError(w, "internal error")
		}
		return
	}

	response.Ok(w, out)
}
