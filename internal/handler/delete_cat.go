package handler

import (
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func (s *Service) DeleteCat(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		log.Print(err.Error())
		return c.JSON(http.StatusInternalServerError, err)
	}

	deleted, err := s.catalog.Delete(id)
	if err != nil {
		log.Print(err.Error())
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, deleted)

}
