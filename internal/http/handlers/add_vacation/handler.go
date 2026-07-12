package add_vacation

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"salary_calculator/internal/dto/add_vacation"
	"salary_calculator/internal/pkg/http/response"
	"salary_calculator/internal/services/vacation_pay"
	addVacationUC "salary_calculator/internal/usecase/vacations/add"
)

type Handler struct {
	usecase usecase
}

func NewHandler(usecase usecase) *Handler {
	return &Handler{usecase: usecase}
}

// ServeHTTP godoc
// @Summary      Добавить отпуск
// @Description  Создает период отпуска. Пересекающиеся периоды запрещены.
// @Tags         vacations
// @Accept       json
// @Produce      json
// @Param        input body add_vacation.In true "Период отпуска (даты в формате YYYY-MM-DD)"
// @Success      200  {object}  add_vacation.Out
// @Failure      400  {object}  map[string]string "error"
// @Failure      409  {object}  map[string]string "error"
// @Failure      500  {object}  map[string]string "error"
// @Router       /vacations/ [post]
func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req add_vacation.In
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, fmt.Sprintf("invalid JSON: %s", err.Error()))
		return
	}

	out, err := h.usecase.Do(r.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, vacation_pay.ErrInvalidPeriod):
			response.BadRequest(w, "invalid vacation period")
		case errors.Is(err, addVacationUC.ErrOverlappingVacation):
			response.Conflict(w, "vacation overlaps with existing one")
		default:
			response.InternalServerError(w, "internal error")
		}
		return
	}

	response.Ok(w, out)
}
