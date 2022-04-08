package main

import (
	"log"
	"net/http"
	"os"
	"ssh-vault/internal/model"
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
}

func main() {
	store, err := store.Open("./vault.db")
	if err != nil {
		log.Fatal(err)
	}
	defer store.Close()

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

	app.Get("/api/credentials/:host", func(c *fiber.Ctx) error {
		credential, err := store.Get(c.Params("host"))
		if err != nil {
			return fiber.NewError(http.StatusNotFound, err.Error())
		}
		if credential == nil || credential.Host == "" {
			return fiber.NewError(http.StatusNotFound, "Credential not found")
		}
		return c.Status(http.StatusOK).JSON(credential)
	})

	app.Post("/api/credentials", func(c *fiber.Ctx) error {
		cred := new(model.Credential)
		if err := c.BodyParser(cred); err != nil {
			return fiber.NewError(http.StatusBadRequest, err.Error())
		}
		if cred.Host == "" {
			return fiber.NewError(http.StatusBadRequest, "Host is required")
		}
		if cred.User == "" {
			return fiber.NewError(http.StatusBadRequest, "User is required")
		}

		if err := store.Add(*cred); err != nil {
			return fiber.NewError(http.StatusInternalServerError, err.Error())
		}

		type response struct {
			Host string `json:"host"`
			Port int    `json:"port"`
		}

		return c.Status(http.StatusCreated).JSON(response{
			Host: cred.Host,
			Port: cred.Port,
		})
	})

	app.All("*", func(c *fiber.Ctx) error {
		return c.SendFile("/public/index.html")
	})

	log.Fatal(app.Listen(":3000"))
}
