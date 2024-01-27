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

func (d *ArticleDAO) Insert(ctx context.Context, article Article) (Article, error) {
	result := d.db.WithContext(ctx).Create(&article)
	if result.Error != nil {
		var err *pgconn.PgError
		if errors.As(result.Error, &err) && err.Code == pgerrcode.UniqueViolation {
			return Article{}, ErrArticleDuplicated
		}

		return Article{}, result.Error
	}

	return article, nil
}

func (d *ArticleDAO) FindByID(ctx context.Context, id uint) (Article, error) {
	var article Article

	result := d.db.WithContext(ctx).First(&article, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return Article{}, ErrArticleNotFound
		}

		return Article{}, result.Error
	}

	return article, nil
}

func (d *ArticleDAO) FindAll(ctx context.Context, page, perPage uint) ([]Article, error) {
	var articles []Article

	// page number is starting from 1.
	offset := (page - 1) * perPage
	result := d.db.WithContext(ctx).Offset(int(offset)).Limit(int(perPage)).Find(&articles)
	if result.Error != nil {
		return nil, result.Error
	}

	return articles, nil
}

func (d *ArticleDAO) Search(ctx context.Context, title, content string) ([]Article, error) {
	var articles []Article

	finder := d.db.WithContext(ctx)
	if title != "" {
		finder = finder.Where("title LIKE ?", "%"+title+"%")
	}
	if content != "" {
		finder = finder.Where("content LIKE ?", "%"+content+"%")
	}

	result := finder.Find(&articles)
	if result.Error != nil {
		return nil, result.Error
	}

	return articles, nil
}
