package performance_test

import (
    "context"
    "encoding/json"
    "math/rand"
    "net/http/httptest"
    "testing"
    "time"

    vegeta "github.com/tsenart/vegeta/v12/lib"
    "github.com/labstack/echo/v4"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"

    "parcial_final/handlers"
    "parcial_final/models"
)

// Generador de string aleatorios
func randomString(n int) string {
    letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
    s := make([]rune, n)
    for i := range s {
        s[i] = letters[rand.Intn(len(letters))]
    }
    return string(s)
}

// Prueba de performance para la creacion de 50 productos
func TestPerformanceCreateCategory(t *testing.T) {
    rand.Seed(time.Now().UnixNano())

    // Define time out
    ctxSetup, cancelSetup := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancelSetup()

	// Configuracion para conexion a base de datos
    client, err := mongo.Connect(ctxSetup, options.Client().ApplyURI("mongodb://localhost:27017"))
    if err != nil {
        t.Fatal(err)
    }

	// Define la base de datos y la coleccion
    userColl := client.Database("test_db_final").Collection("categories")
    if err := userColl.Drop(ctxSetup); err != nil {
        t.Fatal(err)
    }

    // Levanta server en httptest
    h := handlers.NewHandler(nil, userColl, nil)

	// Instancia de echo
    e := echo.New()

	// Ruta para crear productos
    e.POST("/categories", h.CreateCategory)
    ts := httptest.NewServer(e)
    defer ts.Close()

	// Define la frecuencia de request de creacion
    rate := vegeta.Rate{
		Freq: 10, 
		Per: time.Second,
	}

	// Define la duracion de la prueba
    duration := 5 * time.Second

	// Prepara peticiones
    attacker := vegeta.NewAttacker()

	// Instancia de las metricas de las pruebas
    var metrics vegeta.Metrics

	// Organiza peticiones
    targeter := func() vegeta.Targeter {
        return func(tgt *vegeta.Target) error {
			// Instancia seteada de producto
            u := models.Category{
                Name       : "50 registros",
                Description : randomString(11),
            }
			// Codifica 
            body, _ := json.Marshal(u)
			// Organiza la peticion
            *tgt = vegeta.Target{
                Method : "POST",
                URL    : ts.URL + "/categories",
                Body   : body,
                Header : map[string][]string{"Content-Type": {"application/json"}},
            }
            return nil
        }
    }()

	// Ejecuta pruebas
    for res := range attacker.Attack(targeter, rate, duration, "perf-create-products") {
        metrics.Add(res)
    }
    metrics.Close()

    // Crea un nuevo contexto para consulta (o usa context.Background())
    ctxQuery, cancelQuery := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancelQuery()

	// Cuenta los documentos alterados
    count, err := userColl.CountDocuments(ctxQuery, bson.M{})
    if err != nil {
        t.Fatal(err)
    }
    if count == 0 {
        t.Fatal("No se creó ninguna categoria en la base de datos")
    }
    t.Logf("Se crearon %d categorias con éxito", count)
}