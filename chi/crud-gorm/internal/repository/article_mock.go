package repository

import (
	"context"

	"github.com/yizeng/gab/chi/crud-gorm/internal/domain"
)

type ArticleRepositoryMock struct {
	MockCreate   func(ctx context.Context, article *domain.Article) (*domain.Article, error)
	MockFindByID func(ctx context.Context, id uint) (*domain.Article, error)
	MockFindAll  func(ctx context.Context) ([]domain.Article, error)
}

func (m *ArticleRepositoryMock) Create(ctx context.Context, article *domain.Article) (*domain.Article, error) {
	return m.MockCreate(ctx, article)
}

func (m *ArticleRepositoryMock) FindByID(ctx context.Context, id uint) (*domain.Article, error) {
	return m.MockFindByID(ctx, id)
}

func (m *ArticleRepositoryMock) FindAll(ctx context.Context) ([]domain.Article, error) {
	return m.MockFindAll(ctx)
}
