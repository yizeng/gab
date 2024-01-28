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

	addr := ":" + s.Config.Port
	zap.L().Info(fmt.Sprintf("starting server at %v", addr))
	if err = http.ListenAndServe(addr, s.Router); err != nil {
		panic(fmt.Sprintf("failed to start the server -> %v", err))
	}
}
