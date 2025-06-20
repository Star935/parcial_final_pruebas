package handlers

import (
	"context"
	"net/http"
	"strings"

	"parcial_final/models"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Recupera todos los categoryos
func (h *Handler) GetCategories(c echo.Context) error {
	// Valida la conexion a la coleccion
	if h.Categories == nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"status"  : http.StatusNotFound, 
			"message" : "Sin conexion a la colección",
			"data"	  : nil,
		})
	}

	// Recupera todos los categoryos
	cur, err := h.Categories.Find(context.Background(), bson.M{})
	// Valuda si recupero los categoryos
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"status" : http.StatusInternalServerError,
			"message": err.Error(),
			"data"   : nil,
		})
	}

	// Lista de categoryos
	var categorys []models.Category

	// Almacena en la lista de categoryos todos los categoryos recuperados y valida si la operacion es exitosa
	if err := cur.All(context.Background(), &categorys); err != nil{
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"status"  : http.StatusInternalServerError,
			"message" : err.Error(),
			"data"    : nil,
		})
	}

	// Valida si se recupero algun categoryo
	if len(categorys) == 0 {
		return c.JSON(http.StatusNotFound, echo.Map{
			"status"  : http.StatusNotFound,
			"message" : "No se encontraron categoryos",
			"data"	  : nil,
		})
	}

	// Retorna estado de respuesta ok y todos los categoryos recuperados
	return c.JSON(http.StatusFound, echo.Map{
		"status"  : http.StatusFound,
		"message" : "Lista de categoryos encontrada",
		"data"    : categorys,
	})
}

// Recupera un categoryo mediante su id
func (h *Handler) GetCategoryById(c echo.Context) error {
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
			"data"    : nil,
		})
	}

	// Instancia de Category
	var category models.Category

	// Recupera el categoryo mediante su id y lo decodifica en el espacio de memoria de la instancia de Category
	err = h.Categories.FindOne(context.Background(), bson.M{"_id": id}).Decode(&category)
	// Valuda si no existe el documento
	if err == mongo.ErrNoDocuments {
		return c.JSON(http.StatusNotFound, echo.Map{
			"status"  : http.StatusNotFound,
			"message" : "Categoria no encontrado",
			"data"    : nil,
		})
	} else if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"status"  : http.StatusInternalServerError,
			"message" : err.Error(),
			"data"    : nil,
		})
	}

	// Retorna estado de respuesta ok y todos los categoryos recuperados
	return c.JSON(http.StatusFound, echo.Map{
		"status"  : http.StatusFound,
		"message" : "Categoria encontrado",
		"data"    : category,
	})
}

// Crea un nuevo categoria
func (h *Handler) CreateCategory(c echo.Context) error {

	// Instancia de Category
	var category models.Category

	// Ingresa los datos de la request a la memoria de la instancia de Category
	if err := c.Bind(&category); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"status"  : http.StatusBadRequest, 
			"message" : "Input invalido",
			"data"    : nil,
		})
	}

	// Valida el nombre
	if strings.TrimSpace(category.Name) == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"status"  : http.StatusBadRequest, 
			"message" : "El nombre es obligatorio",
			"data"    : nil,
		})
	}

	// Valida la descripcion
	if strings.TrimSpace(category.Description) == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"status"  : http.StatusBadRequest, 
			"message" : "La descripcion es obligatoria",
			"data"    : nil,
		})
	}

	// Realiza la insercion en la base de datos
	_, err := h.Categories.InsertOne(context.Background(), category)
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
		"message" : "Categoria creado exitosamente",
		"data"	  : category,
	})
}

// Actualiza un categoryo existente
func (h *Handler) UpdateCategory(c echo.Context) error {
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

	// Instancia de la estructura de Category
	var category models.Category

	// Ingresa los datos de la request a la memoria de la instancia de Category
	if err := c.Bind(&category); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"status"  : http.StatusBadRequest, 
			"message" : "Input invalido",
			"data"	  : nil,
		})
	}

	// Valida que el nombre no este vacia
	if strings.TrimSpace(category.Name) == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"status"  : http.StatusBadRequest, 
			"message" : "El nombre es obligatorio",
			"data"    : nil,
		})
	}

	// Valida la descripcion
	if strings.TrimSpace(category.Description) == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"status"  : http.StatusBadRequest, 
			"message" : "La descripcion es obligatoria",
			"data"    : nil,
		})
	}

	// Prepara filtro y documento de actualización
    filter := bson.M{"_id": id}
    update := bson.M{
        "$set": bson.M{
            "name" 		  : category.Name,
            "description" : category.Description,
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
        "data"    : category,
    })
}

// Elimina un categoryo
func (h *Handler) DeleteCategory(c echo.Context) error {
	// Valida la conexion a la coleccion
	if h.Categories == nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"status": http.StatusNotFound, 
			"message": "Sin conexion a la colección",
			"data": nil,
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
			"data" 	  : nil,
		})
	}

	// Realiza operacion de eliminado mediante el id recuperado del parametro de consulta
	res, err := h.Categories.DeleteOne(context.Background(), bson.M{"_id": id})
	// Valida si la operacion fue exitosa
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"status"  : http.StatusInternalServerError, 
			"message" : err.Error(),
			"data" 	  : nil,
		})
	}
	
	// Valida si se elimino algun documento
	if res.DeletedCount == 0 {
		return c.JSON(http.StatusNotFound, echo.Map{
			"status"  : http.StatusNotFound, 
			"message" : "Categoria no encontrado",
			"data" 	  : nil,
		})
	}

	// Retorna respuesta 200
	return c.JSON(http.StatusOK, echo.Map{
        "status"  : http.StatusOK,
        "message" : "Categoria eliminado exitosamente",
        "data"    : nil,
    })
}