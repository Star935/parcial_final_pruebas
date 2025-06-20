package handlers

import "go.mongodb.org/mongo-driver/mongo"

// Estructura de las colecciones de mongo
type Handler struct {
	Products   *mongo.Collection
	Categories *mongo.Collection
	Movements     *mongo.Collection
}

// Constructor de los handlers
func NewHandler(products, categories, movements *mongo.Collection) *Handler {
	return &Handler{
		Products   : products, 
		Categories : categories,
		Movements  : movements,
	}
}