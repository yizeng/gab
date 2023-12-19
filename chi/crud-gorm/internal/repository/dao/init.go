package dao

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/spf13/viper"
)

func InitDB(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("gorm.Open -> %w", err)
	}

	err = initTables(db)
	if err != nil {
		return nil, fmt.Errorf("dao.InitTables -> %w", err)
	}
	return db, nil
}

func BuildDSNFromENV() string {
	host := viper.GetString("POSTGRES_HOST")
	user := viper.GetString("POSTGRES_USER")
	password := viper.GetString("POSTGRES_PASSWORD")
	dbName := viper.GetString("POSTGRES_DB_NAME")
	port := viper.GetInt("POSTGRES_PORT")

	dsn := fmt.Sprintf(
		"host=%v user=%v password=%v dbname=%v port=%v sslmode=disable",
		host, user, password, dbName, port,
	)

	return dsn
}

func initTables(db *gorm.DB) error {
	return db.AutoMigrate()
}
