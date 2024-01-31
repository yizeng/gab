package v1

import (
	"context"
	"net/http"

	"github.com/yizeng/gab/chi/minimum/internal/api/handler/v1/request"
	"github.com/yizeng/gab/chi/minimum/internal/api/handler/v1/response"
	"github.com/yizeng/gab/chi/minimum/internal/domain"

	"github.com/go-chi/render"
)

type CountryService interface {
	SumPopulationByState(ctx context.Context, states []domain.State) uint
}

type CountryHandler struct {
	svc CountryService
}

func NewCountryHandler(svc CountryService) *CountryHandler {
	return &CountryHandler{
		svc: svc,
	}
}

// HandleSumPopulationByState godoc
// @Summary      Sum the total population by state
// @Tags         countries
// @Produce      json
// @Param        request   body      request.SumPopulationByState true "request body"
// @Success      200      {object}   []domain.State
// @Failure      400      {object}   response.ErrResponse
// @Failure      500      {object}   response.ErrResponse
// @Router       /countries/sum-population-by-state [post]
func (h *CountryHandler) HandleSumPopulationByState(w http.ResponseWriter, r *http.Request) {
	req := request.SumPopulationByState{}
	if err := render.Bind(r, &req); err != nil {
		_ = render.Render(w, r, response.NewBadRequest(err))

		return
	}

	totalPopulation := h.svc.SumPopulationByState(r.Context(), req.States)

	if err := render.Render(w, r, response.NewSumPopulationByStateResult(totalPopulation)); err != nil {
		_ = render.Render(w, r, response.NewInternalServerError(err))

		return
	}
}
