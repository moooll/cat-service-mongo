// Package repository contains utilities for integration with the database
package repository

import (
	"context"
	"log"

	"github.com/moooll/cat-service-mongo/internal/models"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoCatalog is used for db integration from handlers
type MongoCatalog struct {
	collection *mongo.Collection
}

// NewCatalog is used to create new *Catalog entity
func NewCatalog(coll *mongo.Collection) *MongoCatalog {
	return &MongoCatalog{
		collection: coll,
	}
}

// Save saves a doc to the database, returns the error
func (c *MongoCatalog) Save(ctx context.Context, cat models.Cat) error {
	_, err := c.collection.InsertOne(ctx, cat)
	if err != nil {
		return err
	}

	return nil
}

// Get retieves a doc by id from the db and returns the doc and the error
func (c *MongoCatalog) Get(ctx context.Context, id uuid.UUID) (cat models.Cat, err error) {
	err = c.collection.FindOne(ctx, bson.M{"id": id}).Decode(&cat)
	if err != nil {
		return models.Cat{}, err
	}

	return cat, nil
}

// GetAll retireves all docs from the db and returns models.Cat slice and error
func (c *MongoCatalog) GetAll(ctx context.Context) (cats []models.Cat, err error) {
	cur, err := c.collection.Find(ctx, bson.M{})
	if err != nil {
		return []models.Cat{}, err
	}

	cat := models.Cat{}
	for cur.Next(ctx) {
		err = cur.Decode(&cat)
		if err != nil {
			return []models.Cat{}, err
		}

		cats = append(cats, cat)
	}

	return cats, nil
}

// Delete deletes the doc by id from the database
func (c *MongoCatalog) Delete(ctx context.Context, id uuid.UUID) (deleted bson.M, err error) {
	err = c.collection.FindOneAndDelete(ctx, bson.M{"id": id}).Decode(&deleted)
	if err != nil {
		return bson.M{}, err
	}

	return deleted, nil
}

// Update updates the doc by id in the db, and if is not present, creates it
func (c *MongoCatalog) Update(ctx context.Context, cat models.Cat) error {
	opts := options.FindOneAndUpdate().SetUpsert(true)
	filter := bson.D{{Key: "id", Value: cat.ID}}
	upd := bson.D{primitive.E{Key: "$set", Value: bson.D{
		primitive.E{Key: "id", Value: cat.ID},
		primitive.E{Key: "name", Value: cat.Name},
		primitive.E{Key: "breed", Value: cat.Breed},
		primitive.E{Key: "color", Value: cat.Color},
		primitive.E{Key: "age", Value: cat.Age},
		primitive.E{Key: "price", Value: cat.Price},
	}}}
	log.Print(cat.ID)
	updated := models.Cat{}
	err := c.collection.FindOneAndUpdate(ctx, filter, upd, opts).Decode(&updated)
	if err != nil {
		return err
	}

	return nil
}
