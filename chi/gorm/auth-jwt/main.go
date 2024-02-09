package main

import (
	_ "github.com/joho/godotenv/autoload" // Autoload .env file.

	"github.com/yizeng/gab/chi/gorm/auth-jwt/cmd/app"
)

func main() {
	if err := app.Start(); err != nil {
		panic(err)
	}
}
