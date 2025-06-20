package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Estructura de categoria
type Category struct {
	ID    		primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name  		string 			   `json:"name" bson:"name"`
	Description string 			   `json:"description" bson:"description"`
}