package main

import (
	"fmt"
	"net/http"
	"os"
	"ssh-vault/internal/model"
	"ssh-vault/internal/store"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
	// auth0 "github.com/satishbabariya/go-echo-auth0-middleware"
)

func init() {
	godotenv.Load()

	if os.Getenv("AUTH0_ISSUER") == "" {
		logrus.Fatal("AUTH0_ISSUER is not set")
	}

	if os.Getenv("AUTH0_AUDIENCE") == "" {
		logrus.Fatal("AUTH0_AUDIENCE is not set")
	}

	if os.Getenv("VAULT_SECRET") == "" {
		logrus.Fatal("VAULT_SECRET is not set")
	}
}

func main() {
	e := echo.New()
	e.HideBanner = true

	e.Use(middleware.Recover())

	// custom validator
	e.Validator = model.NewValidator()

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${time_rfc3339} | ${status} | ${method} | ${latency_human} | ${uri}\n",
		Output: logrus.New().Out,
	}))

	// e.Use(auth0.Auth0WithConfig(auth0.Auth0Config{
	// 	Issuer:   os.Getenv("AUTH0_ISSUER"),
	// 	Audience: []string{os.Getenv("AUTH0_AUDIENCE")},
	// }))

	e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Root:   "/public",
		Index:  "index.html",
		Browse: false,
		HTML5:  true,
	}))

	store, err := store.Open("./vault.db")
	if err != nil {
		e.Logger.Fatal(err)
	}
	defer store.Close()

	// Routes
	// e.GET("/api/keys", func(c echo.Context) error {
	// 	claims := c.Get("claims").(*validator.ValidatedClaims)
	// 	return c.JSON(http.StatusOK, claims)
	// })

	e.GET("/api/credentials/:host", func(c echo.Context) error {
		credential, err := store.Get(c.Param("host"))
		if err != nil {
			return c.JSON(http.StatusNotFound, err)
		}
		return c.JSON(http.StatusOK, credential)
	})

	e.POST("/api/credentials", func(c echo.Context) error {
		cred := new(model.Credential)
		if err = c.Bind(cred); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		if err = c.Validate(cred); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		if err = store.Add(*cred); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusCreated, fmt.Sprintf("Added %s", cred.Host))
	})

	// Start server
	e.Logger.Fatal(e.Start(":3000"))
}
