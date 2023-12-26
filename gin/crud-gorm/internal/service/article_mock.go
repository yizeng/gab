package service

import (
	"context"

	"github.com/yizeng/gab/gin/crud-gorm/internal/domain"
)

type ArticleServiceMock struct {
	MockCreate       func(ctx context.Context, article *domain.Article) (*domain.Article, error)
	MockGetArticle   func(ctx context.Context, id uint) (*domain.Article, error)
	MockListArticles func(ctx context.Context) ([]domain.Article, error)
}

func NewArticleServiceMock() *ArticleServiceMock {
	return &ArticleServiceMock{}
}

func (m *ArticleServiceMock) CreateArticle(ctx context.Context, article *domain.Article) (*domain.Article, error) {
	return m.MockCreate(ctx, article)
}

func (m *ArticleServiceMock) GetArticle(ctx context.Context, id uint) (*domain.Article, error) {
	return m.MockGetArticle(ctx, id)
}

func (m *ArticleServiceMock) ListArticles(ctx context.Context) ([]domain.Article, error) {
	return m.MockListArticles(ctx)
}
