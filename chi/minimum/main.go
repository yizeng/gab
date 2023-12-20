package main

import (
	"fmt"
	"net/http"

	"github.com/yizeng/gab/chi/minimum/internal/config"
	"github.com/yizeng/gab/chi/minimum/internal/logger"
	"github.com/yizeng/gab/chi/minimum/internal/web"

	"go.uber.org/zap"
)

func main() {
	conf, err := config.Load()
	if err != nil {
		panic(fmt.Sprintf("failed to initialize config -> %v", err))
	}

	err = logger.Init(conf.API.Environment)
	if err != nil {
		panic(fmt.Sprintf("failed to initialize logger -> %v", err))
	}

	s := web.NewServer(conf.API)

	zap.L().Info(fmt.Sprintf("starting server at %v", s.Address))
	err = http.ListenAndServe(s.Address, s.Router)
	if err != nil {
		panic(fmt.Sprintf("failed to start the server -> %v", err))
	}
}
