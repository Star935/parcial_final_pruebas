package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Estructura de producto
type Product struct {
	ID    		primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name  		string 			   `json:"name" bson:"name"`
	Price  		int 			   `json:"price" bson:"price"`
	Stock  		int 			   `json:"stock" bson:"stock"`
	CategoryId  int 			   `json:"category_id" bson:"category_id"`
}