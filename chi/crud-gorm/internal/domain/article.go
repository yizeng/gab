package domain

type Article struct {
	ID uint `json:"id"`

	Title   string `json:"title"`
	Content string `json:"content"`
}
