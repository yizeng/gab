package request

import (
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation"
)

const (
	maxTitleLength   = 128
	maxContentLength = 5000
)

type CreateArticleRequest struct {
	UserID uint `json:"user_id" validate:"required"`

	Title   string `json:"title" validate:"required"`
	Content string `json:"content" validate:"required"`
}

func (req *CreateArticleRequest) Validate() error {
	return validation.ValidateStruct(
		req,
		validation.Field(&req.UserID, validation.Required, validation.Min(uint(1))),
		validation.Field(&req.Title, validation.Required, validation.Length(1, maxTitleLength)),
		validation.Field(&req.Content, validation.Required, validation.Length(1, maxContentLength)),
	)
}

func (req *CreateArticleRequest) Bind(r *http.Request) error {
	err := req.Validate()
	if err != nil {
		return err
	}

	return nil
}
