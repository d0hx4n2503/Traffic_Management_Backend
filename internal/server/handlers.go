package server

import (
	"net/http"
	"os"
	"strings"

	_ "github.com/adohong4/driving-license/docs"
	echoSwagger "github.com/swaggo/echo-swagger"

	authHttp "github.com/adohong4/driving-license/internal/auth/delivery/http"
	authRepository "github.com/adohong4/driving-license/internal/auth/repository"
	authUseCase "github.com/adohong4/driving-license/internal/auth/usecase"

	govAgencyHttp "github.com/adohong4/driving-license/internal/gov_agency/delivery/http"
	govAgencyRepo "github.com/adohong4/driving-license/internal/gov_agency/repository"
	govAgencyUC "github.com/adohong4/driving-license/internal/gov_agency/usecase"

	driverLicenseHttp "github.com/adohong4/driving-license/internal/driver_license/delivery/http"
	driverLicenseRepo "github.com/adohong4/driving-license/internal/driver_license/repository"
	driverLicenseUseCase "github.com/adohong4/driving-license/internal/driver_license/usecase"

	vehicleRegHttp "github.com/adohong4/driving-license/internal/vehicle_registration/delivery/http"
	vehicleReqRepository "github.com/adohong4/driving-license/internal/vehicle_registration/repository"
	vehicleReqUseCase "github.com/adohong4/driving-license/internal/vehicle_registration/usecase"

	trafficVioHttp "github.com/adohong4/driving-license/internal/traffic_violation/delivery/http"
	trafficVioRepository "github.com/adohong4/driving-license/internal/traffic_violation/repository"
	trafficVioUseCase "github.com/adohong4/driving-license/internal/traffic_violation/usecase"

	newsHttp "github.com/adohong4/driving-license/internal/news/delivery/http"
	newsRepository "github.com/adohong4/driving-license/internal/news/repository"
	newsUseCase "github.com/adohong4/driving-license/internal/news/usecase"

	notiHttp "github.com/adohong4/driving-license/internal/notification/delivery/http"
	notiRepository "github.com/adohong4/driving-license/internal/notification/repository"
	notiUseCase "github.com/adohong4/driving-license/internal/notification/usecase"

	apiMiddlewares "github.com/adohong4/driving-license/internal/middleware"
	"github.com/adohong4/driving-license/pkg/utils"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Map Server Handler
func (s *Server) MapHandlers(e *echo.Echo) error {
	// Init Repositories
	aRepo := authRepository.NewAuthRepository(s.db)
	gRepo := govAgencyRepo.NewGovAgencyRepo(s.db)
	dRepo := driverLicenseRepo.NewDriverLicenseRepo(s.db)
	vReRepo := vehicleReqRepository.NewVehicleDocRepository(s.db)
	tRepo := trafficVioRepository.NewTrafficViolationRepo(s.db)
	newsRepo := newsRepository.NewNewsRepo(s.db)
	notiRepo := notiRepository.NewNotificationRepo(s.db)

	// Init Usecase
	authUC := authUseCase.NewAuthUseCase(s.cfg, aRepo, s.logger)
	goAgenUC := govAgencyUC.NewGovAgencyUseCase(s.cfg, gRepo, s.logger)
	dlUC := driverLicenseUseCase.NewDriverLicenseUseCase(s.cfg, dRepo, s.logger)
	vReUC := vehicleReqUseCase.NewVehicleRegUseCase(s.cfg, vReRepo, s.logger)
	tUC := trafficVioUseCase.NewTrafficViolationUseCase(s.cfg, tRepo, s.logger)
	newsUC := newsUseCase.NewNewsUseCase(s.cfg, newsRepo, s.logger)
	notiUC := notiUseCase.NewNotificationUseCase(s.cfg, notiRepo, s.logger)

	// Init Handler
	authHandlers := authHttp.NewAuthHandlers(s.cfg, authUC, s.logger)
	govAgencyHandlers := govAgencyHttp.NewGovAgencyHandlers(s.cfg, goAgenUC, s.logger)
	driverLicenseHandlers := driverLicenseHttp.NewDriverLicenseHandlers(s.cfg, dlUC, s.logger)
	vehiclerReqHandlers := vehicleRegHttp.NewVehicleReqHandlers(s.cfg, vReUC, s.logger)
	trafficVioHandlers := trafficVioHttp.NewTrafficViolationHandlers(s.cfg, tUC, s.logger)
	newsHandlers := newsHttp.NewsHandlers(s.cfg, newsUC, s.logger)
	notiHandlers := notiHttp.NewNotificationHandlers(s.cfg, notiUC, s.logger)

	mw := apiMiddlewares.NewMiddlewareManager(authUC, s.cfg, []string{"*"}, s.logger)

	// middleware
	e.Use(mw.RequestLoggerMiddleware)
	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())

	allowOrigins := []string{"http://localhost:3000", "http://localhost:3001"}
	if corsOrigins := os.Getenv("CORS_ALLOW_ORIGINS"); corsOrigins != "" {
		parsedOrigins := strings.Split(corsOrigins, ",")
		allowOrigins = make([]string, 0, len(parsedOrigins))
		for _, origin := range parsedOrigins {
			origin = strings.TrimSpace(origin)
			if origin != "" {
				allowOrigins = append(allowOrigins, origin)
			}
		}
	}

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     allowOrigins,
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions, http.MethodPatch, http.MethodHead},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization, "X-CSRF-Token"}, // Add "X-CSRF-Token" so you have CSRF middleware
		ExposeHeaders:    []string{"Content-Length", "X-CSRF-Token"},                                                                       // If need expose add header
		AllowCredentials: true,                                                                                                             // Acceptance send cookie/credentials
		MaxAge:           86400,                                                                                                            // 24 giờ,
	}))

	//Swagger
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// CSRF middleware
	if s.cfg.Server.CSRF {
		//if false {
		e.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
			TokenLookup:    "header:X-CSRF-Token",
			CookieName:     s.cfg.Cookie.Name,
			CookieMaxAge:   s.cfg.Cookie.MaxAge,
			CookieSecure:   s.cfg.Cookie.Secure,
			CookieHTTPOnly: s.cfg.Cookie.HTTPOnly,
		}))
	}

	// API v1
	v1 := e.Group("v1/api")

	health := v1.Group("/health")
	authGroup := v1.Group("/auth")
	goAgencyGroup := v1.Group("/agency")
	driverLicenseGroup := v1.Group("/licenses")
	vehicleReqGroup := v1.Group("/vehicle")
	trafficVioGroup := v1.Group("/traffic")
	newsGroup := v1.Group("/news")
	notiGroup := v1.Group("/noti")

	authHttp.MapAuthRoutes(authGroup, authHandlers, mw, s.cfg, authUC)
	govAgencyHttp.MapGovAgencyRoutes(goAgencyGroup, govAgencyHandlers)
	driverLicenseHttp.MapDriverLicenseRoutes(driverLicenseGroup, driverLicenseHandlers, mw, s.cfg, authUC)
	vehicleRegHttp.MapVehicleRegistrationRoutes(vehicleReqGroup, vehiclerReqHandlers, mw, s.cfg, authUC)
	trafficVioHttp.MapTrafficViolationRoutes(trafficVioGroup, trafficVioHandlers, mw, s.cfg, authUC)
	newsHttp.MapNewsRoutes(newsGroup, newsHandlers, mw, authUC, s.cfg)
	notiHttp.MapNotificationRoutes(notiGroup, notiHandlers, mw, s.cfg, authUC)

	health.GET("", func(c echo.Context) error {
		s.logger.Infof("Health check request id: %s", utils.GetRequestId(c))
		return c.JSON(http.StatusOK, map[string]string{"status": "OK"})
	})

	return nil
}
