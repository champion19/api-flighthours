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

var log = logger.NewSlogLogger()

func routing(app *gin.Engine, dependencies *dependency.Dependencies) {
	log.Info(logger.LogRouteConfiguring)

	corsConfig := cors.Config{
		AllowOrigins:     []string{"http://localhost:8080", "http://localhost:8081", "http://localhost:3001"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Request-ID"},
		ExposeHeaders:    []string{"Content-Length", "Location"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	app.Use(cors.New(corsConfig))
	log.Info("CORS middleware configured")

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
		dependencies.AirlineEmployeeInteractor,
		dependencies.EngineInteractor,
		dependencies.RouteInteractor,
		dependencies.ManufacturerInteractor,
		dependencies.AirportInteractor,
		dependencies.AirlineRouteInteractor,
		dependencies.DailyLogbookDetailInteractor,
		dependencies.DailyLogbookInteractor,
		dependencies.AircraftModelInteractor,
		dependencies.LicensePlateInteractor,
	)

	validators, err := schema.NewValidator(&schema.DefaultFileReader{})
	if err != nil {
		log.Error(logger.LogRouteValidatorError, err)
		log.Fatal(logger.LogRouteValidatorError, err)
		return
	}
	log.Success(logger.LogRouteValidatorOK)
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

		public.GET("/engines", handler.ListEngines())

		public.GET("/engines/:id", handler.GetEngineByID())

		public.GET("/routes", handler.ListRoutes())

		public.GET("/routes/:id", handler.GetRouteByID())

		public.GET("/airline-routes", handler.ListAirlineRoutes())

		public.GET("/airports", handler.ListAirports())

		public.GET("/airports/:id", handler.GetAirportByID())

		public.GET("/airport-types/:airport_type", handler.GetAirportsByType())

		public.GET("/manufacturers", handler.ListManufacturers())

		public.GET("/manufacturers/:id", handler.GetManufacturerByID())

		public.GET("/aircraft-models", handler.ListAircraftModels())

		public.GET("/aircraft-models/:id", handler.GetAircraftModelByID())

		public.GET("/aircraft-families/:family", handler.GetAircraftModelsByFamily())

	}
	protected := app.Group("flighthours/api/v1")
	protected.Use(middleware.RequireAuth(dependencies.EmployeeService, dependencies.MessagingCache, dependencies.JWTValidator))
	{
		protected.POST("/auth/change-password", validator.WithValidateChangePassword(), handler.ChangePassword())

		protected.GET("/employees", handler.GetEmployee())

		protected.PUT("/employees", validator.WithValidateUpdateEmployee(), handler.UpdateEmployee())

		protected.DELETE("/employees", handler.DeleteEmployee())

		protected.GET("/employees/airline", handler.GetEmployeeAirlineInfo())

		protected.PUT("/employees/airline", validator.WithValidateAddAirlineEmployee(), handler.AddEmployeeAirlineInfo())

		protected.PUT("/employees/airline-info", validator.WithValidateUpdateAirlineEmployee(), handler.UpdateEmployeeAirlineInfo())

		protected.PATCH("/employees/airline/activate", handler.ActivateEmployeeAirlineInfo())

		protected.PATCH("/employees/airline/deactivate", handler.DeactivateEmployeeAirlineInfo())

		protected.POST("/messages", validator.WithValidateMessage(), handler.CreateMessage())

		protected.PUT("/messages/:id", validator.WithValidateMessage(), handler.UpdateMessage())

		protected.DELETE("/messages/:id", handler.DeleteMessage())

		protected.GET("/messages/:id", handler.GetMessageByID())

		protected.GET("/messages", handler.ListMessages())

		protected.POST("/messages/cache/reload", handler.ReloadMessageCache())

		protected.PATCH("/airlines/:id/activate", handler.ActivateAirline())

		protected.PATCH("/airlines/:id/deactivate", handler.DeactivateAirline())

		protected.PATCH("/airline-routes/:id/activate", handler.ActivateAirlineRoute())

		protected.PATCH("/airline-routes/:id/deactivate", handler.DeactivateAirlineRoute())

		protected.GET("/employees/airline-routes", handler.ListMyAirlineRoutes())

		protected.PATCH("/airports/:id/deactivate", handler.DeactivateAirport())

		protected.PATCH("/airports/:id/activate", handler.ActivateAirport())

		protected.PATCH("/aircraft-models/:id/activate", handler.ActivateAircraftModel())

		protected.PATCH("/aircraft-models/:id/deactivate", handler.DeactivateAircraftModel())

		protected.GET("/license-plates", handler.ListLicensePlates())

		protected.GET("/license-plates/:plate", handler.GetLicensePlateByPlate())

		protected.POST("/license-plates", validator.WithValidateCreateLicensePlate(), handler.CreateLicensePlate())

		protected.PUT("/license-plates/:id", validator.WithValidateUpdateLicensePlate(), handler.UpdateLicensePlate())

		protected.GET("/daily-logbook-details/:id", handler.GetDailyLogbookDetail())

		protected.PUT("/daily-logbook-details/:id", validator.WithValidateUpdateDailyLogbookDetail(), handler.UpdateDailyLogbookDetail())

		protected.GET("/daily-logbooks/:id/details", handler.ListDailyLogbookDetails())

		protected.POST("/daily-logbooks/:id/details", validator.WithValidateCreateDailyLogbookDetail(), handler.CreateDailyLogbookDetail())

		protected.GET("/employees/flights", handler.ListMyFlights())

		protected.GET("/daily-logbooks", handler.ListDailyLogbooks())

		protected.POST("/daily-logbooks", validator.WithValidateCreateDailyLogbook(), handler.CreateDailyLogbook())

		protected.GET("/daily-logbooks/:id", handler.GetDailyLogbookByID())

		protected.PATCH("/daily-logbooks/:id/activate", handler.ActivateDailyLogbook())

		protected.PATCH("/daily-logbooks/:id/deactivate", handler.DeactivateDailyLogbook())

	}
	admin := app.Group("flighthours/api/v1/admin")
	admin.Use(middleware.RequireAuth(dependencies.EmployeeService, dependencies.MessagingCache, dependencies.JWTValidator))
	admin.Use(middleware.RequireRole("admin"))
	{
		admin.GET("/routes", handler.ListRoutes())
		admin.GET("/airlines", handler.ListAirlines())
		admin.GET("/airlines/:id", handler.GetAirlineByID())
		admin.PATCH("/airlines/:id/activate", handler.ActivateAirline())
		admin.PATCH("/airlines/:id/deactivate", handler.DeactivateAirline())
		admin.PATCH("/airline-routes/:id/activate", handler.ActivateAirlineRoute())
		admin.PATCH("/airline-routes/:id/deactivate", handler.DeactivateAirlineRoute())
	}
	log.Success(logger.LogRouteConfigured)
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
