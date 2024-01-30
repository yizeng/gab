package app

import (
	"fmt"
	"net/http"

	"go.uber.org/zap"

	"github.com/yizeng/gab/chi/minimum/internal/api"
	"github.com/yizeng/gab/chi/minimum/internal/config"
	"github.com/yizeng/gab/chi/minimum/internal/logger"
)

func Start() error {
	conf, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to initialize config -> %w", err)
	}

	if err = logger.Init(conf.API.Environment); err != nil {
		return fmt.Errorf("failed to initialize logger -> %w", err)
	}

	s := api.NewServer(conf.API)

	addr := ":" + s.Config.Port
	zap.L().Info(fmt.Sprintf("starting server at %v", addr))
	if err = http.ListenAndServe(addr, s.Router); err != nil {
		return fmt.Errorf("failed to start the server -> %w", err)
	}

	return nil
}
