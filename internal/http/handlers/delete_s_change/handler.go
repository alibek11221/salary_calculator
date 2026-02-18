package delete_s_change

import (
	"encoding/json"
	"fmt"
	"net/http"

	"salary_calculator/internal/dto/delete_salary_change"
	"salary_calculator/internal/pkg/http/response"
)

type Handler struct {
	usecase usecase
}

func NewHandler(usecase usecase) *Handler {
	return &Handler{usecase: usecase}
}

// ServeHTTP godoc
// @Summary      Удалить изменение зарплаты
// @Description  Удаляет запись об изменении оклада
// @Tags         salary
// @Accept       json
// @Produce      json
// @Param        input body delete_salary_change.In true "ID записи"
// @Success      200  {object}  delete_salary_change.Out
// @Failure      400  {object}  map[string]string "error"
// @Failure      500  {object}  map[string]string "error"
// @Router       /salary/changes/ [delete]
func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req delete_salary_change.In
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
