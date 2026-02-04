package middlewares

import (
	"fmt"
	"time"

	"go-boilerplate/apps/internal/database"

	"github.com/gofiber/fiber/v2"
)

func RateLimiter(limit int, window time.Duration) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.Context()

		key := fmt.Sprintf("rl:%s:%s", c.IP(), c.Path())

		count, err := database.Redis.Incr(ctx, key).Result()
		if err != nil {
			return fiber.ErrInternalServerError
		}

		if count == 1 {
			_ = database.Redis.Expire(ctx, key, window)
		}

		remaining := limit - int(count)
		remaining = max(limit-int(count), 0)

		c.Set("X-RateLimit-Limit", fmt.Sprintf("%d", limit))
		c.Set("X-RateLimit-Remaining", fmt.Sprintf("%d", remaining))

		if count > int64(limit) {
			ttl, _ := database.Redis.TTL(ctx, key).Result()
			c.Set("Retry-After", fmt.Sprintf("%d", int(ttl.Seconds())))

			return fiber.ErrTooManyRequests
		}

		return c.Next()
	}
}
