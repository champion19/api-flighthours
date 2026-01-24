# Guía de Docker para Swagger

Esta guía detalla cómo usar Docker para trabajar con Swagger en el proyecto Flighthours API.

---

## 🚀 Guía Rápida - 3 Pasos

### ⚡ Si Solo Quieres Ver la Documentación:

```bash
# 1. Iniciar Swagger UI
docker-compose -f docker-compose.swagger.yml up -d swagger-ui

# 2. Abrir navegador
open http://localhost:8082

# 3. Cuando termines, detener
docker-compose -f docker-compose.swagger.yml down
```

### 🔄 Si Necesitas Regenerar la Documentación:

```bash
# Opción más fácil: usar el script
./swagger.sh generate

# Opción Docker:
docker-compose -f docker-compose.swagger.yml --profile tools run --rm swagger-generator
```

---

## 📦 Contenedores Disponibles

### 1. Swagger UI (Interfaz Web)

Contenedor standalone con la interfaz de Swagger UI para visualizar y probar la API.

### 2. Swagger Generator (Herramienta CLI)

Contenedor con `swag` CLI para generar documentación desde anotaciones en el código.

## 🚀 Guía de Uso

### Opción 1: Swagger UI Standalone

Inicia un contenedor con Swagger UI que lee los archivos de documentación:

```bash
# Iniciar Swagger UI
docker-compose -f docker-compose.swagger.yml up -d swagger-ui

# Ver logs
docker-compose -f docker-compose.swagger.yml logs -f swagger-ui

# Acceder a la UI
open http://localhost:8082

# Detener
docker-compose -f docker-compose.swagger.yml down
```

**Características**:
- Puerto: `8082`
- Lee: `./platform/swaggo/swagger.json`
- Permite probar endpoints directamente
- Actualización automática al recargar página (si regeneras docs)

### Opción 2: Generar Documentación con Docker

#### Método A: Docker Compose (Recomendado)

```bash
# Ejecutar generador una vez
docker-compose -f docker-compose.swagger.yml --profile tools run --rm swagger-generator

# Los archivos se generan en: ./platform/swaggo/
```

#### Método B: Docker Run Directo

```bash
# 1. Construir imagen del generador
docker build -f platform/swaggo/Dockerfile.swag -t flighthours-swag .

# 2. Generar documentación
docker run --rm \
  -v $(pwd):/app \
  flighthours-swag \
  init -g cmd/main.go -o platform/swaggo --parseDependency --parseInternal

# 3. Verificar archivos generados
ls -la platform/swaggo/
```

### Opción 3: Ambos Servicios Juntos

```bash
# Iniciar Swagger UI y tener el generador disponible
docker-compose -f docker-compose.swagger.yml --profile tools up -d

# Regenerar docs en cualquier momento
docker-compose -f docker-compose.swagger.yml exec swagger-generator \
  init -g cmd/main.go -o platform/swaggo --parseDependency --parseInternal
```

## 🔧 Configuración Detallada

### Variables de Entorno de Swagger UI

Edita `docker-compose.swagger.yml` para personalizar:

```yaml
environment:
  # Archivo de especificación
  SWAGGER_JSON: /swagger/swagger.json          # Usar JSON
  # SWAGGER_YAML: /swagger/swagger.yaml        # O usar YAML

  # URL base de tu API
  API_URL: http://localhost:8081

  # Configuración de interfaz
  DEEP_LINKING: "true"                         # URLs profundas
  DISPLAY_OPERATION_ID: "true"                 # Mostrar IDs de operación
  DEFAULT_MODELS_EXPAND_DEPTH: 3               # Profundidad de modelos
  DISPLAY_REQUEST_DURATION: "true"             # Duración de requests
  DOC_EXPANSION: "list"                        # Expandir: none/list/full
  FILTER: "true"                               # Filtro de búsqueda
  TRY_IT_OUT_ENABLED: "true"                   # Habilitar "Try it out"

  # Validación
  VALIDATOR_URL: "https://validator.swagger.io/validator"
```

### Puertos Personalizados

Para cambiar el puerto de Swagger UI:

