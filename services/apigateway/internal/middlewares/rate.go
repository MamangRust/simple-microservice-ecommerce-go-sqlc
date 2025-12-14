package middlewares

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/MamangRust/simple_microservice_ecommerce/apigateway/pkg/logger"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type RateLimiterMiddleware struct {
	redis  *redis.Client
	logger logger.LoggerInterface
	rps    int
	burst  int
	window time.Duration
}

func NewRateLimiterMiddleware(redis *redis.Client, logger logger.LoggerInterface, rps int, burst int) *RateLimiterMiddleware {
	return &RateLimiterMiddleware{
		redis:  redis,
		logger: logger,
		rps:    rps,
		burst:  burst,
		window: time.Second,
	}
}

func (rl *RateLimiterMiddleware) Limit(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ip := c.RealIP()

		allowed, remaining, resetTime, err := rl.allowRequest(c.Request().Context(), ip)
		if err != nil {
			rl.logger.Error("Rate limiter error", zap.Error(err), zap.String("ip", ip))
			return next(c)
		}

		c.Response().Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", rl.rps))
		c.Response().Header().Set("X-RateLimit-Remaining", fmt.Sprintf("%d", remaining))
		c.Response().Header().Set("X-RateLimit-Reset", fmt.Sprintf("%d", resetTime.Unix()))

		if !allowed {
			c.Response().Header().Set("Retry-After", fmt.Sprintf("%d", int(time.Until(resetTime).Seconds())))
			return c.JSON(http.StatusTooManyRequests, map[string]interface{}{
				"error":      "Too many requests, please try again later",
				"retryAfter": int(time.Until(resetTime).Seconds()),
			})
		}

		return next(c)
	}
}

func (rl *RateLimiterMiddleware) allowRequest(ctx context.Context, key string) (bool, int, time.Time, error) {
	now := time.Now()

	redisKey := fmt.Sprintf("ratelimit:%s", key)

	script := redis.NewScript(`
		local key = KEYS[1]
		local now = tonumber(ARGV[1])
		local window = tonumber(ARGV[2])
		local limit = tonumber(ARGV[3])
		local windowStart = now - window

		-- Remove old entries outside the window
		redis.call('ZREMRANGEBYSCORE', key, 0, windowStart)

		-- Count current requests in window
		local current = redis.call('ZCARD', key)

		if current < limit then
			-- Add current request
			redis.call('ZADD', key, now, now)
			redis.call('EXPIRE', key, window)
			return {1, limit - current - 1}
		else
			return {0, 0}
		end
	`)

	result, err := script.Run(ctx, rl.redis,
		[]string{redisKey},
		now.UnixNano(),
		rl.window.Nanoseconds(),
		rl.rps,
	).Result()

	if err != nil {
		return false, 0, now, err
	}

	resultSlice := result.([]interface{})
	allowed := resultSlice[0].(int64) == 1
	remaining := int(resultSlice[1].(int64))

	resetTime := now.Add(rl.window)

	return allowed, remaining, resetTime, nil
}
