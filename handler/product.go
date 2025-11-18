package handler

import (
	"github.com/gofiber/fiber/v2"
	"somapay-backend/ent"
	"somapay-backend/ent/product"
	"strconv"
)

func GetProductHandler(client *ent.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		productID, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
		}

		p, err := client.Product.Get(c.Context(), productID)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "not found"})
		}

		return c.JSON(p)
	}
}

func CreateProductHandler(client *ent.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req struct {
			BoothID     int    `json:"booth_id"`
			Name        string `json:"name"`
			Description string `json:"description"`
			Price       int    `json:"price"`
			Quantity    int    `json:"quantity"`
		}
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
		}

		if !canManageProduct(c, req.BoothID, client) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "forbidden"})
		}

		p, err := client.Product.
			Create().
			SetName(req.Name).
			SetDescription(req.Description).
			SetPrice(int64(req.Price)).
			SetQuantity(int64(req.Quantity)).
			Save(c.Context())
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "create failed"})
		}

		return c.JSON(p)
	}
}

func UpdateProductHandler(client *ent.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		productID, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
		}

		boothID, err := client.Product.
			Query().
			Where(product.IDEQ(productID)).
			QueryBooth().
			OnlyID(c.Context())
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "booth not found"})
		}

		if !canManageProduct(c, boothID, client) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "forbidden"})
		}

		var req struct {
			Name        *string `json:"name"`
			Description *string `json:"description"`
			Price       *int    `json:"price"`
			Quantity    *int    `json:"quantity"`
		}
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
		}

		q := client.Product.UpdateOneID(productID)
		if req.Name != nil {
			q.SetName(*req.Name)
		}
		if req.Description != nil {
			q.SetDescription(*req.Description)
		}
		if req.Price != nil {
			q.SetPrice(int64(*req.Price))
		}
		if req.Quantity != nil {
			q.SetQuantity(int64(*req.Quantity))
		}

		updated, err := q.Save(c.Context())
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "update failed"})
		}

		return c.JSON(updated)
	}
}

func DeleteProductHandler(client *ent.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		productID, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
		}

		boothID, err := client.Product.
			Query().
			Where(product.IDEQ(productID)).
			QueryBooth().
			OnlyID(c.Context())
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "booth not found"})
		}

		if !canManageProduct(c, boothID, client) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "forbidden"})
		}

		if err := client.Product.DeleteOneID(productID).Exec(c.Context()); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "delete failed"})
		}

		return c.SendStatus(fiber.StatusNoContent)
	}
}
