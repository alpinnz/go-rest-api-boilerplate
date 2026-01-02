package router

import (
	"log"

	"github.com/alpinnz/go-rest-api-boilerplate/internal/delivery/http/handler"
	"github.com/alpinnz/go-rest-api-boilerplate/internal/localization"
	"github.com/alpinnz/go-rest-api-boilerplate/internal/middleware"
	"github.com/gin-gonic/gin"
)

type Router struct {
	engine         *gin.Engine
	authHandler    *handler.AuthHandler
	userHandler    *handler.UserHandler
	roleHandler    *handler.RoleHandler
	healthHandler  *handler.HealthHandler
	authMiddleware *middleware.AuthMiddleware
}

func NewRouter(
	authHandler *handler.AuthHandler,
	userHandler *handler.UserHandler,
	roleHandler *handler.RoleHandler,
	healthHandler *handler.HealthHandler,
	authMiddleware *middleware.AuthMiddleware,
) *Router {
	return &Router{
		engine:         gin.Default(),
		authHandler:    authHandler,
		userHandler:    userHandler,
		roleHandler:    roleHandler,
		healthHandler:  healthHandler,
		authMiddleware: authMiddleware,
	}
}

func (r *Router) Setup() *gin.Engine {
	// Load localization bundles
	bundle := localization.NewBundle("en")
	if err := bundle.Load("en", "internal/localization/lang/en.json"); err != nil {
		log.Printf("Warning: failed to load en.json: %v", err)
	}
	if err := bundle.Load("id", "internal/localization/lang/id.json"); err != nil {
		log.Printf("Warning: failed to load id.json: %v", err)
	}

	// Global middlewares
	r.engine.Use(middleware.Logger())
	r.engine.Use(middleware.Language(bundle))
	r.engine.Use(middleware.CORS())
	r.engine.Use(middleware.Recovery())

	// Serve OpenAPI specification
	r.engine.StaticFile("/docs/swagger.json", "./docs/swagger.json")

	v1 := r.engine.Group("/api/v1")

	{
		// Health check endpoint
		health := v1.Group("/health")
		{
			health.GET("", r.healthHandler.Check)
		}

		// Authentication endpoints
		auth := v1.Group("/auth")
		{
			auth.POST("/register", r.authHandler.Register)
			auth.POST("/login", r.authHandler.Login)
			auth.POST("/refresh-token", r.authHandler.RefreshToken)
			auth.POST("/logout", r.authMiddleware.Authenticate(), r.authHandler.Logout)
		}

		// User management endpoints
		users := v1.Group("/users")
		users.Use(r.authMiddleware.Authenticate())
		{
			users.GET("/me", r.userHandler.GetProfile)
			users.GET("", r.userHandler.List)
			users.GET("/:id", r.userHandler.GetByID)
			users.PUT("/:id", r.userHandler.Update)
			users.DELETE("/:id", r.userHandler.Delete)
			users.POST("/:id/roles", r.userHandler.AssignRole)
			users.DELETE("/:id/roles/:roleId", r.userHandler.RemoveRole)
		}

		// Role management endpoints
		roles := v1.Group("/roles")
		roles.Use(r.authMiddleware.Authenticate())
		{
			roles.GET("", r.roleHandler.List)
			roles.POST("", r.roleHandler.Create)
			roles.GET("/:id", r.roleHandler.GetByID)
			roles.PUT("/:id", r.roleHandler.Update)
			roles.DELETE("/:id", r.roleHandler.Delete)
		}
	}

	return r.engine
}
