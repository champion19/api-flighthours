package main

import (
	"log/slog"

	"github.com/champion19/flighthours-api/platform/logger"
	_ "github.com/champion19/flighthours-api/platform/swaggo" // Importar documentos generados por Swag
	"github.com/champion19/flighthours-api/server"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// @title           FlightHours API
// @version         1.0
// @description     API para gestión de empleados de aerolínea y control de horas de vuelo
// @termsOfService  http://swagger.io/terms/

// @contact.name   FlightHours Support
// @contact.email  support@flighthours.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8082
// @BasePath  /flighthours/api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Ingrese el token en formato: Bearer {token}

func main() {
	gin.SetMode(gin.ReleaseMode)
	app := gin.New()
	app.Use(gin.Logger())
	app.Use(gin.Recovery())

	app.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:8081", "http://localhost:8082", "http://localhost:8080", "http://localhost:3001"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	dependencies := server.Bootstrap(app)
	serverAddr := dependencies.Config.GetServerAddress()
	slog.Info(logger.LogAppServerStarting, slog.String("address", serverAddr))

	if err := app.Run(serverAddr); err != nil {
		slog.Error(logger.LogAppServerStartError, slog.String("error", err.Error()))
		return
	}

}
