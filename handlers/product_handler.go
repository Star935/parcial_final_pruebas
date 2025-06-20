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

// Recupera todos los productos
func (h *Handler) GetProducts(c echo.Context) error {
	// Valida la conexion a la coleccion
	if h.Products == nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"status"  : http.StatusNotFound, 
			"message" : "Sin conexion a la colección",
			"data"	  : nil,
		})
	}

	// Recupera todos los productos
	cur, err := h.Products.Find(context.Background(), bson.M{})
	// Valuda si recupero los productos
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"status" : http.StatusInternalServerError,
			"message": err.Error(),
			"data"   : nil,
		})
	}

	// Lista de productos
	var products []models.Product

	// Almacena en la lista de productos todos los productos recuperados y valida si la operacion es exitosa
	if err := cur.All(context.Background(), &products); err != nil{
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"status"  : http.StatusInternalServerError,
			"message" : err.Error(),
			"data"    : nil,
		})
	}

	// Valida si se recupero algun producto
	if len(products) == 0 {
		return c.JSON(http.StatusNotFound, echo.Map{
			"status"  : http.StatusNotFound,
			"message" : "No se encontraron productos",
			"data"	  : nil,
		})
	}

	// Retorna estado de respuesta ok y todos los productos recuperados
	return c.JSON(http.StatusFound, echo.Map{
		"status"  : http.StatusFound,
		"message" : "Lista de productos encontrada",
		"data"    : products,
	})
}

// Recupera un producto mediante su id
func (h *Handler) GetProductById(c echo.Context) error {
	// Valida la conexion a la coleccion
	if h.Products == nil {
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

	// Instancia de Product
	var product models.Product

	// Recupera el producto mediante su id y lo decodifica en el espacio de memoria de la instancia de Product
	err = h.Products.FindOne(context.Background(), bson.M{"_id": id}).Decode(&product)
	// Valuda si no existe el documento
	if err == mongo.ErrNoDocuments {
		return c.JSON(http.StatusNotFound, echo.Map{
			"status"  : http.StatusNotFound,
			"message" : "Producto no encontrado",
			"data"    : nil,
		})
	} else if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"status"  : http.StatusInternalServerError,
			"message" : err.Error(),
			"data"    : nil,
		})
	}

	// Retorna estado de respuesta ok y todos los productos recuperados
	return c.JSON(http.StatusFound, echo.Map{
		"status"  : http.StatusFound,
		"message" : "Producto encontrado",
		"data"    : product,
	})
}

// Crea un nuevo producto
func (h *Handler) CreateProduct(c echo.Context) error {

	// Instancia de Product
	var product models.Product

	// Ingresa los datos de la request a la memoria de la instancia de Product
	if err := c.Bind(&product); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"status"  : http.StatusBadRequest, 
			"message" : "Input invalido",
			"data"    : nil,
		})
	}

	// Valida el nombre
	if strings.TrimSpace(product.Name) == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"status"  : http.StatusBadRequest, 
			"message" : "El nombre es obligatorio",
			"data"    : nil,
		})
	}

	// Valida el precio
	if product.Price <= 0 {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"status"  : http.StatusBadRequest, 
			"message" : "El precio es obligatorio y debe ser valido",
			"data"    : nil,
		})
	}

	// Valida el stock
	if product.Stock <= 0 {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"status"  : http.StatusBadRequest, 
			"message" : "El stock es obligatorio y debe ser valido",
			"data"    : nil,
		})
	}

	// Valida el id de la categoria
	if product.CategoryId <= 0 {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"status"  : http.StatusBadRequest, 
			"message" : "El category_id es obligatorio y debe ser valido",
			"data"    : nil,
		})
	}

	// Realiza la insercion en la base de datos
	_, err := h.Products.InsertOne(context.Background(), product)
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
		"message" : "Producto creado exitosamente",
		"data"	  : product,
	})
}

// Actualiza un producto existente
func (h *Handler) UpdateProduct(c echo.Context) error {
	// Valida la conexion a la coleccion
	if h.Products == nil {
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

	// Instancia de la estructura de Product
	var product models.Product

	// Ingresa los datos de la request a la memoria de la instancia de Product
	if err := c.Bind(&product); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"status"  : http.StatusBadRequest, 
			"message" : "Input invalido",
			"data"	  : nil,
		})
	}

	// Valida que el nombre no este vacia
	if strings.TrimSpace(product.Name) == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"status"  : http.StatusBadRequest, 
			"message" : "El nombre es obligatorio",
			"data"    : nil,
		})
	}

	// Valida el precio
	if product.Price <= 0 {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"status"  : http.StatusBadRequest, 
			"message" : "El precio es obligatorio y debe ser valido",
			"data"    : nil,
		})
	}

	// Valida el stock
	if product.Stock <= 0 {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"status"  : http.StatusBadRequest, 
			"message" : "El stock es obligatorio y debe ser valido",
			"data"    : nil,
		})
	}

	// Valida el id de la categoria
	if product.CategoryId <= 0 {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"status"  : http.StatusBadRequest, 
			"message" : "El category_id es obligatorio y debe ser valido",
			"data"    : nil,
		})
	}

	// Prepara filtro y documento de actualización
    filter := bson.M{"_id": id}
    update := bson.M{
        "$set": bson.M{
            "name" 		  : product.Name,
            "price" 	  : product.Price,
            "stock" 	  : product.Stock,
            "category_id" : product.CategoryId,
        },
    }

	// Actualiza el documento
    res, err := h.Products.UpdateOne(context.Background(), filter, update)
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
            "message" : "Producto no encontrado",
            "data"    : nil,
        })
    }

	// Retorna respuesta 201
	return c.JSON(http.StatusCreated, echo.Map{
        "status"  : http.StatusCreated,
        "message" : "Producto actualizado exitosamente",
        "data"    : product,
    })
}

// Elimina un producto
func (h *Handler) DeleteProduct(c echo.Context) error {
	// Valida la conexion a la coleccion
	if h.Products == nil {
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
	res, err := h.Products.DeleteOne(context.Background(), bson.M{"_id": id})
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
			"message" : "Producto no encontrado",
			"data" 	  : nil,
		})
	}

	// Retorna respuesta 200
	return c.JSON(http.StatusOK, echo.Map{
        "status"  : http.StatusOK,
        "message" : "Producto eliminado exitosamente",
        "data"    : nil,
    })
}