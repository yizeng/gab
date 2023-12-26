package app

import (
	"fmt"
	"net/http"

	"go.uber.org/zap"

	"github.com/yizeng/gab/gin/crud-gorm/internal/api"
	"github.com/yizeng/gab/gin/crud-gorm/internal/config"
	"github.com/yizeng/gab/gin/crud-gorm/internal/db"
	"github.com/yizeng/gab/gin/crud-gorm/internal/logger"
)

func Start() {
	conf, err := config.Load("./cmd/app/config.yml")
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

	s := api.NewServer(conf, postgresDB)

	zap.L().Info("starting server at...", zap.String("address", s.Address))

	if err = http.ListenAndServe(s.Address, s.Router); err != nil {
		panic(fmt.Sprintf("failed to start the server -> %v", err))
	}
}
