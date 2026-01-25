# Implementación del Endpoint `GET /employees` - Lista de Empleados para Testing

## Fecha: 2026-01-25

## Descripción
Se implementó un nuevo endpoint `GET /employees` que permite listar todos los empleados registrados en el sistema. Este endpoint está diseñado principalmente para facilitar el testing, ya que permite obtener los IDs ofuscados de los empleados que luego pueden usarse en otros endpoints.

## Cambios Realizados

### 1. Controlador HTTP (`handlers/employee_controller.go`)
- **Nuevo método:** `ListEmployees()` - Handler que procesa la solicitud GET y retorna la lista de empleados con sus IDs ofuscados
- Incluye HATEOAS links para navegación

### 2. Interface del Servicio (`core/ports/input/service.go`)
- **Nuevo método:** `ListEmployees(ctx context.Context) ([]domain.Employee, error)`

### 3. Implementación del Servicio (`core/interactor/services/employee_services.go`)
- **Nuevo método:** `ListEmployees()` - Delega al repositorio para obtener la lista de empleados

### 4. Interface del Repositorio (`core/ports/output/repository.go`)
- **Nuevo método:** `ListEmployees(ctx context.Context) ([]domain.Employee, error)`

### 5. Implementación del Repositorio (`platform/databases/repositories/employee/`)
- **Nuevo archivo:** `list.go` - Implementación del método ListEmployees
- **Actualización en `repository.go`:** Nueva constante `QueryList` con la query SQL

### 6. Mensajes de Log (`platform/logger/log_messages.go`)
- **Nuevas constantes:**
  - `LogEmployeeList` - "Listando empleados"
  - `LogEmployeeListOK` - "Empleados listados exitosamente"
  - `LogEmployeeListError` - "Error listando empleados"

### 7. Mensajes de Dominio (`core/interactor/services/domain/erros.go`)
- **Nueva constante:** `MsgEmployeeListOK = "MOD_U_LIST_EXI_00006"` - Código para respuesta exitosa

### 8. Configuración de Rutas (`server/server.go`)
- **Nueva ruta:** `GET /employees` en el grupo protegido (requiere autenticación)

### 9. Documentación Bruno (`docs/employee_bruno/employee/List Employees.bru`)
- **Nuevo archivo:** Request Bruno para probar el endpoint
- Incluye script post-response que guarda automáticamente el primer employee-id

### 10. Mocks de Tests
- **Actualizados:**
  - `handlers/http_test.go` - fakeService y fakeServiceErr
  - `core/interactor/interactor_test.go` - fakeService
  - `core/interactor/interactor_readonly_test.go` - fakeServiceForReadOnly
  - `core/interactor/services/employee_services_test.go` - fakeRepo
  - `mocks/mock_repository.go` - MockRepository

## Uso del Endpoint

### Request
```http
GET /flighthours/api/v1/employees
Authorization: Bearer {access_token}
Content-Type: application/json
```

### Response (200 OK)
```json
{
  "success": true,
  "code": "MOD_U_LIST_EXI_00006",
  "message": "Lista de empleados obtenida exitosamente",
  "data": {
    "employees": [
      {
        "id": "obfuscated_id",
        "name": "Emma",
        "email": "emma@yopmail.com",
        "role": "pilot",
        "active": true,
        "_links": [
          {
            "rel": "self",
            "href": "/employees/{id}",
            "method": "GET"
          }
        ]
      }
    ],
    "count": 1,
    "_links": [
      {
        "rel": "register",
        "href": "/register",
        "method": "POST"
      }
    ]
  }
}
```

## Verificación
- ✅ Build exitoso: `go build ./...`
- ✅ Tests pasando: `go test ./...`

## Notas
- Este endpoint requiere autenticación (token JWT válido)
- Los IDs retornados están ofuscados mediante el sistema HashIDs
- El script post-response de Bruno guarda automáticamente el primer employee-id para facilitar testing de otros endpoints
