package edit_duty

import (
	"encoding/json"
	"fmt"
	"net/http"

	"salary_calculator/internal/dto/edit_duty"
	"salary_calculator/internal/pkg/http/response"
)

type Handler struct {
	usecase usecase
}

func NewHandler(usecase usecase) *Handler {
	return &Handler{usecase: usecase}
}

// ServeHTTP godoc
// @Summary      Редактировать дежурство
// @Description  Обновляет данные существующего дежурства
// @Tags         duties
// @Accept       json
// @Produce      json
// @Param        input body edit_duty.In true "Данные дежурства"
// @Success      200  {object}  edit_duty.Out
// @Failure      400  {object}  map[string]string "error"
// @Failure      500  {object}  map[string]string "error"
// @Router       /duties/ [put]
func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req edit_duty.In
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
