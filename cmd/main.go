package main

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/alpinnz/go-rest-api-boilerplate/internal/config"
	"github.com/alpinnz/go-rest-api-boilerplate/internal/infrastructure/databases"
	"github.com/alpinnz/go-rest-api-boilerplate/internal/server"
	"github.com/alpinnz/go-rest-api-boilerplate/pkg/app"
	"github.com/alpinnz/go-rest-api-boilerplate/pkg/logger"
)

func main() {
	// Init global logger
	logger.Init()

	env := config.NewEnv()
	logger.Log.Info("Connecting to DB", "host", env.Postgres.Host, "port", env.Postgres.Port)

	// Connect DB
	postgres, err := databases.NewPostgres(env.Postgres)
	if err != nil {
		logger.Log.Error("Postgres connection error", "error", err)
		return
	}

	// App context with signal handling
	ctx, cancel := app.NewAppContext()
	defer cancel()

	// Migration & seeders
	if err := databases.Migrate(postgres); err != nil {
		logger.Log.Error("Migration error", "error", err)
		return
	}
	if err := databases.SeedRoles(ctx, postgres); err != nil {
		logger.Log.Error("Seed roles error", "error", err)
		return
	}
	if err := databases.SeedUsers(ctx, postgres, env.Auth.PasswordSecret); err != nil {
		logger.Log.Error("Seed users error", "error", err)
		return
	}

	// Build HTTP server
	httpServer, err := server.NewHTTPServer(&env, postgres)
	if err != nil {
		logger.Log.Error("Failed to init HTTP server", "error", err)
		return
	}

	// Run server
	go func() {
		logger.Log.Info("Server started", "address", env.Server.Address)
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Log.Error("Server error", "error", err)
		}
	}()

	// Wait for cancel (SIGINT/SIGTERM)
	<-ctx.Done()
	logger.Log.Warn("Shutting down server...")

	// Graceful shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		logger.Log.Error("Server forced to shutdown", "error", err)
	}
	logger.Log.Info("Server exited gracefully")
}
