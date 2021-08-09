package main

import (
	"cat-service/internal/handler"
	"cat-service/internal/repository"
	"context"

	"log"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
)
 
func main() {
	client, err := repository.Connect()
	if err != nil {
		log.Print("could not connect to the db\n", err.Error())
	}

	defer repository.Close(client)

	collection := client.Database("catalog").Collection("cats2")
	service := handler.NewService(repository.NewCatalog(collection))
	dbs, err := client.ListDatabases(context.Background(), bson.M{})
	if err != nil {
		log.Print("error listing dbs ", err.Error())
	}
	collections, _ := client.Database("catalog").ListCollectionNames(context.Background(), bson.M{})
	log.Print(collections)

	log.Print(dbs)
	e := echo.New()
	e.POST("/cats", service.AddCat)
	e.GET("/cats", service.GetAllCats)
	e.GET("/cats/:id", service.GetCat)
	e.PUT("/cats", service.UpdateCat)
	e.DELETE("/cats/:id", service.DeleteCat)
	e.GET("/cats/get-rand-cat", handler.GetRandCat)
	if err := e.Start(":8081"); err != nil {
		log.Print("could not start server\n", err.Error())
	}
}
