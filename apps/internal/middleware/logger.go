package middlewares

import (
	"time"

	"go-boilerplate/apps/internal/utils"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func RequestLogger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		err := c.Next()

		latency := time.Since(start)

		utils.Logger.Info("http request",
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
			zap.Int("status", c.Response().StatusCode()),
			zap.Duration("latency", latency),
			zap.String("ip", c.IP()),
		)

		return err
	}
}
