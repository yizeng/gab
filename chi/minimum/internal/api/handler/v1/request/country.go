package request

import (
	"net/http"

	"github.com/yizeng/gab/chi/minimum/internal/domain"
)

type SumPopulationByState struct {
	States []domain.State `json:"states" validate:"required"`
}

func (c *SumPopulationByState) Bind(r *http.Request) error {
	return nil
}
