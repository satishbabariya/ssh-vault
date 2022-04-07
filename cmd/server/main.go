package main

import (
	"net/http"

	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	// auth0 "github.com/satishbabariya/go-echo-auth0-middleware"
)

func main() {
	e := echo.New()
	e.HideBanner = true

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// e.Use(auth0.Auth0WithConfig(auth0.Auth0Config{
	// 	Issuer:   "https://<your tenant domain>/",
	// 	Audience: []string{"<your api identifier>"},
	// }))

	e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Root:   "/public",
		Index:  "index.html",
		Browse: true,
		HTML5:  true,
	}))

	// Routes
	e.GET("/api/validate", func(c echo.Context) error {
		claims := c.Get("claims").(*validator.ValidatedClaims)
		return c.JSON(http.StatusOK, claims)
	})

	// Start server
	e.Logger.Fatal(e.Start(":3000"))
}
