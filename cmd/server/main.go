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
	// Initialize the database.
	store, err := store.NewStore(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer store.Close()

	err = store.Connect(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	// Auto migrate the database. (Create tables if not exists)
	err = store.AutoMigrate(context.Background())
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
	// app.Use(cors.New())

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

	// api.Get("/credentials/:host", func(c *fiber.Ctx) error {
	// 	credential, err := store.Get(c.Params("host"))
	// 	if err != nil {
	// 		return fiber.NewError(http.StatusNotFound, err.Error())
	// 	}
	// 	if credential == nil || credential.Host == "" {
	// 		return fiber.NewError(http.StatusNotFound, "Credential not found")
	// 	}
	// 	return c.Status(http.StatusOK).JSON(credential)
	// })

	// api.Get("/credentials", func(c *fiber.Ctx) error {
	// 	credentials, err := store.Remotes()
	// 	if err != nil {
	// 		return fiber.NewError(http.StatusInternalServerError, err.Error())
	// 	}
	// 	return c.Status(http.StatusOK).JSON(credentials)
	// })

	// api.Post("/credentials", func(c *fiber.Ctx) error {
	// 	cred := new(model.Credential)
	// 	if err := c.BodyParser(cred); err != nil {
	// 		return fiber.NewError(http.StatusBadRequest, err.Error())
	// 	}
	// 	if cred.Host == "" {
	// 		return fiber.NewError(http.StatusBadRequest, "Host is required")
	// 	}
	// 	if cred.User == "" {
	// 		return fiber.NewError(http.StatusBadRequest, "User is required")
	// 	}

	// 	if err := store.Add(*cred); err != nil {
	// 		return fiber.NewError(http.StatusInternalServerError, err.Error())
	// 	}

	// 	type response struct {
	// 		Host string `json:"host"`
	// 		Port int    `json:"port"`
	// 	}

	// 	return c.Status(http.StatusCreated).JSON(response{
	// 		Host: cred.Host,
	// 		Port: cred.Port,
	// 	})
	// })

	log.Fatal(app.Listen(":1203"))
}
