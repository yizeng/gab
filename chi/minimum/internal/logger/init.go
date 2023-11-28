package logger

import (
	"go.uber.org/zap"

	"github.com/spf13/viper"
)

func Init() error {
	logger, err := zap.NewProduction()
	if err != nil {
		return err
	}

	if env := viper.GetString("API_ENV"); env == "development" {
		logger, err = zap.NewDevelopment()

		if err != nil {
			return err
		}
	}

	defer logger.Sync()

	zap.ReplaceGlobals(logger)

	return nil
}
