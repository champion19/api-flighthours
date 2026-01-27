package server

import (
	"log/slog"
	"time"

	"github.com/champion19/api-flighthours/cmd/dependency"
	"github.com/champion19/api-flighthours/handlers"
	"github.com/champion19/api-flighthours/middleware"
	"github.com/champion19/api-flighthours/platform/logger"
	"github.com/champion19/api-flighthours/platform/schema"
	_ "github.com/champion19/api-flighthours/platform/swaggo"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func routing(app *gin.Engine, dependencies *dependency.Dependencies) {
	dependencies.Logger.Info(logger.LogRouteConfiguring)

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

	app.GET("/metrics", gin.WrapH(promhttp.Handler()))

	app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	app.Use(middleware.RequestID())

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
		dependencies.AirlineInteractor,
	)

	validators, err := schema.NewValidator(&schema.DefaultFileReader{})
	if err != nil {
		dependencies.Logger.Error(logger.LogRouteValidatorError, err)
		dependencies.Logger.Fatal(logger.LogRouteValidatorError, err)
		return
	}
	dependencies.Logger.Success(logger.LogRouteValidatorOK)
	validator := middleware.NewMiddlewareValidator(validators)

	app.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "flighthours-backend",
		})
	})

	app.NoRoute(middleware.NotFoundHandler())

	public := app.Group("flighthours/api/v1")
	{

		public.POST("/register", validator.WithValidateRegister(), handler.RegisterEmployee())

		public.POST("/login", handler.Login())

		public.POST("/auth/resend-verification", validator.WithValidateResendVerificationEmail(), handler.ResendVerificationEmail())

		public.POST("/auth/verify-email", handler.VerifyEmailByToken())

		public.POST("/auth/password-reset", validator.WithValidatePasswordResetRequest(), handler.PasswordReset())

		public.POST("/auth/update-password", validator.WithValidateUpdatePassword(), handler.UpdatePassword())

		public.GET("/airlines", handler.ListAirlines())

		public.GET("/airlines/:id", handler.GetAirlineByID())


	}
	protected := app.Group("flighthours/api/v1")
	protected.Use(middleware.RequireAuth(dependencies.EmployeeService, dependencies.MessagingCache, dependencies.JWTValidator))
	{
		protected.POST("/auth/change-password", validator.WithValidateChangePassword(), handler.ChangePassword())

		protected.GET("/employee/me", handler.GetEmployee())

		protected.PUT("/employees/me", validator.WithValidateUpdateEmployee(), handler.UpdateEmployee())

		protected.DELETE("/employees/me", handler.DeleteEmployee())

		protected.POST("/messages", validator.WithValidateMessage(), handler.CreateMessage())

		protected.PUT("/messages/:id", validator.WithValidateMessage(), handler.UpdateMessage())

		protected.DELETE("/messages/:id", handler.DeleteMessage())

		protected.GET("/messages/:id", handler.GetMessageByID())

		protected.GET("/messages", handler.ListMessages())

		protected.POST("/messages/cache/reload", handler.ReloadMessageCache())

	}

	dependencies.Logger.Success(logger.LogRouteConfigured)
}

func Bootstrap(app *gin.Engine) *dependency.Dependencies {
	dependencies, err := dependency.Init()
	if err != nil {
		slog.Error(logger.LogDepInitError, slog.String("error", err.Error()))
		panic(err)
	}
	routing(app, dependencies)
	return dependencies
}
