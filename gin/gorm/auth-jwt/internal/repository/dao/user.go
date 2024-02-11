package dao

import (
	"context"
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrUserEmailExists = errors.New("user already exists")
	ErrUserNotFound    = errors.New("user not found")
)

type User struct {
	ID uint `gorm:"primaryKey"`

	Email    string `gorm:"unique;not null"`
	Password string `gorm:"not null"`

	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
}

type UserDAO struct {
	db *gorm.DB
}

func NewUserDAO(db *gorm.DB) *UserDAO {
	return &UserDAO{
		db: db,
	}
}

func (d *UserDAO) Insert(ctx context.Context, user User) (User, error) {
	result := d.db.WithContext(ctx).Create(&user)
	if result.Error != nil {
		var err *pgconn.PgError
		if errors.As(result.Error, &err) &&
			err.Code == pgerrcode.UniqueViolation &&
			strings.Contains(err.Message, `unique constraint "users_email_key"`) {
			return User{}, ErrUserEmailExists
		}

		return User{}, result.Error
	}

	return user, nil
}

func (d *UserDAO) FindByID(ctx context.Context, id uint) (User, error) {
	var user User

	result := d.db.WithContext(ctx).First(&user, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return User{}, ErrUserNotFound
		}

		return User{}, result.Error
	}

	return user, nil
}

func (d *UserDAO) FindByEmail(ctx context.Context, email string) (User, error) {
	var user User

	result := d.db.WithContext(ctx).First(&user, "email = ?", email)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return User{}, ErrUserNotFound
		}

		return User{}, result.Error
	}

	return user, nil
}
