package list_bonuses

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
// @Summary      Список бонусов
// @Description  Возвращает список всех зарегистрированных бонусов
// @Tags         bonuses
// @Produce      json
// @Success      200  {object}  salary_calculator_internal_dto_list_bonuses.Out
// @Failure      500  {object}  map[string]string "error"
// @Router       /bonuses/ [get]
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	out, err := h.u.Do(r.Context())
	if err != nil {
		response.InternalServerError(w, err.Error())
		return
	}

	response.Ok(w, out)
}
