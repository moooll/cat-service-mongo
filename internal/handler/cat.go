// Package handler contains handlers for http requests
package handler

import (
	"log"
	"net/http"

	"github.com/moooll/cat-service-mongo/internal/models"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// AddCat endpoint receives new cat in request body and puts it in the database, if succeeds, sends OK, "created"
func (s *Service) AddCat(c echo.Context) error {
	cat := &models.Cat{}
	if err := (&echo.DefaultBinder{}).BindBody(c, &cat); err != nil {
		log.Println("bind body ", err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	err := s.cache.SetToHash(*cat)
	if err != nil {
		log.Println("save to cache error ", err)
		return c.JSON(http.StatusInternalServerError, "error saving cat(")
	}

	err = s.catalog.Save(*cat)
	if err != nil {
		log.Println("save ", err)
		return c.JSON(http.StatusInternalServerError, "error saving cat(")
	}

	return c.JSON(http.StatusCreated, "created")
}

// DeleteCat endpoints receives id as a url param, and deletes the document with the corresponding id
func (s *Service) DeleteCat(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		log.Println(err.Error())
		return c.JSON(http.StatusInternalServerError, err)
	}

	err = s.cache.DeleteFromCache(id)
	if err != nil {
		log.Println("delete from cache error ", err)
		return c.JSON(http.StatusInternalServerError, err)
	}
	_, err = s.catalog.Delete(id)
	if err != nil {
		log.Println(err.Error())
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, "deleted")
}

// GetAllCats endpoint sends all cats
func (s *Service) GetAllCats(c echo.Context) error {
	var cats []models.Cat
	cats, err := s.cache.GetAllFromCache()
	if err != nil || len(cats) == 0 {
		// cats, err = s.catalog.GetAll()
		// if err != nil {
		// 	log.Println(err.Error())
		// 	return c.JSON(http.StatusInternalServerError, err)
		// }

		err = s.cache.SetAllToHash(cats)
		if err != nil {
			log.Println("save to cache error ", err)
			return c.JSON(http.StatusInternalServerError, err)
		}
	}
	

	return c.JSON(http.StatusOK, cats)
}

// GetCat sends cat by id from url params
func (s *Service) GetCat(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		log.Println(err.Error())
		return c.JSON(http.StatusInternalServerError, err)
	}

	cat, err := s.cache.GetFromCache(id)
	if err != nil || (models.Cat{}) == cat {
		// cat, err = s.catalog.Get(id)
		// if err != nil {
		// 	log.Println(err.Error())
		// 	return c.JSON(http.StatusInternalServerError, err)
		// }

		err = s.cache.SetToHash(cat)
		if err != nil {
			log.Println("save to cache error ", err)
			return c.JSON(http.StatusInternalServerError, err)
		}
	}

	return c.JSON(http.StatusOK, cat)
}

// UpdateCat endpoint updates cat specified in body
func (s *Service) UpdateCat(c echo.Context) error {
	req := &models.Cat{}
	if err := (&echo.DefaultBinder{}).BindBody(c, req); err != nil {
		log.Print("bind body: ", err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	err := s.cache.SetToHash(*req)
	if err != nil {
		log.Println("save to cache error ", err)
		return c.JSON(http.StatusInternalServerError, "error saving to cache(")
	}

	err = s.catalog.Update(*req)
	if err != nil {
		log.Print("update: ", err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, "updated")
}
