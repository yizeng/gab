package main

import (
	"fmt"
	"net/http"

	"github.com/yizeng/gab/chi/crud-gorm/internal/config"
	"github.com/yizeng/gab/chi/crud-gorm/internal/db"
	"github.com/yizeng/gab/chi/crud-gorm/internal/logger"
	"github.com/yizeng/gab/chi/crud-gorm/internal/web"

	"go.uber.org/zap"
)

func main() {
	conf, err := config.Load("./config.yml")
	if err != nil {
		panic(fmt.Sprintf("failed to initialize config -> %v", err))
	}

	if err = logger.Init(conf.API.Environment); err != nil {
		panic(fmt.Sprintf("failed to initialize logger -> %v", err))
	}

	postgresDB, err := db.OpenPostgres(conf.Postgres)
	if err != nil {
		panic(fmt.Sprintf("failed to initialize database -> %v", err))
	}

	s := web.NewServer(conf.API, postgresDB)

	zap.L().Info("starting server at...", zap.String("address", s.Address))

	if err = http.ListenAndServe(s.Address, s.Router); err != nil {
		panic(fmt.Sprintf("failed to start the server -> %v", err))
	}
}
