package model

import (
	"net/http"

	"github.com/go-playground/validator"
	"github.com/labstack/echo"
)

type Credential struct {
	Host      string  `json:"host" validate:"required"`
	Port      int     `json:"port,omitempty"`
	User      string  `json:"user" validate:"required"`
	PrivatKey *string `json:"private_key,omitempty"`
	Password  *string `json:"password,omitempty"`
}

type CustomValidator struct {
	validator *validator.Validate
}

func NewValidator() *CustomValidator {
	return &CustomValidator{validator: validator.New()}
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}
