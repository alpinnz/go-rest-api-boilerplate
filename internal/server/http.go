package server

import (
	"net/http"

	"github.com/alpinnz/go-rest-api-boilerplate/internal/config"
	"github.com/alpinnz/go-rest-api-boilerplate/internal/infrastructure/repositories"
	"github.com/alpinnz/go-rest-api-boilerplate/internal/interfaces/http/controllers"
	"github.com/alpinnz/go-rest-api-boilerplate/internal/interfaces/http/middleware"
	"github.com/alpinnz/go-rest-api-boilerplate/internal/interfaces/http/router"
	"github.com/alpinnz/go-rest-api-boilerplate/internal/usecase"
	"github.com/alpinnz/go-rest-api-boilerplate/pkg/translations"
	"github.com/alpinnz/go-rest-api-boilerplate/pkg/validation"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewHTTPServer(env *config.Env, db *gorm.DB) (*http.Server, error) {
	// Load translations
	tr := translations.NewStore()
	if err := tr.LoadAllModules("locales"); err != nil {
		return nil, err
	}

	// Init validator
	validator, err := validation.NewValidator(tr)
	if err != nil {
		return nil, err
	}

	// Init repositories
	authSessionRepo := repositories.NewAuthSessionRepository()
	userRepo := repositories.NewUserRepository()
	roleRepo := repositories.NewRoleRepository()
	userRoleRepo := repositories.NewUserRoleRepository()

	// Init usecases
	authUsecase := usecase.NewAuthUsecase(env, db, tr, userRepo, authSessionRepo, roleRepo, userRoleRepo)
	userUsecase := usecase.NewUserUsecase(env, db, tr, userRepo, roleRepo, userRoleRepo)
	roleUsecase := usecase.NewRoleUsecase(env, db, tr, roleRepo)

	// Init controllers
	authController := controllers.NewAuthController(validator, tr, authUsecase)
	userController := controllers.NewUserController(validator, tr, userUsecase)
	roleController := controllers.NewRoleController(validator, tr, roleUsecase)

	// Gin engine
	engine := gin.New()
	engine.Use(gin.Logger())
	engine.Use(gin.Recovery())
	_ = engine.SetTrustedProxies(nil)
	engine.Use(middleware.CorsMiddleware())
	engine.Use(middleware.RateLimiterMiddleware())
	engine.Use(translations.Middleware(translations.Locale(env.App.DefaultLocale)))
	engine.Use(middleware.LocateMiddleware())
	engine.Use(middleware.DeviceMiddleware())
	engine.Use(middleware.PlatformMiddleware())

	// Register routes
	api := engine.Group("/api")
	v1 := api.Group("/v1")
	router.NewAuthRoutes(v1, authController)
	router.NewUserRoutes(v1, userController)
	router.NewRoleRoutes(v1, roleController)

	return &http.Server{
		Addr:    env.Server.Address,
		Handler: engine,
	}, nil
}
