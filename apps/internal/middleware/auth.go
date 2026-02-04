package middlewares

import (
	"strings"

	"go-boilerplate/apps/internal/database"
	"go-boilerplate/apps/internal/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

func AuthRequired(redis *redis.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return utils.Unauthorized(c, "authorization header tidak ada", nil)
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return utils.Unauthorized(c, "format token tidak valid", nil)
		}

		claims, err := utils.ParseAccessToken(parts[1])
		if err != nil {
			return utils.Unauthorized(c, "Token tidak valid", err.Error())
		}

		userID := claims["sub"].(string)
		sessionID := claims["sid"].(string)

		// ============================
		// CEK SESSION KE REDIS (WAJIB)
		// ============================
		key := "auth:session:" + userID
		storedSID, err := redis.Get(c.Context(), key).Result()
		if err != nil || storedSID != sessionID {
			return utils.Unauthorized(c, "session sudah logout", err.Error())
		}

		c.Locals("user_id", userID)
		c.Locals("email", claims["email"].(string))
		c.Locals("session_id", sessionID)

		return c.Next()
	}
}

func SessionAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		auth := c.Get("Authorization")

		if auth == "" {
			return c.Status(401).JSON(fiber.Map{
				"message": "authorization header missing",
			})
		}

		parts := strings.Split(auth, " ")
		if len(parts) != 2 {
			return c.Status(401).JSON(fiber.Map{
				"message": "invalid token",
			})
		}

		session, err := utils.GetSession(
			c.Context(),
			database.Redis,
			parts[1],
		)
		if err != nil {
			return c.Status(401).JSON(fiber.Map{
				"message": "session expired",
			})
		}

		c.Locals("user_id", session.UserID)
		c.Locals("email", session.Email)

		return c.Next()
	}
}
