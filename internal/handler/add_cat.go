package handler

import (
	"cat-service/internal/models"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (s *Service) AddCat(c echo.Context) error {
	cat := &models.Cat{}
	if err := (&echo.DefaultBinder{}).BindBody(c, &cat); err != nil {
		log.Print("bind body ", err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	err := s.catalog.Save(*cat)
	if err != nil {
		log.Print("save ", err)
		return c.JSON(http.StatusInternalServerError, "error saving cat(")
	}

	return c.JSON(http.StatusCreated, "created")
}
