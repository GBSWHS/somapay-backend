package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"somapay-backend/ent"
	"somapay-backend/ent/user"
	"somapay-backend/storage"
)

func LoginHandler(client *ent.Client, sessionStore *storage.SessionStore) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req struct {
			Username string `json:"student_number"`
			Password string `json:"password"`
		}

		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
		}

		u, err := client.User.
			Query().
			Where(user.UsernameEQ(req.Username)).
			Only(c.Context())
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "user not found"})
		}

		if !checkPasswordHash(req.Password, u.Password) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid password"})
		}

		// UUID 생성
		token := uuid.New().String()

		// 메모리 스토어에 저장
		sessionStore.Lock()
		sessionStore.Data[token] = u.ID
		sessionStore.Unlock()

		return c.JSON(fiber.Map{"userId": u.ID, "token": token})
	}
}
