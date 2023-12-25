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

	zap.L().Info(fmt.Sprintf("starting server at %v", s.Address))
	err = s.Router.Run(s.Address)
	if err != nil {
		panic(fmt.Sprintf("failed to start the server -> %v", err))
	}
}
