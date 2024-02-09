package app

import (
	"fmt"
	"net/http"

	"go.uber.org/zap"

	"github.com/yizeng/gab/chi/gorm/auth-jwt/internal/api"
	"github.com/yizeng/gab/chi/gorm/auth-jwt/internal/config"
	"github.com/yizeng/gab/chi/gorm/auth-jwt/internal/db"
	"github.com/yizeng/gab/chi/gorm/auth-jwt/internal/logger"
)

func Start() error {
	conf, err := config.Load("./cmd/app/config.yml")
	if err != nil {
		return fmt.Errorf("failed to initialize config -> %w", err)
	}

	if err = logger.Init(conf.API.Environment); err != nil {
		return fmt.Errorf("failed to initialize logger -> %w", err)
	}

	postgresDB, err := db.OpenPostgres(conf.Postgres)
	if err != nil {
		return fmt.Errorf("failed to initialize database -> %w", err)
	}

	s := api.NewServer(conf, postgresDB)

	addr := ":" + s.Config.API.Port
	zap.L().Info(fmt.Sprintf("starting server at %v", addr))
	if err = http.ListenAndServe(addr, s.Router); err != nil {
		return fmt.Errorf("failed to start the server -> %w", err)
	}

	return nil
}
