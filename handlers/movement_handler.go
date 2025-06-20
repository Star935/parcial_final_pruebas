package handlers

import (
	"context"
	"net/http"
	"time"

	"parcial_final/models"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Recupera todos los movimientos
func (h *Handler) GetMovements(c echo.Context) error {
	// Valida la conexion a la coleccion
	if h.Categories == nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"status"  : http.StatusNotFound, 
			"message" : "Sin conexion a la colección",
			"data"	  : nil,
		})
	}

	// Recupera todos los movimientos
	cur, err := h.Categories.Find(context.Background(), bson.M{})
	// Valuda si recupero los movimientos
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"status" : http.StatusInternalServerError,
			"message": err.Error(),
			"data"   : nil,
		})
	}

	// Lista de movimientos
	var movements []models.Movement

	// Almacena en la lista de movimientos todos los movimientos recuperados y valida si la operacion es exitosa
	if err := cur.All(context.Background(), &movements); err != nil{
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"status"  : http.StatusInternalServerError,
			"message" : err.Error(),
			"data"    : nil,
		})
	}

	// Valida si se recupero algun movemento
	if len(movements) == 0 {
		return c.JSON(http.StatusNotFound, echo.Map{
			"status"  : http.StatusNotFound,
			"message" : "No se encontraron movimientos",
			"data"	  : nil,
		})
	}

	// Retorna estado de respuesta ok y todos los movimientos recuperados
	return c.JSON(http.StatusFound, echo.Map{
		"status"  : http.StatusFound,
		"message" : "Lista de movimientos encontrada",
		"data"    : movements,
	})
}

// Registra nuevo movimiento
func (h *Handler) RegisterMovement(c echo.Context) error {

	// Instancia de Movement
	var movement models.Movement

	// Ingresa los datos de la request a la memoria de la instancia de Movement
	if err := c.Bind(&movement); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"status"  : http.StatusBadRequest, 
			"message" : "Input invalido",
			"data"    : nil,
		})
	}

	// Valida el nombre
	if movement.Stock <= 0 {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"status"  : http.StatusBadRequest, 
			"message" : "El stock es obligatorio y debe ser valido",
			"data"    : nil,
		})
	}

	// Setea el tiempo de entrada
	movement.Entrance = time.Now()
	movement.Exit     = false

	// Realiza la insercion en la base de datos
	_, err := h.Categories.InsertOne(context.Background(), movement)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"status"  : http.StatusInternalServerError, 
			"message" : err.Error(),
			"data"	  : nil,
		})
	}

	// Retorna json de respuesta
	return c.JSON(http.StatusCreated, echo.Map{
		"status"  : http.StatusCreated, 
		"message" : "Movimiento registrado exitosamente",
		"data"	  : movement,
	})
}

// Registra salida
func (h *Handler) ExitMovement(c echo.Context) error {
	// Valida la conexion a la coleccion
	if h.Categories == nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"status"  : http.StatusNotFound, 
			"message" : "Sin conexion a la colección",
			"data"	  : nil,
		})
	}

	// Recupera el parametro de consulta el id
	idParam := c.Param("id")

	// Convierte el id en ObjectID
	id, err := primitive.ObjectIDFromHex(idParam)
	// Valida si la operacion fue exitosa
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"status"  : http.StatusBadRequest, 
			"message" : "Id invalido",
			"data"	  : nil,
		})
	}

	// Instancia de la estructura de Movement
	var movement models.Movement

	// Recupera el movemento mediante su id y lo decodifica en el espacio de memoria de la instancia de Movement
	err = h.Categories.FindOne(context.Background(), bson.M{"_id": id}).Decode(&movement)
	// Valuda si no existe el documento
	if err == mongo.ErrNoDocuments {
		return c.JSON(http.StatusNotFound, echo.Map{
			"status"  : http.StatusNotFound,
			"message" : "Movimiento no encontrado",
			"data"    : nil,
		})
	} else if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"status"  : http.StatusInternalServerError,
			"message" : err.Error(),
			"data"    : nil,
		})
	}

	// Registra la salida
	movement.Exit = true;

	// Prepara filtro y documento de actualización
    filter := bson.M{"_id": id}
    update := bson.M{
        "$set": bson.M{
            "stock"    : movement.Stock,
            "entrance" : movement.Entrance,
            "exit" 	   : movement.Exit,
        },
    }

	// Actualiza el documento
    res, err := h.Categories.UpdateOne(context.Background(), filter, update)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, echo.Map{
            "status"  : http.StatusInternalServerError,
            "message" : err.Error(),
            "data"    : nil,
        })
    }

	// Valida si se afecto algun documento
    if res.MatchedCount == 0 {
        return c.JSON(http.StatusNotFound, echo.Map{
            "status"  : http.StatusNotFound,
            "message" : "Categoria no encontrado",
            "data"    : nil,
        })
    }

	// Retorna respuesta 201
	return c.JSON(http.StatusCreated, echo.Map{
        "status"  : http.StatusCreated,
        "message" : "Categoria actualizada exitosamente",
        "data"    : movement,
    })
}