package middlewares

import (
	"fmt"
	"time"

	"go-boilerplate/apps/internal/database"

	"github.com/gofiber/fiber/v2"
)

func RateLimiter(limit int, window time.Duration) fiber.Handler {
	return func(c *fiber.Ctx) error {
		key := fmt.Sprintf("rl:%s:%s", c.IP(), c.Path())

		count, err := database.Redis.Incr(c.Context(), key).Result()
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"status":  500,
				"error":   err.Error(),
				"message": "Redis error",
			})
		}

		if count == 1 {
			database.Redis.Expire(c.Context(), key, window)
		}

		if count > int64(limit) {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"status":  429,
				"error":   "Too many requests",
				"message": "Terlalu banyak request, coba lagi nanti",
			})
		}

		return c.Next()
	}
}
