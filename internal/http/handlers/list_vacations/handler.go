package list_vacations

import (
	"net/http"

	"salary_calculator/internal/pkg/http/response"
)

type Handler struct {
	u usecase
}

func New(u usecase) *Handler {
	return &Handler{u}
}

// ServeHTTP godoc
// @Summary      Список отпусков
// @Description  Возвращает все периоды отпусков
// @Tags         vacations
// @Produce      json
// @Success      200  {object}  salary_calculator_internal_dto_list_vacations.Out
// @Failure      500  {object}  map[string]string "error"
// @Router       /vacations/ [get]
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	out, err := h.u.Do(r.Context())
	if err != nil {
		response.InternalServerError(w, "internal error")
		return
	}

	response.Ok(w, out)
}
