package handler

import (
	"net/http"

	"github.com/moooll/cat-service-mongo/internal/repository"

	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

// GetRandCat endpoint returns random generated cat
func GetRandCat(c echo.Context) error {
	cat, err := repository.RandCat()
	if err != nil {
		log.Errorln(err.Error())
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(200, cat)
}
