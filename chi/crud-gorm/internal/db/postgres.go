package db

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/yizeng/gab/chi/crud-gorm/internal/config"
	"github.com/yizeng/gab/chi/crud-gorm/internal/repository/dao"
)

func OpenPostgres(conf *config.PostgresConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%v port=%v user=%v password=%v dbname=%v sslmode=disable",
		conf.Host, conf.Port, conf.User, conf.Password, conf.DB,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("gorm.Open -> %w", err)
	}

	if err = dao.InitTables(db); err != nil {
		return nil, fmt.Errorf("dao.InitTables -> %w", err)
	}

	return db, nil
}
