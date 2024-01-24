package service

import (
	"context"

	"github.com/yizeng/gab/chi/crud-gorm/internal/domain"
)

type ArticleServiceMock struct {
	MockCreate         func(ctx context.Context, article *domain.Article) (*domain.Article, error)
	MockGetArticle     func(ctx context.Context, id uint) (*domain.Article, error)
	MockListArticles   func(ctx context.Context, page, perPage uint) ([]domain.Article, error)
	MockSearchArticles func(ctx context.Context, title, content string) ([]domain.Article, error)
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

func (m *ArticleServiceMock) ListArticles(ctx context.Context, page, perPage uint) ([]domain.Article, error) {
	return m.MockListArticles(ctx, page, perPage)
}

func (m *ArticleServiceMock) SearchArticles(ctx context.Context, title, content string) ([]domain.Article, error) {
	return m.MockSearchArticles(ctx, title, content)
}
