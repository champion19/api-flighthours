package server

import (
	"log/slog"
	"time"

	"github.com/champion19/flighthours-api/cmd/dependency"
	"github.com/champion19/flighthours-api/handlers"
	"github.com/champion19/flighthours-api/middleware"
	"github.com/champion19/flighthours-api/platform/logger"
	"github.com/champion19/flighthours-api/platform/schema"
	_ "github.com/champion19/flighthours-api/platform/swaggo"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func routing(app *gin.Engine, dependencies *dependency.Dependencies) {
	dependencies.Logger.Info(logger.LogRouteConfiguring)

	// CORS configuration - Allow requests from Keycloak (localhost:8080) and other origins
	// This is required for the email verification flow from Keycloak's theme pages
	corsConfig := cors.Config{
		AllowOrigins:     []string{"http://localhost:8080", "http://localhost:8081", "http://localhost:3001"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Request-ID"},
		ExposeHeaders:    []string{"Content-Length", "Location"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	app.Use(cors.New(corsConfig))
	dependencies.Logger.Info("CORS middleware configured")

	// Endpoint de métricas de Prometheus
	app.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Documentación Swagger
	app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Middleware de Prometheus para capturar métricas
	app.Use(middleware.RequestID())

	// Apply Prometheus metrics tracking middleware
	app.Use(middleware.TrackMetrics())

	errorHandler := middleware.NewErrorHandler(dependencies.MessagingCache)
	app.Use(errorHandler.Handle())

	handler := handlers.New(
		dependencies.EmployeeService,
		dependencies.Interactor,
		dependencies.IDEncoder,
		dependencies.ResponseHandler,
		dependencies.MessageInteractor,
		dependencies.MessagingCache,
	)

	validators, err := schema.NewValidator(&schema.DefaultFileReader{})
	if err != nil {
		dependencies.Logger.Error(logger.LogRouteValidatorError, err)
		dependencies.Logger.Fatal(logger.LogRouteValidatorError, err)
		return
	}
	dependencies.Logger.Success(logger.LogRouteValidatorOK)
	validator := middleware.NewMiddlewareValidator(validators)

	// Health check endpoint (no authentication required)
	app.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "flighthours-backend",
		})
	})

	app.NoRoute(middleware.NotFoundHandler())

	// ===========================================
	// PUBLIC ROUTES (no authentication required)
	// ===========================================
	public := app.Group("flighthours/api/v1")
	{
		// ---- Authentication ----
		// POST /register - New user registration (HU17)
		public.POST("/register", validator.WithValidateRegister(), handler.RegisterEmployee())

		// POST /login - Returns JWT tokens (HU21)
		public.POST("/login", handler.Login())

		// ---- Email Verification (Helper) ----
		// POST /auth/resend-verification - Resend verification email
		public.POST("/auth/resend-verification", validator.WithValidateResendVerificationEmail(), handler.ResendVerificationEmail())

		// POST /auth/verify-email - Verify email with token
		public.POST("/auth/verify-email", handler.VerifyEmailByToken())
	}

	// ===========================================
	// PROTECTED ROUTES (authentication required)
	// ===========================================
	protected := app.Group("flighthours/api/v1")
	// Use the RequireAuth middleware from jwt_middleware.go
	// This validates JWT tokens and injects the authenticated user into context
	protected.Use(middleware.RequireAuth(dependencies.EmployeeService, dependencies.MessagingCache, dependencies.JWTValidator))
	{
		// ---- Authenticated User Endpoints ----
		// POST /auth/change-password - Change password (HU20)
		protected.POST("/auth/change-password", validator.WithValidateChangePassword(), handler.ChangePassword())

		// GET /employee/me - Get authenticated user profile (secure: no ID exposure)
		protected.GET("/employee/me", handler.GetMe())

		// ---- Messages Management (Protected - Admin) ----
		// POST /messages - Create new message
		protected.POST("/messages", validator.WithValidateMessage(), handler.CreateMessage())

		// PUT /messages/:id - Update existing message
		protected.PUT("/messages/:id", validator.WithValidateMessage(), handler.UpdateMessage())

		// DELETE /messages/:id - Delete message
		protected.DELETE("/messages/:id", handler.DeleteMessage())

		// GET /messages/:id - Get message by ID
		protected.GET("/messages/:id", handler.GetMessageByID())

		// GET /messages - List messages (with optional filters)
		protected.GET("/messages", handler.ListMessages())

		// POST /messages/cache/reload - Reload message cache from DB
		protected.POST("/messages/cache/reload", handler.ReloadMessageCache())
	}

	dependencies.Logger.Success(logger.LogRouteConfigured)
}

func Bootstrap(app *gin.Engine) *dependency.Dependencies {
	// Initialize Prometheus metrics
	dependencies, err := dependency.Init()
	if err != nil {
		slog.Error(logger.LogDepInitError, slog.String("error", err.Error()))
		panic(err)
	}
	routing(app, dependencies)
	return dependencies
}
