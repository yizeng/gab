package request

import (
	"github.com/yizeng/gab/gin/minimum/internal/domain"
)

type SumPopulationByState struct {
	States []domain.State `json:"states" binding:"required,dive"`
}
