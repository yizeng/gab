package request

import (
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation"

	"github.com/yizeng/gab/chi/minimum/internal/domain"
)

type SumPopulationByState struct {
	States []domain.State `json:"states" validate:"required"`
}

func (req *SumPopulationByState) Validate() error {
	return validation.ValidateStruct(
		req,
		validation.Field(&req.States, validation.Required),
		validation.Field(&req.States, validation.By(eachState)),
	)
}

func eachState(value any) error {
	states, _ := value.([]domain.State)
	for _, s := range states {
		err := s.Validate()
		if err != nil {
			return err
		}
	}

	return nil
}

func (req *SumPopulationByState) Bind(r *http.Request) error {
	err := req.Validate()
	if err != nil {
		return err
	}

	return nil
}
