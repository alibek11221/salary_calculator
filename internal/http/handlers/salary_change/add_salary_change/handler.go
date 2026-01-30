package add_salary_change

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"salary_calculator/internal/dto/add_salay_change"
	"salary_calculator/internal/pkg/http/response"
	"salary_calculator/internal/usecase/add_salary_change"
)

type Handler struct {
	usecase usecase
}

func NewHandler(usecase usecase) *Handler {
	return &Handler{usecase: usecase}
}

// ServeHTTP godoc
// @Summary      Добавить изменение зарплаты
// @Description  Создает новую запись об изменении оклада
// @Tags         salary
// @Accept       json
// @Produce      json
// @Param        input body add_salay_change.In true "Данные изменения"
// @Success      200  {object}  add_salay_change.Out
// @Failure      400  {object}  map[string]string "error"
// @Failure      500  {object}  map[string]string "error"
// @Router       /salary/changes/ [post]
func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req add_salay_change.In
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, fmt.Sprintf("invalid JSON: %s", err.Error()))
		return
	}
	out, err := h.usecase.Do(r.Context(), req)
	if err != nil {
		if errors.Is(err, add_salary_change.ErrDuplicateSalaryChange) {
			response.BadRequest(w, "salary change already exists for this date")
			return
		}

		response.InternalServerError(w, "internal error")
		return
	}

	response.Ok(w, out)
}
