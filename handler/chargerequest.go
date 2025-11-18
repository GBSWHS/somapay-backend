package handler

import (
	"github.com/gofiber/fiber/v2"
	"somapay-backend/ent"
	"somapay-backend/ent/chargerequest"
	"strconv"
)

func CreateChargeRequestHandler(client *ent.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if !isUser(c) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "only users can create charge requests"})
		}

		var req struct {
			Amount int64 `json:"amount"`
		}
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
		}

		u := c.Locals("user").(*ent.User)

		cr, err := client.ChargeRequest.
			Create().
			SetAmount(req.Amount).
			SetUser(u). // 이게 이제 가능함
			Save(c.Context())

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to create request"})
		}

		return c.JSON(cr)
	}
}

func UpdateChargeRequestHandler(client *ent.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if !isAdmin(c) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "only admin can approve/reject"})
		}

		chargeID, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
		}

		var req struct {
			Status string `json:"status"` // APPROVED / REJECTED
		}
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
		}

		if req.Status != "APPROVED" && req.Status != "REJECTED" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid status"})
		}

		cr, err := client.ChargeRequest.
			Query().
			Where(chargerequest.IDEQ(chargeID)).
			WithUser().
			Only(c.Context())

		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "request not found"})
		}

		// 승인 시 유저 포인트 증가
		if req.Status == "APPROVED" {
			_, err := client.User.
				UpdateOne(cr.Edges.User).
				AddPoint(cr.Amount).
				Save(c.Context())
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to update user points"})
			}
		}

		updated, err := cr.Update().
			SetStatus(req.Status).
			Save(c.Context())

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to update request"})
		}

		return c.JSON(updated)
	}
}

func GetChargeRequestHandler(client *ent.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		chargeID, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
		}

		cr, err := client.ChargeRequest.
			Query().
			Where(chargerequest.IDEQ(chargeID)).
			WithUser().
			Only(c.Context())

		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "not found"})
		}

		u := c.Locals("user").(*ent.User)

		if !isAdmin(c) && u.ID != cr.Edges.User.ID {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "forbidden"})
		}

		return c.JSON(cr)
	}
}

func ListChargeRequestsHandler(client *ent.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if !isAdmin(c) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "only admin allowed"})
		}

		crs, err := client.ChargeRequest.
			Query().
			WithUser().
			All(c.Context())

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "query failed"})
		}

		return c.JSON(crs)
	}
}