```yaml
services:
  swagger-ui:
    ports:
      - "9090:8080"  # Ahora accede en http://localhost:9090
```

### Volúmenes

El volumen monta la documentación generada:

```yaml
volumes:
  - ./platform/swaggo:/swagger:ro  # :ro = read-only
```

## 📝 Comandos del Generador (swag CLI)

### Comandos Básicos

```bash
# Generar docs (completo)
docker run --rm -v $(pwd):/app flighthours-swag \
  init -g cmd/main.go -o platform/swaggo

# Con parsing de dependencias
docker run --rm -v $(pwd):/app flighthours-swag \
  init -g cmd/main.go -o platform/swaggo --parseDependency --parseInternal

# Formatear anotaciones
docker run --rm -v $(pwd):/app flighthours-swag fmt

# Ver ayuda
docker run --rm flighthours-swag --help

# Ver versión
docker run --rm flighthours-swag --version
```

### Flags Útiles

- `-g, --generalInfo`: Archivo con info general (default: `main.go`)
- `-o, --output`: Directorio de salida (default: `./docs`)
- `--parseDependency`: Parsear dependencias externas
- `--parseInternal`: Parsear paquetes internos
- `--parseDepth`: Profundidad de parsing (default: 100)
- `--instanceName`: Nombre de instancia (default: `swagger`)

## 🔄 Workflows Comunes

### 1. Desarrollo Local

```bash
# Terminal 1: Iniciar Swagger UI
docker-compose -f docker-compose.swagger.yml up swagger-ui

# Terminal 2: Hacer cambios en código
# ... editar handlers con anotaciones ...

# Terminal 3: Regenerar docs
docker-compose -f docker-compose.swagger.yml --profile tools run --rm swagger-generator

# Recargar http://localhost:8082 para ver cambios
```

### 2. CI/CD Pipeline

```yaml
# .github/workflows/swagger.yml (ejemplo)
name: Validate Swagger Docs

on: [push, pull_request]

jobs:
  swagger:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3

      - name: Build Swagger Generator
        run: docker build -f platform/swaggo/Dockerfile.swag -t swag-gen .

      - name: Generate Docs
        run: |
          docker run --rm -v $(pwd):/app swag-gen \
            init -g cmd/main.go -o platform/swaggo

      - name: Check for changes
        run: |
          git diff --exit-code platform/swaggo/ || \
          (echo "Swagger docs not up to date!" && exit 1)

      - name: Validate JSON
        run: |
          docker run --rm -v $(pwd):/app swag-gen \
            fmt -d platform/swaggo
```

### 3. Integración con Hot Reload

```bash
# Usar con air/nodemon para auto-regenerar

# Instalar air (opcional)
docker run --rm -v $(pwd):/app cosmtrek/air:latest init

# Configurar .air.toml
[build]
  cmd = "docker-compose -f docker-compose.swagger.yml --profile tools run --rm swagger-generator && go build -o tmp/main cmd/main.go"

# Iniciar con hot reload
air
```

## 🐛 Troubleshooting

### Error: "Cannot connect to Docker daemon"

```bash
# Verificar que Docker está corriendo
docker ps

# En Mac/Windows, asegúrate que Docker Desktop está iniciado
```

### Error: "Permission denied" al generar docs

```bash
# Dar permisos al directorio
chmod -R 755 platform/swaggo/

# O ejecutar con usuario actual
docker run --rm -v $(pwd):/app --user $(id -u):$(id -g) flighthours-swag \
  init -g cmd/main.go -o platform/swaggo
```

### Swagger UI muestra "Failed to load API definition"

```bash
# 1. Verificar que swagger.json existe
ls -la platform/swaggo/swagger.json

# 2. Verificar formato JSON
docker run --rm -v $(pwd):/app flighthours-swag fmt

# 3. Ver logs del contenedor
docker-compose -f docker-compose.swagger.yml logs swagger-ui

# 4. Verificar volumen montado
docker-compose -f docker-compose.swagger.yml exec swagger-ui ls -la /swagger/
```

### Cambios no se reflejan en Swagger UI

