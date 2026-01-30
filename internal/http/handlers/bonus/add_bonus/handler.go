package add_bonus

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"salary_calculator/internal/dto/add_bonus"
	"salary_calculator/internal/pkg/http/response"
	addBonusUC "salary_calculator/internal/usecase/add_bonus"
)

type Handler struct {
	usecase usecase
}

func NewHandler(usecase usecase) *Handler {
	return &Handler{usecase: usecase}
}

// ServeHTTP godoc
// @Summary      Добавить бонус
// @Description  Создает новую запись о бонусе
// @Tags         bonuses
// @Accept       json
// @Produce      json
// @Param        input body add_bonus.In true "Данные бонуса"
// @Success      200  {object}  add_bonus.Out
// @Failure      400  {object}  map[string]string "error"
// @Failure      500  {object}  map[string]string "error"
// @Router       /bonuses/ [post]
func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req add_bonus.In
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, fmt.Sprintf("invalid JSON: %s", err.Error()))
		return
	}
	out, err := h.usecase.Do(r.Context(), req)
	if err != nil {
		if errors.Is(err, addBonusUC.ErrDuplicateBonus) {
			response.BadRequest(w, "bonus already exists for this date")
			return
		}

		response.InternalServerError(w, "internal error")
		return
	}

	response.Ok(w, out)
}
