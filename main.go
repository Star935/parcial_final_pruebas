package main

import (
	"context"
	"log"

	"parcial_final/handlers"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Punto de entrada del programa
func main() {
	// Instancia de Echo
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Conexion a la base de datos
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	if err := client.Ping(context.Background(), nil); err != nil {
		log.Fatal(err)
	}

	// Define la base de datos
	db := client.Database("parcial_final")

	// Crea las colecciones
	h := handlers.NewHandler(db.Collection("products"), db.Collection("categories"), db.Collection("movements"))

	// Utiliza el html en la raiz
	e.GET("/", func(c echo.Context) error {
		return c.File("index.html")
	})

	// Rutas para la gestion de productos
	e.GET("/products", h.GetProducts)
	e.GET("/products/:id", h.GetProductById)
	e.POST("/products", h.CreateProduct)
	e.PUT("/products/:id", h.UpdateProduct)
	e.DELETE("/products/:id", h.DeleteProduct)

	// Rutas para la gestion de categorias
	e.GET("/categories", h.GetCategories)
	e.GET("/categories/:id", h.GetCategoryById)
	e.POST("/categories", h.CreateCategory)
	e.DELETE("/categories/:id", h.DeleteCategory)

	// Rutas para la gestion de movimientos
	e.GET("/movements", h.GetMovements)
	e.GET("/register-movement", h.RegisterMovement)
	e.DELETE("/exit-movement/:id", h.ExitMovement)

	e.Logger.Fatal(e.Start(":8080"))
}