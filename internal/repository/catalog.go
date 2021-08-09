package repository

import (
	"cat-service/internal/models"
	"context"
	"log"
	"math/rand"

	"github.com/google/uuid"
	"github.com/thanhpk/randstr"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Catalog struct {
	collection *mongo.Collection
}

func NewCatalog(coll *mongo.Collection) *Catalog {
	return &Catalog{
		collection: coll,
	}
}

func (c *Catalog) Save(cat models.Cat) error {
	_, err := c.collection.InsertOne(context.Background(), cat)
	if err != nil {
		return err
	}

	return nil
}

func (c *Catalog) Get(id uuid.UUID) (cat models.Cat, err error) {
	err = c.collection.FindOne(context.Background(), bson.M{"id": id}).Decode(&cat)
	if err != nil {
		return models.Cat{}, err
	}

	return cat, nil
}

func (c *Catalog) GetAll() (cats []models.Cat, err error) {
	cur, err := c.collection.Find(context.Background(), bson.M{})
	cat := models.Cat{}
	for cur.Next(context.Background()) {
		err = cur.Decode(&cat)
		if err != nil {
			return []models.Cat{}, err
		}

		cats = append(cats, cat)
	}

	return cats, nil
}

func (c *Catalog) Delete(id uuid.UUID) (deleted bson.M, err error) {
	err = c.collection.FindOneAndDelete(context.Background(), bson.M{"id": id}).Decode(&deleted)
	if err != nil {
		return bson.M{}, err
	}

	return deleted, nil
}

func (c *Catalog) Update(cat models.Cat) (models.Cat, error) {
	opts := options.FindOneAndUpdate().SetUpsert(true)
	filter := bson.D{{"id", cat.ID}}
	upd := bson.D{{"$set", bson.D{
		{"id", cat.ID},
		{"name", cat.Name},
		{"breed", cat.Breed},
		{"color", cat.Color},
		{"age", cat.Age},
		{"price", cat.Price},
	}}}
	log.Print(cat.ID)
	updated := models.Cat{}
	err := c.collection.FindOneAndUpdate(context.Background(), filter, upd, opts).Decode(&updated)
	if err != nil {
		return models.Cat{}, err
	}

	return updated, nil
}

func RandCat() models.Cat {
	id := uuid.New()
	name := randstr.String(8, "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	breed := randstr.String(8, "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	color := randstr.String(8, "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	age := rand.Float32() * 15
	price := rand.Float32() * 15
	return models.Cat{
		ID:    id,
		Name:  name,
		Breed: breed,
		Color: color,
		Age:   age,
		Price: price,
	}
}
