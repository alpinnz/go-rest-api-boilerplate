package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
	ginlimiter "github.com/ulule/limiter/v3/drivers/middleware/gin"
	memory "github.com/ulule/limiter/v3/drivers/store/memory"
)

func RateLimiterMiddleware() gin.HandlerFunc {
	// Example: 10 requests per minute
	rate, err := limiter.NewRateFromFormatted("10-M")
	if err != nil {
		panic(err)
	}

	store := memory.NewStore()
	instance := limiter.New(store, rate)

	// Wrap as Gin middleware
	return ginlimiter.NewMiddleware(instance)
}
