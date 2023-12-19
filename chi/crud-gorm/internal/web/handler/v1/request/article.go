package request

import (
	"errors"
	"net/http"

	"github.com/yizeng/gab/chi/crud-gorm/internal/domain"
)

type CreateArticleRequest struct {
	domain.Article
}

func (req *CreateArticleRequest) Bind(r *http.Request) error {
	if req.UserID == 0 {
		return errors.New("article.user_id cannot be empty")
	}
	if req.Title == "" {
		return errors.New("article.title cannot be empty")
	}
	if req.Content == "" {
		return errors.New("article.content cannot be empty")
	}

	return nil
}
