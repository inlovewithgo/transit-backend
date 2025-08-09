package middlewares

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/inlovewithgo/transit-backend/main/utils"
	"github.com/inlovewithgo/transit-backend/pkg/logger"
)

type RateLimiter struct {
	redisClient *redis.Client
}

func NewRateLimiter() *RateLimiter {
	redisAddr := utils.GetENV("REDIS_ADDR", "localhost:6379")
	redisPassword := utils.GetENV("REDIS_PASSWORD", "")

	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       0,
	})

	ctx := context.Background()
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		logger.Log.Error("Failed to connect to Redis for rate limiting: %v", err)
		return &RateLimiter{redisClient: nil}
	}

	logger.Log.Info("Rate limiter connected to Redis successfully")
	return &RateLimiter{redisClient: rdb}
}

func (rl *RateLimiter) WaitlistRateLimit() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if rl.redisClient == nil {
			logger.Log.Warn("Rate limiting disabled - Redis not available")
			return c.Next()
		}

		clientIP := c.IP()
		if clientIP == "" {
			clientIP = "unknown"
		}

		key := fmt.Sprintf("ratelimit:waitlist:%s", clientIP)
		ctx := context.Background()

		current, err := rl.redisClient.Get(ctx, key).Int()
		if err != nil && err != redis.Nil {
			logger.Log.Error("Failed to get rate limit count: %v", err)
			return c.Next()
		}

		limit := 5
		if current >= limit {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error":       true,
				"message":     "Rate limit exceeded. Please try again later.",
				"retry_after": 60,
			})
		}

		pipe := rl.redisClient.Pipeline()
		pipe.Incr(ctx, key)
		pipe.Expire(ctx, key, time.Minute)
		_, err = pipe.Exec(ctx)

		if err != nil {
			logger.Log.Error("Failed to update rate limit count: %v", err)
		}

		remaining := limit - (current + 1)
		if remaining < 0 {
			remaining = 0
		}

		c.Set("X-RateLimit-Limit", fmt.Sprintf("%d", limit))
		c.Set("X-RateLimit-Remaining", fmt.Sprintf("%d", remaining))
		c.Set("X-RateLimit-Reset", fmt.Sprintf("%d", time.Now().Add(time.Minute).Unix()))

		return c.Next()
	}
}

func (rl *RateLimiter) Close() error {
	if rl.redisClient != nil {
		return rl.redisClient.Close()
	}
	return nil
}
