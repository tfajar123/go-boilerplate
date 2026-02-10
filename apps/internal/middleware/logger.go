package middlewares

import (
	"encoding/json"
	"time"

	"go-boilerplate/apps/internal/utils"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func RequestLogger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		var requestBody []byte
		if c.Method() == "POST" || c.Method() == "PUT" || c.Method() == "PATCH" {
			body := c.Request().Body()
			if len(body) > 0 {
				requestBody = make([]byte, len(body))
				copy(requestBody, body)

				c.Request().SetBody(requestBody)
			}
		}

		err := c.Next()

		latency := time.Since(start)

		logFields := []zap.Field{
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
			zap.Int("status", c.Response().StatusCode()),
			zap.Duration("latency", latency),
			zap.String("ip", c.IP()),
		}

		// Tambahkan detail request jika di development
		if utils.IsDevelopment() {
			if len(requestBody) > 0 {
				var parsedBody any
				if err := json.Unmarshal(requestBody, &parsedBody); err == nil {
					logFields = append(logFields, zap.Any("request_body", parsedBody))
				} else {
					logFields = append(logFields, zap.String("request_body", string(requestBody)))
				}
			}

			queryParams := c.Request().URI().QueryArgs()
			if queryParams.Len() > 0 {
				queryMap := make(map[string]string)
				queryParams.VisitAll(func(key, value []byte) {
					queryMap[string(key)] = string(value)
				})
				logFields = append(logFields, zap.Any("query_params", queryMap))
			}

			// Tambahkan headers penting
			headers := []string{"Content-Type", "Authorization", "User-Agent"}
			headerMap := make(map[string]string)
			for _, header := range headers {
				if value := c.Get(header); value != "" {
					headerMap[header] = value
				}
			}
			if len(headerMap) > 0 {
				logFields = append(logFields, zap.Any("headers", headerMap))
			}
		}

		utils.Logger.Info("http request", logFields...)

		return err
	}
}
