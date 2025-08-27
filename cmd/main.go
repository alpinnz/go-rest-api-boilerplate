package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/alpinnz/go-rest-api-boilerplate/internal/config"
	"github.com/alpinnz/go-rest-api-boilerplate/internal/infrastructure/databases"
	"github.com/alpinnz/go-rest-api-boilerplate/internal/server"
	"github.com/alpinnz/go-rest-api-boilerplate/pkg/app"
)

func main() {
	env := config.NewEnv()
	fmt.Println("Connecting to DB at:", env.Postgres.Host)

	// Connect DB
	postgres, err := databases.NewPostgres(env.Postgres)
	if err != nil {
		log.Fatal("Postgres connection error:", err)
	}

	// App context with signal handling
	ctx, cancel := app.NewAppContext()
	defer cancel()

	// Migration & seeders
	if err := databases.Migrate(postgres); err != nil {
		log.Fatal(err)
	}
	if err := databases.SeedRoles(ctx, postgres); err != nil {
		log.Fatal(err)
	}
	if err := databases.SeedUsers(ctx, postgres, env.Auth.PasswordSecret); err != nil {
		log.Fatal(err)
	}

	// Build HTTP server
	httpServer, err := server.NewHTTPServer(&env, postgres)
	if err != nil {
		log.Fatal("Failed to init HTTP server:", err)
	}

	// Run server
	go func() {
		log.Println("Server started on", env.Server.Address)
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server error: %s\n", err)
		}
	}()

	// Wait for cancel (SIGINT/SIGTERM)
	<-ctx.Done()
	log.Println("Shutting down server...")

	// Graceful shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server forced to shutdown: %s", err)
	}
	log.Println("Server exited gracefully")
}
