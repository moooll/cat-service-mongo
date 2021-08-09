package handler

import (
	"cat-service/internal/repository"

	"github.com/labstack/echo/v4"
)

func GetRandCat(c echo.Context) error {
	cat := repository.RandCat()
	return c.JSON(200, cat)
}
