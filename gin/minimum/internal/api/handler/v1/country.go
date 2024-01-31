package v1

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/yizeng/gab/gin/minimum/internal/api/handler/v1/request"
	"github.com/yizeng/gab/gin/minimum/internal/api/handler/v1/response"
	"github.com/yizeng/gab/gin/minimum/internal/domain"
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
// @Router       /countries/sum-population-by-state [post]
func (h *CountryHandler) HandleSumPopulationByState(ctx *gin.Context) {
	req := request.SumPopulationByState{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.RenderError(ctx, response.NewBadRequest(err))

		return
	}

	totalPopulation := h.svc.SumPopulationByState(ctx.Request.Context(), req.States)

	ctx.JSON(http.StatusOK, response.NewSumPopulationByStateResult(totalPopulation))
}
