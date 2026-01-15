package container

import (
	"database/sql"
	"time"

	"github.com/alpinnz/go-rest-api-boilerplate/config"
	"github.com/alpinnz/go-rest-api-boilerplate/internal/delivery/http/handler"
	"github.com/alpinnz/go-rest-api-boilerplate/internal/domain/repository"
	"github.com/alpinnz/go-rest-api-boilerplate/internal/infrastructure/database"
	"github.com/alpinnz/go-rest-api-boilerplate/internal/middleware"
	repoImpl "github.com/alpinnz/go-rest-api-boilerplate/internal/repository"
	"github.com/alpinnz/go-rest-api-boilerplate/internal/usecase"
	"github.com/alpinnz/go-rest-api-boilerplate/internal/websocket"
	"github.com/alpinnz/go-rest-api-boilerplate/pkg/auth"
	"github.com/alpinnz/go-rest-api-boilerplate/pkg/logger"
	"github.com/redis/go-redis/v9"
)

// Container holds all application dependencies
type Container struct {
	Config *config.Config
	Logger *logger.Logger

	// Infrastructure
	DB          *sql.DB
	RedisClient *redis.Client

	// Repositories
	UserRepo    repository.UserRepository
	RoleRepo    repository.RoleRepository
	SessionRepo repository.SessionRepository

	// Use Cases
	AuthUseCase      *usecase.AuthUseCase
	UserUseCase      *usecase.UserUseCase
	RoleUseCase      *usecase.RoleUseCase
	WebSocketUseCase *usecase.WebSocketUseCase

	// Handlers
	AuthHandler      *handler.AuthHandler
	UserHandler      *handler.UserHandler
	RoleHandler      *handler.RoleHandler
	HealthHandler    *handler.HealthHandler
	WebSocketHandler *handler.WebSocketHandler

	// Middleware
	AuthMiddleware *middleware.AuthMiddleware

	// WebSocket
	WSHub *websocket.Hub
}

// New creates and initializes a new container with all dependencies
func New(cfg *config.Config) (*Container, error) {
	c := &Container{
		Config: cfg,
	}

	// Initialize Logger
	c.Logger = logger.New(logger.Config{
		Level:  cfg.Logger.Level,
		Pretty: cfg.Logger.Pretty,
	})

	// Set JWT config
	auth.SetJWTConfig(
		cfg.JWT.AccessTokenSecret,
		cfg.JWT.AccessTokenExpiration,
		cfg.JWT.RefreshTokenSecret,
		cfg.JWT.RefreshTokenExpiration,
	)

	// Initialize Infrastructure
	if err := c.initInfrastructure(); err != nil {
		return nil, err
	}

	// Initialize Repositories
	c.initRepositories()

	// Initialize Use Cases
	c.initUseCases()

	// Initialize WebSocket Hub
	c.WSHub = websocket.NewHub()
	go c.WSHub.Run()

	// Initialize Handlers
	c.initHandlers()

	// Initialize Middleware
	c.initMiddleware()

	return c, nil
}

// initInfrastructure initializes database and redis connections
func (c *Container) initInfrastructure() error {
	poolConfig := database.PoolConfig{
		MaxOpenConns:    c.Config.Database.MaxOpenConns,
		MaxIdleConns:    c.Config.Database.MaxIdleConns,
		ConnMaxLifetime: c.Config.Database.ConnMaxLifetime,
		ConnMaxIdleTime: c.Config.Database.ConnMaxIdleTime,
	}

	db, err := database.NewPostgresDB(c.Config.Database.DSN(), poolConfig)
	if err != nil {
		return err
	}
	c.DB = db

	redisClient, err := database.NewRedisClient(
		c.Config.Redis.Addr(),
		c.Config.Redis.Password,
		c.Config.Redis.DB,
	)
	if err != nil {
		return err
	}
	c.RedisClient = redisClient

	return nil
}

// initRepositories initializes all repositories
func (c *Container) initRepositories() {
	c.UserRepo = repoImpl.NewUserRepository(c.DB)
	c.RoleRepo = repoImpl.NewRoleRepository(c.DB)
	c.SessionRepo = repoImpl.NewSessionRepository(c.RedisClient)
}

// initUseCases initializes all use cases
func (c *Container) initUseCases() {
	timeout := 10 * time.Second

	c.AuthUseCase = usecase.NewAuthUseCase(
		c.UserRepo,
		c.SessionRepo,
		timeout,
	)

	c.UserUseCase = usecase.NewUserUseCase(
		c.UserRepo,
		c.RoleRepo,
		timeout,
	)

	c.RoleUseCase = usecase.NewRoleUseCase(
		c.RoleRepo,
		timeout,
	)

	c.WebSocketUseCase = usecase.NewWebSocketUseCase(c.WSHub)
}

// initHandlers initializes all HTTP handlers
func (c *Container) initHandlers() {
	c.AuthHandler = handler.NewAuthHandler(c.AuthUseCase)
	c.UserHandler = handler.NewUserHandler(c.UserUseCase)
	c.RoleHandler = handler.NewRoleHandler(c.RoleUseCase)
	c.HealthHandler = handler.NewHealthHandler(c.DB, c.RedisClient)
	c.WebSocketHandler = handler.NewWebSocketHandler(c.WSHub, c.Config.CORS)
}

// initMiddleware initializes all middleware
func (c *Container) initMiddleware() {
	c.AuthMiddleware = middleware.NewAuthMiddleware(c.AuthUseCase)
}

// Close closes all resources
func (c *Container) Close() error {
	c.Logger.Info().Msg("Closing resources...")

	if c.DB != nil {
		if err := c.DB.Close(); err != nil {
			c.Logger.Error().Err(err).Msg("Failed to close database")
			return err
		}
	}

	if c.RedisClient != nil {
		if err := c.RedisClient.Close(); err != nil {
			c.Logger.Error().Err(err).Msg("Failed to close Redis client")
			return err
		}
	}

	c.Logger.Info().Msg("All resources closed successfully")
	return nil
}
