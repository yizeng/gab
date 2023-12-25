package response

type SumPopulationByStateResult struct {
	TotalPopulation uint `json:"total_population"`
}

func NewSumPopulationByStateResult(result uint) *SumPopulationByStateResult {
	resp := &SumPopulationByStateResult{
		TotalPopulation: result,
	}

	return resp
}
