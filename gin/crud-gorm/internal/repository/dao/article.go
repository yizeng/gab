package dao

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrArticleDuplicated = errors.New("article already exists")
	ErrArticleNotFound   = errors.New("article not found")
)

type Article struct {
	gorm.Model

	UserID  uint   `gorm:"uniqueIndex:idx_user_id_title,not null"`
	Title   string `gorm:"uniqueIndex:idx_user_id_title,not null"`
	Content string `gorm:"not null"`
}

type ArticleDAO struct {
	db *gorm.DB
}

func NewArticleDAO(db *gorm.DB) *ArticleDAO {
	return &ArticleDAO{
		db: db,
	}
}

func (d *ArticleDAO) Create(ctx context.Context, article *Article) (*Article, error) {
	result := d.db.WithContext(ctx).Create(article)
	if result.Error != nil {
		var err *pgconn.PgError
		if errors.As(result.Error, &err) && err.Code == pgerrcode.UniqueViolation {
			return nil, ErrArticleDuplicated
		}

		return nil, result.Error
	}

	return article, nil
}

func (d *ArticleDAO) FindByID(ctx context.Context, id uint) (*Article, error) {
	var article Article

	result := d.db.WithContext(ctx).First(&article, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrArticleNotFound
		}

		return nil, result.Error
	}

	return &article, nil
}

func (d *ArticleDAO) FindAll(ctx context.Context) ([]Article, error) {
	var articles []Article

	result := d.db.WithContext(ctx).Find(&articles)
	if result.Error != nil {
		return nil, result.Error
	}

	return articles, nil
}
