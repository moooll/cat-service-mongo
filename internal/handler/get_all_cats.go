package handler

import (
	"cat-service/internal/models"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (s *Service) GetAllCats(c echo.Context) error {
	var cats []models.Cat
	cats, err := s.catalog.GetAll()
	if err != nil {
		log.Print(err.Error())
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, cats)
}
