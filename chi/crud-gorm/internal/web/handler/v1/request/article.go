package request

import (
	"net/http"

	"github.com/yizeng/gab/chi/crud-gorm/internal/domain"
)

type CreateArticleRequest struct {
	domain.Article
}

func (req *CreateArticleRequest) Bind(r *http.Request) error {
	return nil
}
