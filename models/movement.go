package models

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Estructura de Movements
type Movement struct {
	ID    		primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Stock 		int				   `json:"stock" bson:"stock"`
	Entrance    time.Time   	   `bson:"entrance"`
	Exit  		bool		       `bson:"exit"`
}