package main

import (
	"log"

	"github.com/Cheasezz/anSpace/backend/config"
	"github.com/Cheasezz/anSpace/backend/internal/app"
	"github.com/joho/godotenv"
)

// @title AnSpace App API
// @version 1.0
// @description API Server for AnSpace Application

// @host localhost:8000
// @BasePath /

// @securityDefinitions.apikey bearerAuth
// @in header
// @name Authorization
// @description Enter the token with the `Bearer: ` prefix

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("error loading env variables: %s", err.Error())
	}

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	app.Run(cfg)
}
