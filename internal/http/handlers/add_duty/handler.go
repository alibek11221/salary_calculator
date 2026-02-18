package add_duty

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"salary_calculator/internal/dto/add_duty"
	"salary_calculator/internal/pkg/http/response"
	addDutyUC "salary_calculator/internal/usecase/duties/add"
)

type Handler struct {
	usecase usecase
}

func NewHandler(usecase usecase) *Handler {
	return &Handler{usecase: usecase}
}

// ServeHTTP godoc
// @Summary      Добавить дежурство
// @Description  Создает новую запись о дежурстве
// @Tags         duties
// @Accept       json
// @Produce      json
// @Param        input body add_duty.In true "Данные дежурства"
// @Success      200  {object}  add_duty.Out
// @Failure      400  {object}  map[string]string "error"
// @Failure      500  {object}  map[string]string "error"
// @Router       /duties/ [post]
func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req add_duty.In
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, fmt.Sprintf("invalid JSON: %s", err.Error()))
		return
	}
	out, err := h.usecase.Do(r.Context(), req)
	if err != nil {
		if errors.Is(err, addDutyUC.ErrDuplicateDuty) {
			response.BadRequest(w, "duty already exists for this date")
			return
		}

		response.InternalServerError(w, "internal error")
		return
	}

	response.Ok(w, out)
}