```bash
# 1. Regenerar docs
docker-compose -f docker-compose.swagger.yml --profile tools run --rm swagger-generator

# 2. Reiniciar Swagger UI
docker-compose -f docker-compose.swagger.yml restart swagger-ui

# 3. Limpiar caché del navegador (Ctrl+Shift+R)
```

### Contenedor no inicia

```bash
# Ver logs detallados
docker-compose -f docker-compose.swagger.yml logs --tail=100 swagger-ui

# Verificar puertos ocupados
lsof -i :8082
netstat -an | grep 8082

# Cambiar puerto si está ocupado
docker-compose -f docker-compose.swagger.yml up swagger-ui -p 8083:8080
```

## 🔐 Seguridad

### Producción

**IMPORTANTE**: No expongas Swagger UI en producción sin autenticación.

```yaml
# docker-compose.swagger.yml (ejemplo con reverse proxy)
services:
  swagger-ui:
    # ... configuración ...
    labels:
      - "traefik.enable=true"
      - "traefik.http.routers.swagger.rule=Host(`swagger.internal.example.com`)"
      - "traefik.http.routers.swagger.middlewares=auth"
      - "traefik.http.middlewares.auth.basicauth.users=admin:$$apr1$$..."
```

### Protección con Basic Auth

```bash
# Generar usuario/contraseña
htpasswd -nb admin password

# Agregar a docker-compose.yml
environment:
  SWAGGER_BASICAUTH: "admin:$$apr1$$..."
```

## 📊 Monitoreo y Logs

```bash
# Ver logs en tiempo real
docker-compose -f docker-compose.swagger.yml logs -f

# Solo Swagger UI
docker-compose -f docker-compose.swagger.yml logs -f swagger-ui

# Ver estadísticas de recursos
docker stats flighthours-swagger-ui

# Inspeccionar contenedor
docker inspect flighthours-swagger-ui
```

## 🧹 Limpieza

```bash
# Detener servicios
docker-compose -f docker-compose.swagger.yml down

# Remover volúmenes
docker-compose -f docker-compose.swagger.yml down -v

# Remover imágenes
docker rmi flighthours-swag
docker rmi swaggerapi/swagger-ui:latest

# Limpiar todo (cuidado)
docker system prune -a
```

## 🚀 Optimizaciones

### Multi-stage Build (Opcional)

Crear un Dockerfile optimizado:

```dockerfile
# Dockerfile.swagger-optimized
FROM golang:1.23-alpine AS builder
RUN go install github.com/swaggo/swag/cmd/swag@latest

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /go/bin/swag /usr/local/bin/swag
WORKDIR /app
ENTRYPOINT ["swag"]
```

```bash
# Construir versión optimizada
docker build -f platform/swaggo/Dockerfile.swagger-optimized -t flighthours-swag:optimized .
```

### Caché de Capas

```bash
# Reutilizar caché entre builds
docker build --cache-from flighthours-swag:latest \
  -f platform/swaggo/Dockerfile.swag \
  -t flighthours-swag:latest .
```

## 🌐 Integración con Otros Servicios

### Con Nginx

```nginx
# /etc/nginx/conf.d/swagger.conf
server {
    listen 80;
    server_name swagger.example.com;

    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

### Con Traefik

Ya configurado en `docker-compose.swagger.yml` con labels.

### Con Kubernetes

```yaml
# kubernetes/swagger-ui-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: swagger-ui
spec:
  replicas: 1
  selector:
    matchLabels:
      app: swagger-ui
  template:
    metadata:
      labels:
        app: swagger-ui
    spec:
      containers:
      - name: swagger-ui
        image: swaggerapi/swagger-ui:latest
        ports:
        - containerPort: 8080
        env:
        - name: SWAGGER_JSON
          value: /swagger/swagger.json
        volumeMounts:
        - name: swagger-docs
          mountPath: /swagger
      volumes:
      - name: swagger-docs
        configMap:
          name: swagger-docs
```

## 📚 Referencias

- [Docker Compose Documentation](https://docs.docker.com/compose/)
- [Swagger UI Docker](https://hub.docker.com/r/swaggerapi/swagger-ui)
- [swaggo/swag](https://github.com/swaggo/swag)

---

**Última actualización**: Diciembre 2024
**Mantenido por**: Flighthours Team
