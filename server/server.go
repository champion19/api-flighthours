package server

import (
	"log/slog"
	"strings"
	"time"

	"github.com/champion19/api-flighthours/cmd/dependency"
	"github.com/champion19/api-flighthours/handlers"
	"github.com/champion19/api-flighthours/middleware"
	"github.com/champion19/api-flighthours/platform/cookie"
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
		AllowOriginFunc: func(origin string) bool {
			// In development: allow any localhost port (Flutter Web, browser, etc.)
			return strings.HasPrefix(origin, "http://localhost") ||
				strings.HasPrefix(origin, "http://127.0.0.1")
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Request-ID"},
		ExposeHeaders:    []string{"Content-Length", "Location", "Set-Cookie"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	app.Use(cors.New(corsConfig))
	log.Info(logger.LogAppCORSConfigured)

	app.GET("/metrics", gin.WrapH(promhttp.Handler()))

	app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	app.Use(middleware.RequestID())

	app.Use(middleware.TrackMetrics())

	errorHandler := middleware.NewErrorHandler(dependencies.MessagingCache)
	app.Use(errorHandler.Handle())

	cookieManager := cookie.NewManager(dependencies.Config.Cookie)

	handler := handlers.New(handlers.HandlerDeps{
		Service:                      dependencies.EmployeeService,
		EmployeeInteractor:           dependencies.Interactor,
		IDEncoder:                    dependencies.IDEncoder,
		Response:                     dependencies.ResponseHandler,
		CookieManager:                cookieManager,
		MessageInteractor:            dependencies.MessageInteractor,
		MessagingCache:               dependencies.MessagingCache,
		AirlineInteractor:            dependencies.AirlineInteractor,
		AirlineEmployeeInteractor:    dependencies.AirlineEmployeeInteractor,
		EngineInteractor:             dependencies.EngineInteractor,
		RouteInteractor:              dependencies.RouteInteractor,
		ManufacturerInteractor:       dependencies.ManufacturerInteractor,
		AirportInteractor:            dependencies.AirportInteractor,
		AirlineRouteInteractor:       dependencies.AirlineRouteInteractor,
		DailyLogbookDetailInteractor: dependencies.DailyLogbookDetailInteractor,
		DailyLogbookInteractor:       dependencies.DailyLogbookInteractor,
		AircraftModelInteractor:      dependencies.AircraftModelInteractor,
		TailNumberInteractor:         dependencies.TailNumberInteractor,
		FlightSummaryInteractor:      dependencies.FlightSummaryInteractor,
	})

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
	authLimiter := middleware.NewRateLimiter(0.2, 5) // 1 req/5s, burst 5
	{

		public.POST("/register", validator.WithValidateRegister(), handler.RegisterEmployee())

		public.POST("/login", authLimiter.Limit(), validator.WithValidateLogin(), handler.Login())

		public.POST("/auth/resend-verification", validator.WithValidateResendVerificationEmail(), handler.ResendVerificationEmail())

		public.POST("/auth/verify-email", validator.WithValidateVerifyEmail(), handler.VerifyEmailByToken())

		public.POST("/auth/password-reset", authLimiter.Limit(), validator.WithValidatePasswordResetRequest(), handler.PasswordReset())

		public.POST("/auth/update-password", authLimiter.Limit(), validator.WithValidateUpdatePassword(), handler.UpdatePassword())

		public.POST("/auth/refresh", authLimiter.Limit(), validator.WithValidateRefreshToken(), handler.RefreshToken())

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

		public.GET("/crew-member-types", handler.GetCrewMemberTypes())

		public.GET("/manufacturers", handler.ListManufacturers())

		public.GET("/manufacturers/:id", handler.GetManufacturerByID())

		public.GET("/aircraft-models", handler.ListAircraftModels())

		public.GET("/aircraft-models/:id", handler.GetAircraftModelByID())

		public.GET("/aircraft-families/:family", handler.GetAircraftModelsByFamily())

	}
	protected := app.Group("flighthours/api/v1")
	protected.Use(middleware.RequireAuth(dependencies.EmployeeService, dependencies.MessagingCache, dependencies.JWTValidator))

	// ── Authenticated endpoints (any role) ─(Employee)─────────────────────────────
	protected.GET("/employees", handler.GetEmployee())

	// ── Pilot-only endpoints (role: "pilot") ────────────────────────────
	pilot := protected.Group("")
	pilot.Use(middleware.RequireRole("pilot"))
	{
		// Employee management (HU17, HU18, HU20, HU22, HU23)
		pilot.POST("/auth/change-password", validator.WithValidateChangePassword(), handler.ChangePassword())
		pilot.PUT("/employees", validator.WithValidateUpdateEmployee(), handler.UpdateEmployee())
		pilot.DELETE("/employees", handler.DeleteEmployee())

		// Airline-Employee management (HU24, HU25, HU26, HU27, HU28)
		pilot.GET("/employees/airline", handler.GetEmployeeAirlineInfo())
		pilot.PUT("/employees/airline", validator.WithValidateAddAirlineEmployee(), handler.AddEmployeeAirlineInfo())
		pilot.PUT("/employees/airline-info", validator.WithValidateUpdateAirlineEmployee(), handler.UpdateEmployeeAirlineInfo())
		pilot.PATCH("/employees/airline/activate", handler.ActivateEmployeeAirlineInfo())
		pilot.PATCH("/employees/airline/deactivate", handler.DeactivateEmployeeAirlineInfo())

		// Airline routes consulta (HU37)
		pilot.GET("/employees/airline-routes", handler.ListMyAirlineRoutes())

		// Tail numbers consulta y creación (HU31, HU32)
		pilot.GET("/tail-numbers", handler.ListTailNumbers())
		pilot.GET("/tail-numbers/:plate", handler.GetTailNumberByPlate())
		pilot.POST("/tail-numbers", validator.WithValidateCreateTailNumber(), handler.CreateTailNumber())

		// Daily logbook (HU7, HU8, HU9, HU10, HU11, HU12)
		pilot.GET("/daily-logbooks", handler.ListDailyLogbooks())
		pilot.POST("/daily-logbooks", validator.WithValidateCreateDailyLogbook(), handler.CreateDailyLogbook())
		pilot.GET("/daily-logbooks/:id", handler.GetDailyLogbookByID())
		pilot.DELETE("/daily-logbooks/:id", handler.DeleteDailyLogbook())
		pilot.PUT("/daily-logbooks/:id", validator.WithValidateUpdateDailyLogbook(), handler.UpdateDailyLogbook())
		pilot.PATCH("/daily-logbooks/:id/activate", handler.ActivateDailyLogbook())
		pilot.PATCH("/daily-logbooks/:id/deactivate", handler.DeactivateDailyLogbook())

		// Daily logbook details (HU13, HU14, HU15, HU16)
		pilot.GET("/daily-logbook-details/:id", handler.GetDailyLogbookDetail())
		pilot.PUT("/daily-logbook-details/:id", validator.WithValidateUpdateDailyLogbookDetail(), handler.UpdateDailyLogbookDetail())
		pilot.DELETE("/daily-logbook-details/:id", handler.DeleteDailyLogbookDetail())
		pilot.GET("/daily-logbooks/:id/details", handler.ListDailyLogbookDetails())
		pilot.POST("/daily-logbooks/:id/details", validator.WithValidateCreateDailyLogbookDetail(), handler.CreateDailyLogbookDetail())

		// Flights (HU45, HU46, HU47)
		pilot.GET("/employees/flights", handler.ListMyFlights())
		pilot.GET("/employees/flight-hours-summary", handler.GetFlightHoursSummary())
		pilot.GET("/employees/flight-alerts", handler.GetFlightAlerts())
		pilot.GET("/employees/recent-flights", handler.GetRecentFlights())
	}

	// ── Admin-only endpoints (role: "admin") ────────────────────────────
	adminProtected := protected.Group("")
	adminProtected.Use(middleware.RequireRole("admin"))
	{
		// Messages (admin only)
		adminProtected.POST("/messages", validator.WithValidateMessage(), handler.CreateMessage())
		adminProtected.PUT("/messages/:id", validator.WithValidateMessage(), handler.UpdateMessage())
		adminProtected.DELETE("/messages/:id", handler.DeleteMessage())
		adminProtected.GET("/messages/:id", handler.GetMessageByID())
		adminProtected.GET("/messages", handler.ListMessages())
		adminProtected.POST("/messages/cache/reload", handler.ReloadMessageCache())

		// Airline state (HU2, HU3)
		adminProtected.PATCH("/airlines/:id/activate", handler.ActivateAirline())
		adminProtected.PATCH("/airlines/:id/deactivate", handler.DeactivateAirline())

		// Airline route state (HU38, HU39)
		adminProtected.PATCH("/airline-routes/:id/activate", handler.ActivateAirlineRoute())
		adminProtected.PATCH("/airline-routes/:id/deactivate", handler.DeactivateAirlineRoute())

		// Airport state (HU5, HU6)
		adminProtected.PATCH("/airports/:id/activate", handler.ActivateAirport())
		adminProtected.PATCH("/airports/:id/deactivate", handler.DeactivateAirport())

		// Aircraft model state (HU41, HU42)
		adminProtected.PATCH("/aircraft-models/:id/activate", handler.ActivateAircraftModel())
		adminProtected.PATCH("/aircraft-models/:id/deactivate", handler.DeactivateAircraftModel())

		// Tail number update (HU33)
		adminProtected.PUT("/tail-numbers/:id", validator.WithValidateUpdateTailNumber(), handler.UpdateTailNumber())
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
