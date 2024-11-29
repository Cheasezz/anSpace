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

// @securityDefinitions.apikey AccessToken
// @in header
// @name Authorization

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
