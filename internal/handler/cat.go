// Package handler contains handlers for http requests
package handler

import (
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/moooll/cat-service-mongo/internal/models"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

// AddCat endpoint receives new cat in request body and puts it in the database, if succeeds, sends OK, "created"
func (s *Service) AddCat(c echo.Context) error {
	cat := &models.Cat{}
	if err := (&echo.DefaultBinder{}).BindBody(c, &cat); err != nil {
		log.Errorln("add cat bind body ", err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	err := s.storage.SaveToStorage(c.Request().Context(), *cat)
	if err != nil {
		log.Errorln("save to storage: ", err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusCreated, "created")
}

// DeleteCat endpoints receives id as a url param, and deletes the document with the corresponding id
func (s *Service) DeleteCat(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		log.Errorln(err.Error())
		return c.JSON(http.StatusInternalServerError, err)
	}

	err = s.storage.DeleteFromStorage(c.Request().Context(), id)
	if err != nil {
		log.Errorln(err.Error())
		return c.JSON(http.StatusInternalServerError, err)
	}

	data := map[string]interface{}{"act": "delete", "message": fmt.Sprintf("cat id#%s was deleted", id.String())}
	err = s.stream.Push(c.Request().Context(), data)
	if err != nil {
		log.Errorln(err.Error())
	}

	return c.JSON(http.StatusOK, "deleted")
}

// GetAllCats endpoint sends all cats
func (s *Service) GetAllCats(c echo.Context) error {
	var cats []models.Cat
	cats, err := s.storage.GetAllFromStorage(c.Request().Context())
	if err != nil && err != mongo.ErrNoDocuments {
		log.Errorln(err.Error())
		return c.JSON(http.StatusInternalServerError, err)
	} else if err == mongo.ErrNoDocuments {
		return c.String(200, "nothing found")
	}

	return c.JSON(http.StatusOK, cats)
}

// GetCat sends cat by id from url params
func (s *Service) GetCat(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil && err != mongo.ErrNoDocuments {
		log.Errorln(err.Error())
		return c.JSON(http.StatusInternalServerError, err)
	}

	cat, err := s.storage.GetFromStorage(c.Request().Context(), id)
	if err != nil && err != mongo.ErrNoDocuments {
		log.Errorln(err.Error())
		return c.JSON(http.StatusInternalServerError, err)
	} else if err == mongo.ErrNoDocuments {
		return c.String(200, "nothing found")
	}

	return c.JSON(http.StatusOK, cat)
}

// UpdateCat endpoint updates cat specified in body
func (s *Service) UpdateCat(c echo.Context) error {
	cat := models.Cat{}
	if err := (&echo.DefaultBinder{}).BindBody(c, &cat); err != nil {
		log.Errorln("bind body: ", err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	err := s.storage.UpdateStorage(c.Request().Context(), cat)
	if err != nil {
		log.Errorln("update: ", err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, "updated")
}
