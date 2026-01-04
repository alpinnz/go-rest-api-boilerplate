package handler

import (
	"context"
	"database/sql"
	"time"

	"github.com/alpinnz/go-rest-api-boilerplate/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type HealthHandler struct {
	db          *sql.DB
	redisClient *redis.Client
	startTime   time.Time
}

func NewHealthHandler(db *sql.DB, redisClient *redis.Client) *HealthHandler {
	return &HealthHandler{
		db:          db,
		redisClient: redisClient,
		startTime:   time.Now(),
	}
}

// HealthStatus represents the health status response
type HealthStatus struct {
	Status    string                 `json:"status"`
	Service   string                 `json:"service"`
	Timestamp string                 `json:"timestamp"`
	Uptime    string                 `json:"uptime"`
	Checks    map[string]CheckResult `json:"checks"`
}

// CheckResult represents individual check result
type CheckResult struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

// Check godoc
// @Summary      Health check
// @Description  Check if the API is running with all dependencies
// @Tags         Health
// @Accept       json
// @Produce      json
// @Success      200  {object}  response.Response{data=HealthStatus}
// @Failure      503  {object}  response.Response{data=HealthStatus}
// @Router       /health [get]
func (h *HealthHandler) Check(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	checks := make(map[string]CheckResult)
	overallStatus := "healthy"

	// Check database
	dbResult := h.checkDatabase(ctx)
	checks["database"] = dbResult
	if dbResult.Status != "healthy" {
		overallStatus = "unhealthy"
	}

	// Check Redis
	redisResult := h.checkRedis(ctx)
	checks["redis"] = redisResult
	if redisResult.Status != "healthy" {
		overallStatus = "unhealthy"
	}

	status := HealthStatus{
		Status:    overallStatus,
		Service:   "go-rest-api-boilerplate",
		Timestamp: time.Now().Format(time.RFC3339),
		Uptime:    time.Since(h.startTime).String(),
		Checks:    checks,
	}

	if overallStatus == "unhealthy" {
		response.ServiceUnavailable(c, status)
		return
	}

	response.Success(c, status)
}

// Liveness godoc
// @Summary      Liveness probe
// @Description  Check if the service is alive (for Kubernetes)
// @Tags         Health
// @Produce      json
// @Success      200  {object}  response.Response
// @Router       /health/live [get]
func (h *HealthHandler) Liveness(c *gin.Context) {
	response.Success(c, gin.H{
		"status": "alive",
	})
}

// Readiness godoc
// @Summary      Readiness probe
// @Description  Check if the service is ready to serve traffic (for Kubernetes)
// @Tags         Health
// @Produce      json
// @Success      200  {object}  response.Response
// @Failure      503  {object}  response.Response
// @Router       /health/ready [get]
func (h *HealthHandler) Readiness(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	// Check if database is ready
	if err := h.db.PingContext(ctx); err != nil {
		response.ServiceUnavailable(c, gin.H{
			"status":  "not_ready",
			"reason":  "database",
			"message": "database not ready",
		})
		return
	}

	// Check if Redis is ready
	if err := h.redisClient.Ping(ctx).Err(); err != nil {
		response.ServiceUnavailable(c, gin.H{
			"status":  "not_ready",
			"reason":  "redis",
			"message": "redis not ready",
		})
		return
	}

	response.Success(c, gin.H{
		"status": "ready",
	})
}

func (h *HealthHandler) checkDatabase(ctx context.Context) CheckResult {
	if err := h.db.PingContext(ctx); err != nil {
		return CheckResult{
			Status:  "unhealthy",
			Message: "database connection failed: " + err.Error(),
		}
	}

	return CheckResult{
		Status:  "healthy",
		Message: "connected",
	}
}

func (h *HealthHandler) checkRedis(ctx context.Context) CheckResult {
	if err := h.redisClient.Ping(ctx).Err(); err != nil {
		return CheckResult{
			Status:  "unhealthy",
			Message: "redis connection failed: " + err.Error(),
		}
	}

	return CheckResult{
		Status:  "healthy",
		Message: "connected",
	}
}
