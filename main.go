package main

import (
	"os"

	"github.com/Enthys/url_shortener/http/controller"
	"github.com/Enthys/url_shortener/pkg"
	"github.com/Enthys/url_shortener/pkg/repository"
	"github.com/Enthys/url_shortener/pkg/repository/redis"
	"github.com/Enthys/url_shortener/pkg/services"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

// createLinkRepository creates a `repository.LinkRepository` instance which works with the said `DATABASE` type.
//
// If the environment variable `DATABASE` is not set a `pkg.ErrorMissingDatabaseType` error will be returned.
//
// If the environment variable `DATABASE` is set but contains an unknown database type a `pkg.ErrorUnkownDatabaseType`
// error will be returned.
func createLinkRepository() (repository.LinkRepository, error) {
	switch os.Getenv("DATABASE") {
	case "redis":
		return redis.NewRedisRepository()

	default:
		if os.Getenv("DATABASE") == "" {
			return nil, pkg.ErrorMissingDatabaseType{}
		}

		return nil, pkg.ErrorUnkownDatabaseType{Type: os.Getenv("DATABASE")}
	}
}

func main() {
	// Load the .env if such exists
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	// Repository creation
	linkRepository, err := createLinkRepository()
	if err != nil {
		panic(err)
	}

	// Service setup
	linkService := services.NewLinkService(linkRepository)

	// Controller setup
	linkController := controller.NewLinkController(linkService)

	// Server setup
	server := echo.New()
	server.Logger.SetLevel(log.DEBUG)
	linkController.SetupRoutes(server)

	// Server start
	server.Logger.Fatal(server.Start(":80"))
}
