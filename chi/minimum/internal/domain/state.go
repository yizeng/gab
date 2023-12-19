package domain

type State struct {
	Name       string `json:"name" validate:"required"`
	Population uint   `json:"population" validate:"required"`
}
