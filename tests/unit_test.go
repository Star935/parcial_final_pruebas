package tests

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"

    "parcial_final/handlers"
    "parcial_final/models"
    "github.com/labstack/echo/v4"
)

// Prueba de validacion de creacion de categoria sin el nombre
func TestCreateCategoryValidationNameEmpty(t *testing.T) {
	// Instancia de echo
    e := echo.New()
    h := &handlers.Handler{Categories: nil}

	// Instancia de Category seteada sin Name
    book := models.Category{
		Name        : "", 
		Description : "Name empty",
	}

	// Codifica los datos de la Category
    body, _ := json.Marshal(book)

	// Define la request
    req := httptest.NewRequest(http.MethodPost, "/categories", bytes.NewReader(body))
    req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
    rec := httptest.NewRecorder()
    c := e.NewContext(req, rec)

	// Ejecuta el metodo para crear categoria
    if err := h.CreateCategory(c); err != nil {
        t.Fatal(err)
    }

	// Valida si el codigo de respuesta fue la esperada
    if rec.Code != http.StatusBadRequest {
        t.Errorf("Esperado 400, obtuvo %d", rec.Code)
    }
}

// Prueba de validacion de creacion de cateroria sin la descripcion
func TestCreateCategoryValidationDescriptionEmpty(t *testing.T) {
	// Instancia de echo
    e := echo.New()
    h := &handlers.Handler{Categories: nil}

	// Instancia de Category seteada sin Description
	category := models.Category{
		Name        : "Description empty", 
		Description : "",
	}

	// Define la request
    body, _ := json.Marshal(category)
    req := httptest.NewRequest(http.MethodPost, "/categories", bytes.NewReader(body))
    req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
    rec := httptest.NewRecorder()
    c := e.NewContext(req, rec)

	// Ejecuta el metodo para crear categia
    if err := h.CreateCategory(c); err != nil {
        t.Fatal(err)
    }

	// Valida si el codigo de respuesta fue la esperada
    if rec.Code != http.StatusBadRequest {
        t.Errorf("Esperado 400, obtuvo %d", rec.Code)
    }
}