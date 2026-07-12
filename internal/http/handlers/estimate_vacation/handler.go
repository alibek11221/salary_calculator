package estimate_vacation

import (
	"errors"
	"net/http"

	"salary_calculator/internal/dto/estimate_vacation"
	"salary_calculator/internal/pkg/http/response"
	"salary_calculator/internal/pkg/types"
	"salary_calculator/internal/services/vacation_pay"
)

type Handler struct {
	u usecase
}

func NewHandler(u usecase) *Handler {
	return &Handler{u}
}

// ServeHTTP godoc
// @Summary      Прикидка отпускных
// @Description  Считает отпускные и зарплату затронутых месяцев для гипотетического отпуска, ничего не сохраняя
// @Tags         vacations
// @Produce      json
// @Param        from  query  string  true  "Начало отпуска (YYYY-MM-DD)"
// @Param        to    query  string  true  "Конец отпуска (YYYY-MM-DD)"
// @Success      200  {object}  estimate_vacation.Out
// @Failure      400  {object}  map[string]string "error"
// @Failure      500  {object}  map[string]string "error"
// @Router       /vacations/estimate [get]
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	from, err := types.ParseISODate(query.Get("from"))
	if err != nil {
		response.BadRequest(w, "invalid 'from' parameter: expected YYYY-MM-DD")
		return
	}

	to, err := types.ParseISODate(query.Get("to"))
	if err != nil {
		response.BadRequest(w, "invalid 'to' parameter: expected YYYY-MM-DD")
		return
	}

	out, err := h.u.Do(r.Context(), estimate_vacation.In{From: from, To: to})
	if err != nil {
		switch {
		case errors.Is(err, vacation_pay.ErrInvalidPeriod):
			response.BadRequest(w, "invalid vacation period")
		case errors.Is(err, vacation_pay.ErrNoEarningsData):
			response.InternalServerError(w, "no salary history to compute average earnings")
		default:
			response.InternalServerError(w, err.Error())
		}
		return
	}

	response.Ok(w, out)
}
