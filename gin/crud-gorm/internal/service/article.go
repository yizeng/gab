package service

import (
	"context"
	"fmt"

	"github.com/yizeng/gab/gin/crud-gorm/internal/domain"
	"github.com/yizeng/gab/gin/crud-gorm/internal/repository"
)

var (
	ErrArticleDuplicated = repository.ErrArticleDuplicated
	ErrArticleNotFound   = repository.ErrArticleNotFound
)

type ArticleRepository interface {
	Create(ctx context.Context, article *domain.Article) (*domain.Article, error)
	FindByID(ctx context.Context, id uint) (*domain.Article, error)
	FindAll(ctx context.Context) ([]domain.Article, error)
	Search(ctx context.Context, title, content string) ([]domain.Article, error)
}

type ArticleService struct {
	repo ArticleRepository
}

func NewArticleService(repo ArticleRepository) *ArticleService {
	return &ArticleService{
		repo: repo,
	}
}

func (s *ArticleService) CreateArticle(ctx context.Context, article *domain.Article) (*domain.Article, error) {
	created, err := s.repo.Create(ctx, article)
	if err != nil {
		return nil, fmt.Errorf("s.repo.Create -> %w", err)
	}

	return created, nil
}

func (s *ArticleService) GetArticle(ctx context.Context, id uint) (*domain.Article, error) {
	article, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("s.repo.FindByID -> %w", err)
	}

	return article, nil
}

func (s *ArticleService) ListArticles(ctx context.Context) ([]domain.Article, error) {
	articles, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("s.repo.FindAll -> %w", err)
	}

	return articles, nil
}

func (s *ArticleService) SearchArticles(ctx context.Context, title, content string) ([]domain.Article, error) {
	articles, err := s.repo.Search(ctx, title, content)
	if err != nil {
		return nil, fmt.Errorf("s.repo.Search -> %w", err)
	}

	return articles, nil
}
