package repository

import (
	"context"
	"fmt"

	"github.com/yizeng/gab/chi/crud-gorm/internal/domain"
	"github.com/yizeng/gab/chi/crud-gorm/internal/repository/dao"
)

var (
	ErrArticleDuplicated = dao.ErrArticleDuplicated
	ErrArticleNotFound   = dao.ErrArticleNotFound
)

type ArticleDAO interface {
	Insert(ctx context.Context, article dao.Article) (dao.Article, error)
	FindByID(ctx context.Context, id uint) (dao.Article, error)
	FindAll(ctx context.Context, page, perPage uint) ([]dao.Article, error)
	Search(ctx context.Context, title, content string) ([]dao.Article, error)
}

type ArticleRepository struct {
	dao ArticleDAO
}

func NewArticleRepository(dao ArticleDAO) *ArticleRepository {
	return &ArticleRepository{
		dao: dao,
	}
}

func (r *ArticleRepository) Create(ctx context.Context, article domain.Article) (domain.Article, error) {
	created, err := r.dao.Insert(ctx, dao.Article{
		UserID:  article.UserID,
		Title:   article.Title,
		Content: article.Content,
	})
	if err != nil {
		return domain.Article{}, fmt.Errorf("r.dao.Insert -> %w", err)
	}

	return daoToDomain(created), nil
}

func (r *ArticleRepository) FindByID(ctx context.Context, id uint) (domain.Article, error) {
	found, err := r.dao.FindByID(ctx, id)
	if err != nil {
		return domain.Article{}, fmt.Errorf("r.dao.FindByID -> %w", err)
	}

	return daoToDomain(found), nil
}

func (r *ArticleRepository) FindAll(ctx context.Context, page, perPage uint) ([]domain.Article, error) {
	allArticles, err := r.dao.FindAll(ctx, page, perPage)
	if err != nil {
		return nil, fmt.Errorf("r.dao.FindAll -> %w", err)
	}

	articles := make([]domain.Article, 0, len(allArticles))
	for _, a := range allArticles {
		articles = append(articles, daoToDomain(a))
	}

	return articles, nil
}

func (r *ArticleRepository) Search(ctx context.Context, title, content string) ([]domain.Article, error) {
	allArticles, err := r.dao.Search(ctx, title, content)
	if err != nil {
		return nil, fmt.Errorf("r.dao.Search -> %w", err)
	}

	articles := make([]domain.Article, 0, len(allArticles))
	for _, a := range allArticles {
		articles = append(articles, daoToDomain(a))
	}

	return articles, nil
}

func daoToDomain(a dao.Article) domain.Article {
	return domain.Article{
		ID:        a.ID,
		UserID:    a.UserID,
		Title:     a.Title,
		Content:   a.Content,
		CreatedAt: a.CreatedAt,
		UpdatedAt: a.UpdatedAt,
	}
}
