package domain

type State struct {
	Name       string `json:"name" binding:"required"`
	Population uint   `json:"population" binding:"required"`
}
