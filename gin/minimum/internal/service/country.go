package service

import (
	"context"

	"github.com/yizeng/gab/gin/minimum/internal/domain"
)

type CountryService struct{}

func NewCountryService() *CountryService {
	return &CountryService{}
}

func (s *CountryService) SumPopulationByState(ctx context.Context, states []domain.State) uint {
	var totalCount uint
	for _, v := range states {
		totalCount += v.Population
	}

	return totalCount
}
