package delete_bonus

import (
	"encoding/json"
	"fmt"
	"net/http"

	"salary_calculator/internal/dto/delete_bonus"
	"salary_calculator/internal/pkg/http/response"
)

type Handler struct {
	usecase usecase
}

func NewHandler(usecase usecase) *Handler {
	return &Handler{usecase: usecase}
}

// ServeHTTP godoc
// @Summary      Удалить бонус
// @Description  Удаляет запись о бонусе
// @Tags         bonuses
// @Accept       json
// @Produce      json
// @Param        input body delete_bonus.In true "ID бонуса"
// @Success      200  {object}  delete_bonus.Out
// @Failure      400  {object}  map[string]string "error"
// @Failure      500  {object}  map[string]string "error"
// @Router       /bonuses/ [delete]
func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req delete_bonus.In
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, fmt.Sprintf("invalid JSON: %s", err.Error()))
		return
	}
	out, err := h.usecase.Do(r.Context(), req)
	if err != nil {
		response.InternalServerError(w, err.Error())
		return
	}

	response.Ok(w, out)
}
