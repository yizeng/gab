package service

import (
	"context"
	"fmt"

	"github.com/yizeng/gab/chi/crud-gorm/internal/domain"
)

type ArticleService struct{}

func NewArticleService() *ArticleService {
	return &ArticleService{}
}

func (a *ArticleService) ListArticles(ctx context.Context) ([]domain.Article, error) {
	dummyArticles := []domain.Article{
		{
			ID:      1,
			Title:   "title 1",
			Content: "content 1",
		}, {
			ID:      2,
			Title:   "title 2",
			Content: "content 2",
		},
	}

	return dummyArticles, nil
}

func (a *ArticleService) CreateArticle(ctx context.Context, article domain.Article) (domain.Article, error) {
	return article, nil
}

func (a *ArticleService) GetArticle(ctx context.Context, id uint) (domain.Article, error) {
	return domain.Article{
		ID:      id,
		Title:   fmt.Sprintf("title %v", id),
		Content: fmt.Sprintf("content %v", id),
	}, nil
}
