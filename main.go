package main

import (
	"errors"
	"fmt"
	"github.com/Enthys/url_shortener/http/controller"
	"github.com/Enthys/url_shortener/pkg/repository"
	"github.com/Enthys/url_shortener/pkg/repository/redis"
	"github.com/Enthys/url_shortener/pkg/services"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"os"
)

func createLinkRepository() (repository.LinkRepository, error) {
	switch os.Getenv("DATABASE") {
	case "redis":
		return redis.NewRedisRepository()

	default:
		if os.Getenv("DATABASE") == "" {
			return nil, errors.New(
				"environment variable `DATABASE` is not provided. Please provide the `DATABASE` environment variable",
			)
		}
		return nil, fmt.Errorf("unknown database type '%s'", os.Getenv("DATABASE"))
	}
}

func loadEnv() error {
	return godotenv.Load()
}

func main() {
	if err := loadEnv(); err != nil {
		panic(err)
	}

	linkRepository, err := createLinkRepository()
	if err != nil {
		panic(err)
	}

	linkService := services.NewLinkService(linkRepository)
	linkController := controller.NewLinkController(linkService)

	server := echo.New()
	server.Logger.SetLevel(log.DEBUG)
	linkController.SetupRoutes(server)
	server.Logger.Fatal(server.Start(":80"))
}
