package handler

import (
	"github.com/gofiber/fiber/v2"
	"somapay-backend/ent"
	"strconv"
)

func CreateBoothHandler(client *ent.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if !isAdmin(c) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "forbidden"})
		}

		var req struct {
			Name   string `json:"name"`
			UserID int    `json:"user_id"`
		}

		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
		}

		b, err := client.Booth.
			Create().
			SetName(req.Name).
			SetUserID(req.UserID).
			Save(c.Context())
		if err != nil {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "failed to create"})
		}

		return c.JSON(b)
	}
}

func GetBoothHandler(client *ent.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		boothID64, err := strconv.Atoi(c.Params("id"))
		boothID := int(boothID64)

		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
		}

		if !isAdmin(c) && !isHostOfBooth(c, boothID, client) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "forbidden"})
		}

		b, err := client.Booth.Get(c.Context(), boothID)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "not found"})
		}

		return c.JSON(b)
	}
}

func UpdateBoothHandler(client *ent.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if !isAdmin(c) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "forbidden"})
		}

		boothID, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
		}

		var req struct {
			Name   *string `json:"name"`
			UserID *int    `json:"user_id"`
		}
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
		}

		q := client.Booth.UpdateOneID(boothID)
		if req.Name != nil {
			q.SetName(*req.Name)
		}
		if req.UserID != nil {
			q.SetUserID(*req.UserID)
		}

		b, err := q.Save(c.Context())
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "update failed"})
		}

		return c.JSON(b)
	}
}

func DeleteBoothHandler(client *ent.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if !isAdmin(c) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "forbidden"})
		}

		boothID, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
		}

		err = client.Booth.DeleteOneID(boothID).Exec(c.Context())
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "delete failed"})
		}

		return c.SendStatus(fiber.StatusNoContent)
	}
}
