package v1

import (
	"context"
	"net/http"

	"github.com/yizeng/gab/chi/minimum/internal/domain"
	"github.com/yizeng/gab/chi/minimum/internal/web/handler/v1/request"
	"github.com/yizeng/gab/chi/minimum/internal/web/handler/v1/response"

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

func (h *CountryHandler) HandleSumPopulationByState(w http.ResponseWriter, r *http.Request) {
	req := request.SumPopulationByState{}
	if err := render.Bind(r, &req); err != nil {
		render.Render(w, r, response.NewBadRequest(err.Error()))

		return
	}

	totalPopulation := h.svc.SumPopulationByState(r.Context(), req.States)

	err := render.Render(w, r, response.NewSumPopulationByStateResult(totalPopulation))
	if err != nil {
		render.Render(w, r, response.NewInternalServerError(err))

		return
	}
}
