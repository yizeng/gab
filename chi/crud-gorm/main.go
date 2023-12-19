package main

import (
	"fmt"
	"net/http"

	"github.com/yizeng/gab/chi/crud-gorm/internal/config"
	"github.com/yizeng/gab/chi/crud-gorm/internal/logger"
	"github.com/yizeng/gab/chi/crud-gorm/internal/repository/dao"
	"github.com/yizeng/gab/chi/crud-gorm/internal/web"

	"go.uber.org/zap"
)

func main() {
	err := config.Init()
	if err != nil {
		panic(fmt.Sprintf("failed to initialize config -> %v", err))
	}

	err = logger.Init()
	if err != nil {
		panic(fmt.Sprintf("failed to initialize logger -> %v", err))
	}

	dsn := dao.BuildDSNFromENV()
	db, err := dao.InitDB(dsn)
	if err != nil {
		panic(fmt.Sprintf("failed to initialize database -> %v", err))
	}

	s := web.NewServer(db)

	zap.L().Info("starting server at...", zap.String("address", s.Address))
	err = http.ListenAndServe(s.Address, s.Router)
	if err != nil {
		panic(fmt.Sprintf("failed to start the server -> %v", err))
	}
}
