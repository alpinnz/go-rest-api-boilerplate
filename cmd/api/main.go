package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alpinnz/go-rest-api-boilerplate/config"
	"github.com/alpinnz/go-rest-api-boilerplate/internal/delivery/http/handler"
	"github.com/alpinnz/go-rest-api-boilerplate/internal/delivery/http/router"
	"github.com/alpinnz/go-rest-api-boilerplate/internal/infrastructure/database"
	"github.com/alpinnz/go-rest-api-boilerplate/internal/middleware"
	"github.com/alpinnz/go-rest-api-boilerplate/internal/repository"
	"github.com/alpinnz/go-rest-api-boilerplate/internal/usecase"
	"github.com/alpinnz/go-rest-api-boilerplate/internal/websocket"
	"github.com/alpinnz/go-rest-api-boilerplate/pkg/auth"
)

// @title Go REST API Boilerplate
// @version 1.0
// @description Production-ready REST API with Clean Architecture, JWT authentication, and multi-language support
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@swagger.io

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1
// @schemes http https

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Bearer token authentication. Format: "Bearer {token}"

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	auth.SetJWTConfig(
		cfg.JWT.AccessTokenSecret,
		cfg.JWT.AccessTokenExpiration,
		cfg.JWT.RefreshTokenSecret,
		cfg.JWT.RefreshTokenExpiration,
	)

	poolConfig := database.PoolConfig{
		MaxOpenConns:    cfg.Database.MaxOpenConns,
		MaxIdleConns:    cfg.Database.MaxIdleConns,
		ConnMaxLifetime: cfg.Database.ConnMaxLifetime,
		ConnMaxIdleTime: cfg.Database.ConnMaxIdleTime,
	}

	db, err := database.NewPostgresDB(cfg.Database.DSN(), poolConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	redisClient, err := database.NewRedisClient(cfg.Redis.Addr(), cfg.Redis.Password, cfg.Redis.DB)
	if err != nil {
		log.Fatalf("Failed to connect to redis: %v", err)
	}
	defer redisClient.Close()

	userRepo := repository.NewUserRepository(db)
	roleRepo := repository.NewRoleRepository(db)
	sessionRepo := repository.NewSessionRepository(redisClient)

	authUseCase := usecase.NewAuthUseCase(
		userRepo,
		sessionRepo,
		10*time.Second,
	)

	userUseCase := usecase.NewUserUseCase(
		userRepo,
		roleRepo,
		10*time.Second,
	)

	roleUseCase := usecase.NewRoleUseCase(
		roleRepo,
		10*time.Second,
	)

	// Initialize WebSocket hub
	wsHub := websocket.NewHub()
	go wsHub.Run()

	// WebSocket use case can be used for broadcasting messages
	_ = usecase.NewWebSocketUseCase(wsHub)

	authHandler := handler.NewAuthHandler(authUseCase)
	userHandler := handler.NewUserHandler(userUseCase)
	roleHandler := handler.NewRoleHandler(roleUseCase)
	healthHandler := handler.NewHealthHandler(db, redisClient)
	websocketHandler := handler.NewWebSocketHandler(wsHub)
	authMiddleware := middleware.NewAuthMiddleware(authUseCase)

	r := router.NewRouter(authHandler, userHandler, roleHandler, healthHandler, websocketHandler, authMiddleware)
	engine := r.Setup()

	srv := &http.Server{
		Addr:         ":" + cfg.App.Port,
		Handler:      engine,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	go func() {
		fmt.Printf("Server starting on port %s\n", cfg.App.Port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	fmt.Println("Server exited")
}
