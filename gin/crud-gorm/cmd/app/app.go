package app

import (
	"fmt"

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

	addr := ":" + s.Config.API.Port
	zap.L().Info(fmt.Sprintf("starting server at %v", addr))
	if err = s.Router.Run(addr); err != nil {
		panic(fmt.Sprintf("failed to start the server -> %v", err))
	}
}
