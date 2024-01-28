package config

type AppConfig struct {
	API *APIConfig
	Gin *GinConfig
}

type APIConfig struct {
	Environment        string   `env:"API_ENV,required"`
	Port               string   `env:"API_PORT,required"`
	BaseURL            string   `env:"API_BASE_URL,required"`
	AllowedCORSDomains []string `env:"API_ALLOWED_CORS_DOMAINS"`
}

type GinConfig struct {
	Mode string `env:"GIN_MODE,required"`
}
