package response

import (
	"net/http"
)

type SumPopulationByStateResult struct {
	TotalPopulation uint `json:"total_population"`
}

func NewSumPopulationByStateResult(result uint) *SumPopulationByStateResult {
	resp := &SumPopulationByStateResult{
		TotalPopulation: result,
	}

	return resp
}

func (resp *SumPopulationByStateResult) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
