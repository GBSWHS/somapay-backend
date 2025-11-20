package handler

import (
	"github.com/gofiber/fiber/v2"
	"log"
	"somapay-backend/ent"
	"somapay-backend/ent/booth"
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
		}

		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
		}

		if !canManageProduct(c, req.BoothID, client) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "forbidden"})
		}

		p, err := client.Product.
			Create().
			SetBoothID(req.BoothID).
			SetName(req.Name).
			SetDescription(req.Description).
			SetPrice(int64(req.Price)).
			Save(c.Context())

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "create failed"})
		}

		return c.JSON(p)
	}
}

func ListAllProductsHandler(client *ent.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if !isAdmin(c) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "forbidden"})
		}

		ps, err := client.Product.Query().All(c.Context())
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "query failed"})
		}

		return c.JSON(ps)
	}
}

func ListProductsByBoothHandler(client *ent.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		boothID, err := strconv.Atoi(c.Params("booth_id"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid booth id"})
		}

		ps, err := client.Product.
			Query().
			Where(product.HasBoothWith(booth.IDEQ(boothID))).
			All(c.Context())

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "query failed"})
		}

		return c.JSON(ps)
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
		}

		if err := c.BodyParser(&req); err != nil {
			log.Fatalf("failed to parse body: %v", err)
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
