package logger

import (
	"strings"

	"go.uber.org/zap"

	"github.com/spf13/viper"
)

func Init() error {
	logger, err := zap.NewProduction()
	if err != nil {
		return err
	}

	if strings.EqualFold(viper.GetString("API_ENV"), "development") {
		logger, err = zap.NewDevelopment()

		if err != nil {
			return err
		}
	}

	defer logger.Sync()

	zap.ReplaceGlobals(logger)

	return nil
}
