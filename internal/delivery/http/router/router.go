package router

import (
	"log"
	"time"

	"github.com/alpinnz/go-rest-api-boilerplate/internal/container"
	"github.com/alpinnz/go-rest-api-boilerplate/internal/localization"
	"github.com/alpinnz/go-rest-api-boilerplate/internal/middleware"
	"github.com/gin-gonic/gin"
)

type Router struct {
	engine    *gin.Engine
	container *container.Container
}

func NewRouter(c *container.Container) *Router {
	return &Router{
		engine:    gin.Default(),
		container: c,
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
	r.engine.Use(middleware.Logger(r.container.Logger))
	r.engine.Use(middleware.Language(bundle))
	r.engine.Use(middleware.CORS(r.container.Config.CORS))
	r.engine.Use(middleware.Recovery(r.container.Logger))
	r.engine.Use(middleware.Sanitize())
	r.engine.Use(middleware.Timeout(30 * time.Second))

	// Serve OpenAPI specification
	r.engine.StaticFile("/docs/swagger.json", "./internal/delivery/http/docs/swagger.json")
	r.engine.StaticFile("/docs", "./internal/delivery/http/docs/swagger.html")
	r.engine.StaticFile("/docs/", "./internal/delivery/http/docs/swagger.html")

	v1 := r.engine.Group("/api/v1")

	{
		// Health check endpoint (no rate limit)
		health := v1.Group("/health")
		{
			health.GET("", r.container.HealthHandler.Check)
			health.GET("/live", r.container.HealthHandler.Liveness)
			health.GET("/ready", r.container.HealthHandler.Readiness)
		}

		// Rate limiter for authentication endpoints (prevent brute force)
		var authLimiter *middleware.RateLimiter
		if r.container.Config.RateLimiter.Enabled {
			authLimiter = middleware.NewRateLimiter(
				r.container.Config.RateLimiter.RequestsLimit,
				time.Duration(r.container.Config.RateLimiter.WindowMinutes)*time.Minute,
			)
		}

		// Authentication endpoints
		auth := v1.Group("/auth")
		if authLimiter != nil {
			auth.Use(authLimiter.Limit())
		}
		{
			auth.POST("/register", r.container.AuthHandler.Register)
			auth.POST("/login", r.container.AuthHandler.Login)
			auth.POST("/refresh-token", r.container.AuthHandler.RefreshToken)
			auth.POST("/logout", r.container.AuthMiddleware.Authenticate(), r.container.AuthHandler.Logout)
		}

		// User management endpoints
		users := v1.Group("/users")
		users.Use(r.container.AuthMiddleware.Authenticate())
		{
			users.GET("/me", r.container.UserHandler.GetProfile)
			users.GET("", r.container.UserHandler.List)
			users.GET("/:id", r.container.UserHandler.GetByID)
			users.PUT("/:id", r.container.UserHandler.Update)
			users.DELETE("/:id", r.container.UserHandler.Delete)
			users.POST("/:id/roles", r.container.UserHandler.AssignRole)
			users.DELETE("/:id/roles/:roleId", r.container.UserHandler.RemoveRole)
		}

		// Role management endpoints
		roles := v1.Group("/roles")
		roles.Use(r.container.AuthMiddleware.Authenticate())
		{
			roles.GET("", r.container.RoleHandler.List)
			roles.POST("", r.container.RoleHandler.Create)
			roles.GET("/:id", r.container.RoleHandler.GetByID)
			roles.PUT("/:id", r.container.RoleHandler.Update)
			roles.DELETE("/:id", r.container.RoleHandler.Delete)
		}

		// WebSocket endpoints
		ws := v1.Group("/ws")
		{
			ws.GET("", r.container.AuthMiddleware.Authenticate(), r.container.WebSocketHandler.Connect)
			ws.GET("/stats", r.container.AuthMiddleware.Authenticate(), r.container.WebSocketHandler.GetStats)
		}
	}

	return r.engine
}
