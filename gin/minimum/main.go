package main

import (
	_ "github.com/joho/godotenv/autoload" // Autoload .env file.

	"github.com/yizeng/gab/gin/minimum/cmd/app"
)

func main() {
	app.Start()
}
