package repository

import (
	"context"
	"fmt"

	"github.com/yizeng/gab/gin/crud-gorm/internal/domain"
	"github.com/yizeng/gab/gin/crud-gorm/internal/repository/dao"
)

var (
	ErrArticleDuplicated = dao.ErrArticleDuplicated
	ErrArticleNotFound   = dao.ErrArticleNotFound
)

type ArticleDAO interface {
	Create(ctx context.Context, article *dao.Article) (*dao.Article, error)
	FindByID(ctx context.Context, id uint) (*dao.Article, error)
	FindAll(ctx context.Context) ([]dao.Article, error)
}

type ArticleRepository struct {
	dao ArticleDAO
}

func NewArticleRepository(dao ArticleDAO) *ArticleRepository {
	return &ArticleRepository{
		dao: dao,
	}
}

func (r *ArticleRepository) Create(ctx context.Context, article *domain.Article) (*domain.Article, error) {
	articleDAO := &dao.Article{
		UserID:  article.UserID,
		Title:   article.Title,
		Content: article.Content,
	}

	created, err := r.dao.Create(ctx, articleDAO)
	if err != nil {
		return nil, fmt.Errorf("r.dao.Create -> %w", err)
	}

	return daoToDomain(created), nil
}

func (r *ArticleRepository) FindByID(ctx context.Context, id uint) (*domain.Article, error) {
	found, err := r.dao.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("r.dao.FindByID -> %w", err)
	}

	return daoToDomain(found), nil
}

func (r *ArticleRepository) FindAll(ctx context.Context) ([]domain.Article, error) {
	allArticles, err := r.dao.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("r.dao.FindAll -> %w", err)
	}

	articles := make([]domain.Article, 0, len(allArticles))
	for _, a := range allArticles {
		articles = append(articles, *daoToDomain(&a))
	}

	return articles, nil
}

func daoToDomain(a *dao.Article) *domain.Article {
	return &domain.Article{
		ID:        a.ID,
		UserID:    a.UserID,
		Title:     a.Title,
		Content:   a.Content,
		CreatedAt: a.CreatedAt,
		UpdatedAt: a.UpdatedAt,
	}
}
