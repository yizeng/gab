package app

import (
	"fmt"

	"go.uber.org/zap"

	"github.com/yizeng/gab/gin/minimum/internal/api"
	"github.com/yizeng/gab/gin/minimum/internal/config"
	"github.com/yizeng/gab/gin/minimum/internal/logger"
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

	s := api.NewServer(conf)

	addr := ":" + s.Config.API.Port
	zap.L().Info(fmt.Sprintf("starting server at %v", addr))
	if err = s.Router.Run(addr); err != nil {
		panic(fmt.Sprintf("failed to start the server -> %v", err))
	}
}
