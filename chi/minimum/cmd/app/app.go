package app

import (
	"fmt"
	"net/http"

	"go.uber.org/zap"

	"github.com/yizeng/gab/chi/minimum/internal/api"
	"github.com/yizeng/gab/chi/minimum/internal/config"
	"github.com/yizeng/gab/chi/minimum/internal/logger"
)

func Start() {
	conf, err := config.Load()
	if err != nil {
		panic(fmt.Sprintf("failed to initialize config -> %v", err))
	}

	err = logger.Init(conf.API.Environment)
	if err != nil {
		panic(fmt.Sprintf("failed to initialize logger -> %v", err))
	}

	s := api.NewServer(conf.API)

	zap.L().Info(fmt.Sprintf("starting server at %v", s.Address))
	err = http.ListenAndServe(s.Address, s.Router)
	if err != nil {
		panic(fmt.Sprintf("failed to start the server -> %v", err))
	}
}
