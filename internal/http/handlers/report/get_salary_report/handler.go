package get_salary_report

import (
	"net/http"
	"strconv"

	"salary_calculator/internal/dto/get_salary_report"
	"salary_calculator/internal/pkg/http/response"
)

type Handler struct {
	u usecase
}

func New(u usecase) *Handler {
	return &Handler{u}
}

// ServeHTTP godoc
// @Summary      Отчет по зарплате
// @Description  Рассчитывает зарплату за указанный месяц и год
// @Tags         salary
// @Produce      json
// @Param        year   query     int  true  "Год (напр. 2024)"
// @Param        month  query     int  true  "Месяц (1-12)"
// @Success      200    {object}  get_salary_report.Out
// @Failure      400    {object}  map[string]string "error"
// @Failure      500    {object}  map[string]string "error"
// @Router       /salary/report [get]
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	yearStr := query.Get("year")
	monthStr := query.Get("month")

	if yearStr == "" || monthStr == "" {
		response.BadRequest(w, "missing year or month  parameter")
		return
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		response.BadRequest(w, "invalid year parameter")
		return
	}

	month, err := strconv.Atoi(monthStr)
	if err != nil {
		response.BadRequest(w, "invalid month parameter")
		return
	}

	req := get_salary_report.In{Year: year, Month: month}

	out, err := h.u.Do(r.Context(), req)
	if err != nil {
		response.InternalServerError(w, err.Error())
		return
	}

	response.Ok(w, out)
}
