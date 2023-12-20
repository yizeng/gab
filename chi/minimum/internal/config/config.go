package config

type AppConfig struct {
	API *APIConfig
}

type APIConfig struct {
	Environment string `env:"API_ENV,required"`
	Host        string `env:"API_HOST,required"`
	Port        string `env:"API_PORT,required"`
}
