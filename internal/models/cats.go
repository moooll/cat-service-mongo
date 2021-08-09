package models

import "github.com/google/uuid"

type Cat struct {
	ID    uuid.UUID `bson:"id" json:"id1"`
	Name  string    `bson:"name" json:"name"`
	Breed string    `bson:"breed" json:"breed"`
	Color string    `bson:"color" json:"color"`
	Age   float32   `bson:"age" json:"age"`
	Price float32   `bson:"price" json:"price"`
}
