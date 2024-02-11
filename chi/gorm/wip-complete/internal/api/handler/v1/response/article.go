package response

import (
	"net/http"

	"github.com/go-chi/render"

	"github.com/yizeng/gab/chi/gorm/wip-complete/internal/domain"
)

type ArticleResponse struct {
	*domain.Article
}

func NewArticle(article *domain.Article) *ArticleResponse {
	resp := &ArticleResponse{
		Article: article,
	}

	return resp
}

func NewArticles(articles []domain.Article) []render.Renderer {
	list := make([]render.Renderer, 0, len(articles))
	for _, a := range articles {
		article := a
		list = append(list, NewArticle(&article))
	}

	return list
}

func (resp *ArticleResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
