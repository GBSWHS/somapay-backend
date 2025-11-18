package handler

import (
	"github.com/gofiber/fiber/v2"
	"somapay-backend/ent"
	"strconv"
)

func CreateUserHandler(client *ent.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if !isAdmin(c) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "forbidden"})
		}

		var req struct {
			Username string `json:"username"`
			Password string `json:"password"`
			Pin      string `json:"pin"`
			Role     string `json:"role"`
		}

		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
		}

		u, err := client.User.
			Create().
			SetUsername(req.Username).
			SetPassword(req.Password).
			SetPin(req.Pin).
			SetRole(req.Role).
			SetPoint(0).
			Save(c.Context())

		if err != nil {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "duplicated"})
		}

		return c.JSON(u)
	}
}

func GetUserHandler(client *ent.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		targetID64, err := strconv.ParseInt(c.Params("id"), 10, 64)
		targetID := int(targetID64)

		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
		}

		if !isAdmin(c) && !isSelf(c, targetID) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "forbidden"})
		}

		u, err := client.User.Get(c.Context(), targetID)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "not found"})
		}

		return c.JSON(u)
	}
}

func UpdateUserHandler(client *ent.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		targetID64, err := strconv.ParseInt(c.Params("id"), 10, 64)
		targetID := int(targetID64)

		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
		}

		if !isAdmin(c) && !isSelf(c, targetID) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "forbidden"})
		}

		var req struct {
			Password *string `json:"password"`
			Pin      *string `json:"pin"`
		}

		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
		}

		q := client.User.UpdateOneID(targetID)

		if req.Password != nil {
			q.SetPassword(*req.Password)
		}
		if req.Pin != nil {
			q.SetPin(*req.Pin)
		}

		u, err := q.Save(c.Context())
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "update failed"})
		}

		return c.JSON(u)
	}
}
