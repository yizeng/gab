package main

import (
	_ "github.com/joho/godotenv/autoload" // Autoload .env file.

	"github.com/yizeng/gab/chi/crud-gorm/cmd/app"
)

func main() {
	app.Start()
}
