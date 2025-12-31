package router

import (
	"github.com/alpinnz/go-rest-api-boilerplate/internal/delivery/http/handler"
	"github.com/alpinnz/go-rest-api-boilerplate/internal/middleware"
	"github.com/gin-gonic/gin"
)

type Router struct {
	engine         *gin.Engine
	userHandler    *handler.UserHandler
	healthHandler  *handler.HealthHandler
	authMiddleware *middleware.AuthMiddleware
}

func NewRouter(
	userHandler *handler.UserHandler,
	healthHandler *handler.HealthHandler,
	authMiddleware *middleware.AuthMiddleware,
) *Router {
	return &Router{
		engine:         gin.Default(),
		userHandler:    userHandler,
		healthHandler:  healthHandler,
		authMiddleware: authMiddleware,
	}
}

func (r *Router) Setup() *gin.Engine {
	r.engine.Use(middleware.CORS())
	r.engine.Use(middleware.Logger())
	r.engine.Use(middleware.Recovery())

	v1 := r.engine.Group("/api/v1")
	{
		health := v1.Group("/health")
		{
			health.GET("", r.healthHandler.Check)
		}

		auth := v1.Group("/auth")
		{
			auth.POST("/register", r.userHandler.Register)
			auth.POST("/login", r.userHandler.Login)
			auth.POST("/logout", r.authMiddleware.Authenticate(), r.userHandler.Logout)
		}

		users := v1.Group("/users")
		users.Use(r.authMiddleware.Authenticate())
		{
			users.GET("/me", r.userHandler.GetProfile)
			users.GET("/:id", r.userHandler.GetByID)
		}
	}

	return r.engine
}
