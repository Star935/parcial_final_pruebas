# API de Gestión de Productos

Esta es una API sencillas desarrollada en Golang utilizando el framework **Echo**, que permite realizar operaciones CRUD sobre categorias . Incluye pruebas de **performance** con [Vegeta](https://github.com/tsenart/vegeta)
---

## Ejecución del proyecto
El proyecto corre en localhost

### Requisitos

- Go 1.21+
- Navegador instalado

### Puerto
:8080

### Recursos web disponibles
- GET "/products" Recupera todos los productos
- POST "/products" Crea un nuevo producto
- GET "/products/:id" Recupera un producto por su id
- PUT "/products/:id" Actualiza un producto por su id
- DELETE "/products/:id" Elimina un producto por su id

- GET "/categories" Recupera todas las categorias
- POST "/categories" Crea una nueva categoria
- GET "/categories/:id" Recupera una categoria por su id
- PUT "/categories/:id" Actualiza una categoria por su id
- DELETE "/categories/:id" Elimina una categoria por su id

### Nota importante
El comando del pipeline ejecuta todas las pruebas, incluida la de performance_test la cual falla, el siguiente es el fragmento de la informacion de la prueba de performance
```bash
FAIL	parcial_final/performance_test	5.005s
```
La siguiente linea es el estado de las pruebas:
```bash
ok  	parcial_final/tests	0.004s (es ok en el caso de que funcione)
```

### Comando para iniciar el servidor:
```bash
go run main.go
```
### Comando para ejecutar todas las pruebas:
```bash
go test ./tests -v -count=1
```

### Comando para ejecutar las pruebas de performance:
```bash
go test ./tests -run TestPerformanceCreateProducts -v -count=1
```

### Comando para ejecutar prueba unitaria de crear categoria sin nombre:
```bash
go test ./tests -run TestCreateCategoryValidationNameEmpty -v -count=1
```

### Comando para ejecutar prueba unitaria de crear categoria sin descripcion:
```bash
go test ./tests -run TestCreateCategoryValidationDescriptionEmpty -v -count=1
```

### Comando para ejecutar las pruebas end to end:
```bash
go test ./tests -run TestEndToEndCreateSubject -v -count=1
```

### Comando para ejecutar el analisis estatico del codigo:
```bash
gosec ./...
```

### Comando para ejecutar el analisis de las dependencias del codigo:
```bash
govulncheck ./...
```

### Evidencias de las pruebas de performance:
![Screenshot 2025-06-19 201620](https://github.com/user-attachments/assets/0680756e-dae3-4b75-8067-79a01bfce1ba)
![Screenshot 2025-06-19 201634](https://github.com/user-attachments/assets/cbd2b6fb-2ee5-40eb-9dec-ddf6df23987a)
