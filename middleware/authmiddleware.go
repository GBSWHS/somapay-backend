package middleware

import (
	"github.com/gofiber/fiber/v2"
	"somapay-backend/ent"
	"somapay-backend/storage"
)

func AuthMiddleware(client *ent.Client, sessionStore *storage.SessionStore) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Get("Authorization")
		if token == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "missing token"})
		}

		sessionStore.RLock()
		userID, ok := sessionStore.Data[token]
		sessionStore.RUnlock()

		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token"})
		}

		u, err := client.User.Get(c.Context(), userID)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid user"})
		}

		c.Locals("user", u)
		return c.Next()
	}
}
