package domain

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

type State struct {
	Name       string `json:"name" validate:"required"`
	Population uint   `json:"population" validate:"required"`
}

func (s *State) Validate() error {
	return validation.ValidateStruct(
		s,
		validation.Field(&s.Name, validation.Required),
		validation.Field(&s.Population, validation.Required),
	)
}
