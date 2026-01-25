# Documentación Swagger - Flighthours API

Esta carpeta contiene la configuración y documentación de Swagger para la API de Flighthours, implementada con [swaggo/swag](https://github.com/swaggo/swag) para Go.

---

## 🎯 Guía de Inicio - Instalación Completa

### Paso 1: Verificar Instalación

Swagger ya está integrado en el proyecto. Verifica que todo esté en orden:

```bash
# Verificar que las dependencias están instaladas
go mod tidy

# Verificar que los archivos de Swagger existen
ls -la platform/swaggo/
# Debes ver: docs.go, swagger.json, swagger.yaml
```

### Paso 2: Elegir Método de Uso

Tienes **2 opciones** para ver la documentación Swagger:

#### ✅ Opción A: Swagger UI Integrado (MÁS FÁCIL)

1. Inicia la API normalmente:
   ```bash
   go run cmd/main.go
   ```

2. Abre tu navegador en:
   ```
   http://localhost:8081/swagger/index.html
   ```

3. ¡Listo! Ya puedes ver y probar tus endpoints.

#### ✅ Opción B: Swagger UI Docker (INDEPENDIENTE)

1. Inicia Swagger UI en Docker:
   ```bash
   docker-compose -f docker-compose.swagger.yml up -d swagger-ui
   ```

2. Abre tu navegador en:
   ```
   http://localhost:8082
   ```

3. Para detener:
   ```bash
   docker-compose -f docker-compose.swagger.yml down
   ```

### Paso 3: Probar Endpoints

En Swagger UI puedes:
- ✅ Ver todos los endpoints documentados
- ✅ Probar las peticiones directamente
- ✅ Ver los modelos de datos
- ✅ Copiar ejemplos de cURL

---

## 📁 Estructura de Archivos

```
platform/swaggo/
├── docs.go              # Documentación generada automáticamente por swag
├── swagger.json         # Especificación OpenAPI 2.0 en formato JSON
├── swagger.yaml         # Especificación OpenAPI 2.0 en formato YAML
├── Dockerfile.swag      # Dockerfile para generar docs con swag CLI
├── DOCKER.md           # Guía de uso con Docker
└── README.md           # Este archivo
```

## 🚀 Inicio Rápido

### Opción 1: Acceder a Swagger UI integrado (Recomendado)

La API incluye Swagger UI integrado. Una vez que inicies el servidor:

```bash
go run cmd/main.go
```

Accede a la documentación en:
- **Swagger UI**: http://localhost:8081/swagger/index.html
- **JSON Spec**: http://localhost:8081/swagger/doc.json

### Opción 2: Usar Docker Compose para Swagger UI Standalone

```bash
# Iniciar Swagger UI en puerto 8082
docker-compose -f docker-compose.swagger.yml up -d swagger-ui

# Acceder a: http://localhost:8082
```

## 📝 Generar/Actualizar Documentación

### Método 1: Con CLI local

```bash
# Instalar swag CLI
go install github.com/swaggo/swag/cmd/swag@latest

# Generar documentación
swag init -g cmd/main.go -o platform/swaggo --parseDependency --parseInternal

# Verificar
swag fmt
```

### Método 2: Con Docker

```bash
# Construir imagen
docker build -f platform/swaggo/Dockerfile.swag -t flighthours-swag .

# Generar documentación
docker run --rm -v $(pwd):/app flighthours-swag init -g cmd/main.go -o platform/swaggo --parseDependency --parseInternal
```

### Método 3: Con Docker Compose

```bash
# Usar el perfil tools para ejecutar el generador
docker-compose -f docker-compose.swagger.yml --profile tools run --rm swagger-generator
```

## ✍️ Documentar Endpoints

### Estructura básica de anotaciones

Agrega anotaciones en tus handlers de la siguiente manera:

```go
// RegisterEmployee godoc
// @Summary      Registrar nueva cuenta
// @Description  Crea una nueva cuenta de usuario en el sistema
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Param        account  body      handlers.PersonRequest  true  "Datos de la cuenta"
// @Success      201      {object}  middleware.APIResponse{data=handlers.RegistrationResponse}
// @Failure      400      {object}  middleware.APIResponse
// @Failure      409      {object}  middleware.APIResponse
// @Failure      500      {object}  middleware.APIResponse
// @Router       /accounts [post]
func (h *handler) RegisterEmployee() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Implementación...
    }
}
```

### Anotaciones principales

- `@Summary`: Resumen corto del endpoint
- `@Description`: Descripción detallada
- `@Tags`: Categoría/grupo del endpoint
- `@Accept`: Formato de entrada (json, xml, etc.)
- `@Produce`: Formato de respuesta (json, xml, etc.)
- `@Param`: Parámetros de entrada
- `@Success`: Respuesta exitosa
- `@Failure`: Respuesta de error
- `@Router`: Ruta y método HTTP
- `@Security`: Esquema de autenticación requerido

### Ejemplo con autenticación

```go
// GetEmployee godoc
// @Summary      Obtener empleado por ID
// @Description  Obtiene información de un empleado específico
// @Tags         employees
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Employee ID"
// @Success      200  {object}  middleware.APIResponse{data=handlers.EmployeeResponse}
// @Failure      401  {object}  middleware.APIResponse
// @Failure      404  {object}  middleware.APIResponse
// @Security     BearerAuth
// @Router       /employees/{id} [get]
func (h *handler) GetEmployee() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Implementación...
    }
}
```

## 🔧 Configuración

### Configuración general (en cmd/main.go)

```go
// @title           Flighthours Backend API
// @version         1.0
// @description     API RESTful para la plataforma Flighthours
// @termsOfService  http://swagger.io/terms/

// @contact.name   Flighthours API Support
// @contact.url    https://flighthours.com/support
// @contact.email  support@flighthours.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8081
// @BasePath  /flighthours/api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
```

### Variables de entorno (Docker)

Puedes configurar Swagger UI con variables de entorno en `docker-compose.swagger.yml`:

```yaml
environment:
  SWAGGER_JSON: /swagger/swagger.json
  API_URL: http://localhost:8081
  DEEP_LINKING: "true"
  TRY_IT_OUT_ENABLED: "true"
```

## 🔐 Seguridad

### Bearer Token Authentication

La API usa autenticación JWT. Para probar endpoints protegidos:

1. Haz clic en el botón **Authorize** 🔓 en Swagger UI
2. Ingresa: `Bearer YOUR_JWT_TOKEN`
3. Haz clic en **Authorize**
4. Ahora puedes probar endpoints protegidos

### OAuth2 (Configuración futura)

Si implementas OAuth2 con Keycloak:

```go
// @securityDefinitions.oauth2.implicit OAuth2Implicit
// @authorizationurl https://keycloak.example.com/auth/realms/flighthours/protocol/openid-connect/auth
// @scope.write Grants write access
// @scope.read Grants read access
```

## 📋 Buenas Prácticas

1. **Regenerar docs después de cambios**: Siempre ejecuta `swag init` después de modificar anotaciones
2. **Validar JSON/YAML**: Usa [Swagger Editor](https://editor.swagger.io/) para validar
3. **Ejemplos realistas**: Incluye ejemplos en tus modelos usando `example` tags
4. **Descripciones claras**: Escribe descripciones útiles para desarrolladores
5. **Versionado**: Mantén la versión sincronizada con tu API
6. **Códigos de estado HTTP correctos**: Usa códigos HTTP apropiados en `@Success` y `@Failure`

## 🐛 Troubleshooting

### Error: "Cannot find docs.go"

```bash
# Regenerar documentación
swag init -g cmd/main.go -o platform/swaggo
```

### Error: "404 Not Found" en /swagger/index.html

Verifica que:
1. El import está presente: `_ "github.com/champion19/flighthours-api/platform/swaggo"`
2. La ruta está registrada en `server/server.go`
3. El servidor está corriendo

### Swagger UI no carga los endpoints

1. Verifica que `swagger.json` contiene tus endpoints
2. Revisa que las anotaciones están correctas
3. Asegúrate de regenerar docs con `swag init`

### CORS errors en Swagger UI

Si usas Swagger UI externo, configura CORS en tu API:

```go
import "github.com/gin-contrib/cors"

app.Use(cors.Default())
```

## 📚 Recursos Adicionales

- [Documentación oficial de swaggo](https://github.com/swaggo/swag)
- [Especificación OpenAPI 2.0](https://swagger.io/specification/v2/)
- [Swagger UI](https://swagger.io/tools/swagger-ui/)
- [Ejemplos de anotaciones](https://github.com/swaggo/swag#declarative-comments-format)

## 🔄 Workflow Recomendado

1. **Desarrollo**:
   ```bash
   # Hacer cambios en handlers
   # Agregar/actualizar anotaciones
   swag fmt  # Formatear anotaciones
   swag init -g cmd/main.go -o platform/swaggo
   go run cmd/main.go
   # Probar en http://localhost:8081/swagger/index.html
   ```

2. **CI/CD**:
   ```bash
   # En tu pipeline
   swag init -g cmd/main.go -o platform/swaggo
   # Validar que docs.go fue actualizado
   git diff --exit-code platform/swaggo/docs.go || exit 1
   ```

3. **Producción**:
   - Los archivos generados (`docs.go`, `swagger.json`, `swagger.yaml`) deben estar en el repositorio
   - No es necesario instalar swag CLI en producción
   - Swagger UI está integrado en la aplicación

## 📊 Métricas y Monitoreo

Swagger UI + Prometheus:

```bash
# Métricas de la API
curl http://localhost:8081/metrics

# Swagger UI
open http://localhost:8081/swagger/index.html
```

## 🤝 Contribuir

Al agregar nuevos endpoints:

1. Agrega anotaciones completas
2. Regenera la documentación
3. Verifica en Swagger UI
4. Commit ambos: código + docs generados

---

**Última actualización**: Diciembre 2024
**Mantenido por**: Flighthours Team
