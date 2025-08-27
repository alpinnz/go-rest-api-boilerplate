package server

import (
	"net/http"

	"github.com/alpinnz/go-rest-api-boilerplate/internal/config"
	"github.com/alpinnz/go-rest-api-boilerplate/internal/infrastructure/repositories"
	"github.com/alpinnz/go-rest-api-boilerplate/internal/interfaces/http/controllers"
	"github.com/alpinnz/go-rest-api-boilerplate/internal/interfaces/http/middleware"
	"github.com/alpinnz/go-rest-api-boilerplate/internal/interfaces/http/router"
	"github.com/alpinnz/go-rest-api-boilerplate/internal/usecase"
	"github.com/alpinnz/go-rest-api-boilerplate/pkg/validation"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewHTTPServer(env *config.Env, db *gorm.DB) (*http.Server, error) {
	// Init validator
	validator, err := validation.NewValidator()
	if err != nil {
		return nil, err
	}

	// Init repositories
	authSessionRepo := repositories.NewAuthSessionRepository()
	userRepo := repositories.NewUserRepository()
	roleRepo := repositories.NewRoleRepository()
	userRoleRepo := repositories.NewUserRoleRepository()

	// Init usecases
	authUsecase := usecase.NewAuthUsecase(env, db, userRepo, authSessionRepo, roleRepo, userRoleRepo)
	userUsecase := usecase.NewUserUsecase(env, db, userRepo, roleRepo, userRoleRepo)
	roleUsecase := usecase.NewRoleUsecase(env, db, roleRepo)

	// Init controllers
	authController := controllers.NewAuthController(validator, authUsecase)
	userController := controllers.NewUserController(validator, userUsecase)
	roleController := controllers.NewRoleController(validator, roleUsecase)

	// Gin engine
	engine := gin.New()
	engine.Use(gin.Logger())
	engine.Use(gin.Recovery())
	_ = engine.SetTrustedProxies(nil)
	engine.Use(middleware.RateLimiterMiddleware())
	engine.Use(middleware.CorsMiddleware())

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
