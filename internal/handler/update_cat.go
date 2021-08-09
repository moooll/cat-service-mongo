package handler

import (
	"cat-service/internal/models"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (s *Service) UpdateCat(c echo.Context) error {
	req := &models.Cat{}
	if err := (&echo.DefaultBinder{}).BindBody(c, req); err != nil {
		log.Print("bind body: ", err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	updated, err := s.catalog.Update(*req)
	if err != nil {
		log.Print("update: ", err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, updated)
}
