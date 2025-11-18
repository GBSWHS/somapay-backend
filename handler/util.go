package handler

import (
	"github.com/gofiber/fiber/v2"
	"somapay-backend/ent"
	"somapay-backend/ent/booth"
)

func isUser(c *fiber.Ctx) bool {
	u := c.Locals("user").(*ent.User)
	return u.Role == "USER"
}

func isSelf(c *fiber.Ctx, targetUserID int) bool {
	u := c.Locals("user").(*ent.User)
	return u.ID == targetUserID
}

func isAdmin(c *fiber.Ctx) bool {
	u := c.Locals("user").(*ent.User)
	return u.Role == "ADMIN"
}

func isHostOfBooth(c *fiber.Ctx, boothID int, client *ent.Client) bool {
	u := c.Locals("user").(*ent.User)

	b, err := client.Booth.
		Query().
		Where(booth.IDEQ(boothID)).
		WithUser().
		Only(c.Context())

	if err != nil {
		return false
	}
	return u.Role == "HOST" && b.Edges.User.ID == u.ID
}

func canManageProduct(c *fiber.Ctx, boothID int, client *ent.Client) bool {
	return isAdmin(c) || isHostOfBooth(c, boothID, client)
}
