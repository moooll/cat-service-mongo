// Package models provides models for the database
package models

import "github.com/google/uuid"

// Cat describes cat entity for http interactions
type Cat struct {
	ID    uuid.UUID `bson:"id" json:"id1"`
	Name  string    `bson:"name" json:"name"`
	Breed string    `bson:"breed" json:"breed"`
	Color string    `bson:"color" json:"color"`
	Age   float32   `bson:"age" json:"age"`
	Price float32   `bson:"price" json:"price"`
}
