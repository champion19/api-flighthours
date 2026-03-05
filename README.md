# FlightHours API

[![Go Version](https://img.shields.io/badge/Go-1.25-00ADD8?logo=go&logoColor=white)](https://go.dev/)
[![License: Apache 2.0](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)
[![SonarCloud](https://img.shields.io/badge/SonarCloud-Analyzed-F3702A?logo=sonarcloud)](https://sonarcloud.io/project/overview?id=champion19_api-flighthours)

> API RESTful del backend de **FlightHours** — plataforma que permite a pilotos y tripulantes registrar, consultar y gestionar sus horas de vuelo, bitácoras diarias y datos operacionales de aeronaves, cumpliendo con los estándares de la Aeronáutica Civil.

---

## 📖 Tabla de Contenidos

- [Visión](#-visión)
- [Funcionalidades Principales](#-funcionalidades-principales)
- [Stack Tecnológico](#-stack-tecnológico)
- [Arquitectura](#-arquitectura)
- [Estructura del Proyecto](#-estructura-del-proyecto)
- [Inicio Rápido](#-inicio-rápido)
- [Configuración](#️-configuración)
- [Observabilidad](#-observabilidad)
- [Documentación de la API](#-documentación-de-la-api)
- [Testing](#-testing)
- [Referencia de Puertos](#-referencia-de-puertos)
- [Contribuciones](#-contribuciones)
- [Seguridad](#️-seguridad)
- [Licencia](#-licencia)

---

## 🎯 Visión

Para pilotos y administradores en general que presentan deficiencias en el seguimiento y monitoreo de sus horas de vuelo, lo que puede derivar en problemas de salud, económicos, sociales y legales, **Flight Hours** es una aplicación web y móvil que incorpora las mejores prácticas y técnicas basadas en estudios médicos, físicos, laborales y emocionales; pensada en satisfacer las necesidades del usuario final.

Este sistema garantizará un acceso confiable y seguro a la información, priorizando siempre el derecho a la protección de los datos almacenados.

A diferencia de aplicaciones como **Logten Pro**, **ForeFlight**, **FlightLogger**, **CrewLounge** y **ZuluLog**, nuestro producto ofrecerá una interacción amigable e intuitiva y estará diseñado bajo principios de experiencia de usuario, será multiplataforma (compatible con ordenadores, tabletas y smartphones), y podrá ser administrado y gestionado conforme a la normativa legal vigente aplicable a cada cliente.

---

## ✨ Funcionalidades Principales

- **Gestión de Empleados** — Registro, autenticación, verificación de email y flujos de restablecimiento de contraseña
- **Integración con Identity Provider** — Keycloak para gestión segura de identidad y acceso (OIDC/JWT/JWKS)
- **Bitácora Diaria (Daily Logbook)** — Registro de encabezados y segmentos de vuelo con tiempos operacionales (out, takeoff, landing, in)
- **Gestión de Aeronaves** — Modelos de aeronaves, matrículas (tail numbers) y familias de aeronaves
- **Datos Maestros** — Aerolíneas, aeropuertos, motores, fabricantes y tipos de tripulante
- **Rutas y Rutas por Aerolínea** — Rutas físicas origen-destino y asignación operativa por aerolínea
- **Resúmenes de Horas de Vuelo** — Cálculo automático de tiempos de bloque, aire y servicio
- **Control de Acceso por Roles (RBAC)** — Endpoints diferenciados para pilotos y administradores
- **Mensajería Dinámica** — Sistema de mensajes centralizado y manejado desde base de datos con caché en memoria
- **API HATEOAS** — Respuestas hypermedia (Richardson Maturity Model Level 3)
- **Observabilidad Integral** — Métricas con Prometheus, dashboards en Grafana, logs centralizados con Loki
- **Documentación de API** — Especificaciones OpenAPI/Swagger generadas automáticamente
- **Calidad de Código** — SonarCloud en CI/CD con GitHub Actions

---

## 🛠 Stack Tecnológico

| Categoría | Tecnología |
| --- | --- |
| **Lenguaje** | Go 1.25 |
| **Web Framework** | [Gin](https://github.com/gin-gonic/gin) |
| **Base de Datos** | MySQL 8.0 |
| **Identity Management** | [Keycloak](https://www.keycloak.org/) (OIDC + JWKS validation) |
| **Email** | [Resend](https://resend.com/) |
| **Métricas** | Prometheus |
| **Dashboards** | Grafana |
| **Log Aggregation** | Loki + Promtail |
| **Structured Logging** | `log/slog` |
| **Documentación API** | Swagger (swaggo) |
| **Validación de Entrada** | JSON Schema |
| **Calidad de Código** | SonarCloud, golangci-lint, staticcheck |
| **Testing** | testify, go-sqlmock, testcontainers |
| **Containerización** | Docker + Docker Compose |
| **CI/CD** | GitHub Actions |

---

## 🏗 Arquitectura

El backend de FlightHours sigue una **Clean / Hexagonal Architecture** (Ports & Adapters) con capas bien definidas:

```text
┌─────────────────────────────────────────────────┐
│                   handlers/                      │  ← Presentation Layer (HTTP controllers)
├─────────────────────────────────────────────────┤
│                  middleware/                      │  ← Cross-cutting concerns (auth, CORS, metrics)
├─────────────────────────────────────────────────┤
│                    core/                         │
│   ┌──────────────┐  ┌────────────────────────┐  │
│   │    ports/     │  │     interactor/        │  │  ← Domain Layer (business logic + interfaces)
│   │ (interfaces)  │  │  (use cases / services)│  │
│   └──────────────┘  └────────────────────────┘  │
├─────────────────────────────────────────────────┤
│                  platform/                       │  ← Infrastructure (DB, Keycloak, Prometheus, etc.)
└─────────────────────────────────────────────────┘
```

**Principios clave:**

- Las dependencias apuntan **hacia adentro** — `platform` implementa las interfaces de `core/ports`
- La lógica de negocio en `core/interactor` es **agnóstica al framework**
- `handlers` solo orquesta el mapeo de request/response HTTP
- `middleware` provee autenticación (JWT/JWKS), RBAC, rate limiting, trazabilidad y CORS

---

## 📁 Estructura del Proyecto

```text
api-flighthours/
├── cmd/
│   ├── main.go                 # Punto de entrada de la aplicación
│   └── dependency/             # Wiring de inyección de dependencias
├── config/                     # Archivos de configuración (JSON) + config loader
├── core/
│   ├── ports/                  # Interfaces (contratos de repositorio y servicio)
│   └── interactor/             # Casos de uso / lógica de negocio
│       └── services/
│           └── domain/         # Modelos de dominio puros y constantes de error
├── handlers/                   # Controladores HTTP (Gin handlers) + DTOs + HATEOAS
├── middleware/                  # Auth, CORS, RBAC, rate limiter, request ID, métricas
├── server/                     # Bootstrap del servidor y registro de rutas
├── platform/                   # Implementaciones de infraestructura
│   ├── cache/                  # Caché en memoria (mensajes del sistema)
│   ├── cookie/                 # Utilidades de manejo de cookies HTTP
│   ├── databases/              # Repositorios MySQL
│   ├── grafana/                # Dashboards y provisioning de Grafana
│   ├── identity_provider/      # Integración con Keycloak (gocloak)
│   ├── jwt/                    # Validación JWT / JWKS (keyfunc)
│   ├── k6/                     # Escenarios de pruebas de carga K6
│   ├── log-rotator/            # Servicio de rotación de logs
│   ├── logger/                 # Logging estructurado (slog)
│   ├── loki/                   # Configuración de Loki
│   ├── prometheus/             # Configuración de scrape de Prometheus
│   ├── promtail/               # Configuración de envío de logs con Promtail
│   ├── schema/                 # Validación con JSON Schema
│   └── swaggo/                 # Documentación Swagger generada
├── mocks/                      # Implementaciones mock con testify
├── tools/                      # Utilidades (ID encoder con HashIDs, helpers)
├── scripts/                    # Scripts de automatización (observabilidad, Swagger)
├── docs/                       # Swagger docs + documentación de releases
├── documents/                  # SQL schemas, documentación complementaria
├── .github/workflows/          # CI/CD (SonarCloud analysis)
├── docker-compose.grafana.yml  # Stack de observabilidad (Grafana + Prometheus + Loki)
├── docker-compose.keycloak.yml # Keycloak local
├── docker-compose.swagger.yml  # Swagger UI standalone
├── Containerfile               # Imagen custom de Keycloak
├── .pre-commit-config.yaml     # Pre-commit hooks (golangci-lint)
├── sonar-project.properties    # Integración con SonarCloud
└── swagger.sh                  # Generador de documentación Swagger
```

---

## 🚀 Inicio Rápido

**Tiempo estimado:** 15–20 minutos

### Prerrequisitos

| Herramienta | Versión | Propósito |
| --- | --- | --- |
| **Go** | ≥ 1.25 | Runtime del lenguaje |
| **Docker** | ≥ 24.0 | Runtime de contenedores |
| **Docker Compose** | ≥ 2.20 | Orquestación multi-contenedor |
| **Git** | Cualquiera | Control de versiones |

### Paso 1: Clonar el Repositorio

```bash
git clone https://github.com/champion19/api-flighthours.git
cd api-flighthours
```

### Paso 2: Iniciar MySQL

El proyecto utiliza una instancia MySQL con el contenedor `mysql-flighthours`.

```bash
docker start mysql-flighthours
```

> Si es la primera vez, consultar la documentación de Docker para crear el contenedor con la base de datos `flightDb` en el puerto `3306`.

**Verificación:**

```bash
docker ps --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"
# Esperado: mysql-flighthours (running)
```

### Paso 3: Iniciar Keycloak

```bash
docker compose -f docker-compose.keycloak.yml up -d
```

**Verificación:** Abrir [http://localhost:8080](http://localhost:8080) — debería verse la consola de administración de Keycloak.

> **Esperar** a que el contenedor reporte `healthy` antes de continuar (~30 segundos).

### Paso 4: Configurar la Aplicación

#### 4a. Variables de Entorno

```bash
cp .env.example .env
```

El `.env` configura los siguientes servicios:

| Sección | Variables clave | Propósito |
| --- | --- | --- |
| **Keycloak** | `KEYCLOAK_CLIENT_SECRET`, `KEYCLOAK_ADMIN_PASSWORD` | Identity provider |
| **Resend** | `RESEND_API_KEY` | Envío de emails (reset de contraseña, verificación) |
| **ID Encoder** | `ID_ENCODER_SECRET` | Codificación de HashID |
| **Database** | `DB_PASSWORD` | Contraseña MySQL (si difiere del JSON) |
| **Cookie** | `COOKIE_DOMAIN`, `COOKIE_SECURE` | Configuración de cookies HTTP |

> **Tip:** Las variables de entorno tienen prioridad sobre el archivo JSON. Ver `config/config.go` para detalles.

#### 4b. Configuración por Archivo (JSON)

La aplicación soporta configuración basada en archivos JSON por entorno:

| Entorno | Archivo | Trigger |
| --- | --- | --- |
| **Local** (por defecto) | `config/local-config.json` | `APP_ENV` sin definir o `local` |
| **Railway** | `config/railway-config.json` | `APP_ENV=railway` |

### Paso 5: Instalar Dependencias

```bash
go mod tidy
```

### Paso 6: Ejecutar la Aplicación

```bash
go run ./cmd/main.go
```

### ✅ Verificación

| Comprobación | Cómo verificar |
| --- | --- |
| **Servidor corriendo** | La consola muestra `Listening on 0.0.0.0:8082` |
| **Health endpoint** | `curl http://localhost:8082/health` retorna `200 OK` |
| **Metrics endpoint** | `curl http://localhost:8082/metrics` retorna datos de Prometheus |
| **Swagger UI** | Abrir [http://localhost:8082/swagger/index.html](http://localhost:8082/swagger/index.html) |

---

## ⚙️ Configuración

La aplicación soporta múltiples entornos mediante archivos de configuración JSON y sobreescrituras por variables de entorno.

**Precedencia:** Variables de entorno (`.env`) → Archivo JSON → Valores por defecto.

```bash
# Base de datos
DB_PASSWORD=your-password

# Keycloak
KEYCLOAK_SERVER_URL=http://localhost:8080
KEYCLOAK_CLIENT_SECRET=your-secret
KEYCLOAK_ADMIN_PASSWORD=your-admin-pass

# Servicios externos
RESEND_API_KEY=re_your_api_key

# ID Encoder
ID_ENCODER_SECRET=your-encoder-secret

# Cookie
COOKIE_DOMAIN=localhost
COOKIE_SECURE=false
```

> Ver `.env` para la lista completa de variables de entorno disponibles.

---

## 📊 Observabilidad

### Iniciar el Stack Completo de Observabilidad

```bash
docker compose -f docker-compose.grafana.yml up -d
```

Esto inicia **5 servicios**:

| Servicio | Puerto | Propósito |
| --- | --- | --- |
| **Grafana** | [localhost:3000](http://localhost:3000) | Dashboards y visualización |
| **Prometheus** | [localhost:9090](http://localhost:9090) | Scraping y almacenamiento de métricas |
| **Loki** | localhost:3100 | Agregación de logs |
| **Promtail** | — | Agente de envío de logs |
| **Log Rotator** | — | Rotación automática de logs (retención de 7 días) |

> Credenciales por defecto de Grafana: `admin` / `admin`

### Arquitectura de Logs

```text
Backend (slog) → Archivos JSON → Promtail → Loki → Grafana
                 /tmp/flighthours-logs/
```

---

## 📚 Documentación de la API

### Swagger UI (Integrado)

La documentación de la API se sirve directamente desde la aplicación en ejecución:

```text
http://localhost:8082/swagger/index.html
```

### Swagger UI (Contenedor Independiente)

Para un Swagger UI independiente con hot-reload:

```bash
docker compose -f docker-compose.swagger.yml up -d
# Abrir http://localhost:8082
```

### Generación de Swagger

Para regenerar la documentación Swagger después de modificar anotaciones en handlers:

```bash
./swagger.sh
```

---

## 🧪 Testing

### Ejecutar Todos los Tests

```bash
go test ./... -v
```

### Reporte de Cobertura

```bash
go test ./... -coverprofile=coverage.out -covermode=atomic
go tool cover -html=coverage.out -o coverage.html
# Genera: coverage.html (abrir en navegador)
```

### Threshold de Cobertura

El proyecto exige un **mínimo de 80% de cobertura**, verificado mediante SonarCloud en el pipeline de CI/CD.

### Herramientas de Testing

| Herramienta | Propósito |
| --- | --- |
| `testify` | Assertions y mocks |
| `go-sqlmock` | Mock de conexiones SQL |
| `testcontainers` | Tests de integración con MySQL real |

---

## 🔌 Referencia de Puertos

| Puerto | Servicio | Entorno |
| --- | --- | --- |
| `8082` | FlightHours API | Local |
| `3306` | MySQL (flightDb) | Local |
| `8080` | Keycloak (HTTP) | Local |
| `8443` | Keycloak (HTTPS) | Producción |
| `9000` | Keycloak (Health/Metrics) | Local |
| `3000` | Grafana | Local |
| `9090` | Prometheus | Local |
| `3100` | Loki | Local |

---

## 🤝 Contribuciones

¡Las contribuciones son bienvenidas! Para contribuir:

1. Crear una rama desde `develop` siguiendo Git Flow
2. Utilizar [Conventional Commits](https://www.conventionalcommits.org/) para mensajes
3. Asegurarse de que `go test ./...` pasa sin errores
4. Ejecutar `golangci-lint run` para validar calidad de código
5. Crear un Pull Request hacia `develop`

### Pre-commit Hooks

El proyecto incluye hooks de pre-commit con `golangci-lint`:

```bash
# Los hooks se configuran automáticamente con .pre-commit-config.yaml
pre-commit install
```

---

## 🛡️ Seguridad

### Nota sobre Credenciales Incluidas

> **Este repositorio es un proyecto de grado académico.** Las credenciales incluidas en `config/local-config.json` son exclusivamente para el entorno de desarrollo local y se incluyen intencionalmente para facilitar la ejecución y evaluación del proyecto sin configuración adicional.

En un entorno productivo, **todas las credenciales se inyectan mediante variables de entorno** (archivo `.env` o secretos del proveedor de hosting), nunca desde archivos de configuración commiteados.

| Archivo | ¿Commiteado? | Propósito |
| --- | :---: | --- |
| `config/local-config.json` | ✅ Sí | Configuración local lista para ejecutar |
| `config/railway-config.json` | ✅ Sí | Configuración para Railway (placeholders) |
| `.env` | ❌ No | Variables de entorno con secrets reales |

### Estrategia de Configuración

```text
Prioridad de configuración (de mayor a menor):

  Variables de entorno (.env)  →  Archivo JSON (local/railway-config.json)  →  Valores por defecto
```

Para más detalles, ver la sección [Configuración](#️-configuración) y el archivo `config/config.go`.

---

## 📄 Licencia

Este proyecto está licenciado bajo la [Licencia Apache 2.0](LICENSE).

---

<p align="center">
  <sub>Hecho con ❤️ por Champion19</sub>
</p>
