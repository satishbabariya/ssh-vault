package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"ssh-vault/internal/middleware/auth0"
	"ssh-vault/internal/store"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load()

	if os.Getenv("AUTH0_ISSUER") == "" {
		log.Fatal("AUTH0_ISSUER is not set")
	}

	if os.Getenv("AUTH0_AUDIENCE") == "" {
		log.Fatal("AUTH0_AUDIENCE is not set")
	}

	if os.Getenv("VAULT_SECRET") == "" {
		log.Fatal("VAULT_SECRET is not set")
	}

	if os.Getenv("DATABASE_URL") == "" {
		log.Fatal("DATABASE_URL is not set")
	}
}

func main() {
	store, err := store.NewStore(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer store.Close()

	err = store.Init(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			// Status code defaults to 500
			code := fiber.StatusInternalServerError

			// Retrieve the custom status code if it's an fiber.*Error
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}

			return c.Status(code).JSON(err)
		},
	})

	// Default middleware config
	app.Use(logger.New())

	app.Static("/", "/public")
	// app.All("*", func(c *fiber.Ctx) error {
	// 	return c.SendFile("/public/index.html")
	// })

	api := app.Group("/api")
	api.Use(auth0.New(auth0.Config{
		Issuer:   os.Getenv("AUTH0_ISSUER"),
		Audience: []string{os.Getenv("AUTH0_AUDIENCE")},
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return fiber.NewError(http.StatusUnauthorized, err.Error())
		},
	}))

	log.Fatal(app.Listen(":1203"))
}
