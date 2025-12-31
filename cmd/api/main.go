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
	"github.com/alpinnz/go-rest-api-boilerplate/pkg/auth"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	auth.SetJWTSecret(cfg.JWT.Secret)

	db, err := database.NewPostgresDB(cfg.Database.DSN())
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
	sessionRepo := repository.NewSessionRepository(redisClient)

	userUseCase := usecase.NewUserUseCase(
		userRepo,
		sessionRepo,
		cfg.JWT.Expiration,
		10*time.Second,
	)

	userHandler := handler.NewUserHandler(userUseCase)
	healthHandler := handler.NewHealthHandler()
	authMiddleware := middleware.NewAuthMiddleware(userUseCase)

	r := router.NewRouter(userHandler, healthHandler, authMiddleware)
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
