package handler

import (
	"github.com/gofiber/fiber/v2"
	"somapay-backend/ent"
	"somapay-backend/ent/product"
	"somapay-backend/ent/transaction"
	"strconv"
)

func CreateTransactionHandler(client *ent.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		u := c.Locals("user").(*ent.User)

		var req struct {
			ProductID int    `json:"product_id"`
			Quantity  int    `json:"quantity"`
			PIN       string `json:"pin"`
		}

		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
		}

		if u.Pin != req.PIN {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "invalid pin"})
		}

		tx, err := client.Tx(c.Context())
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "tx begin failed"})
		}

		defer func() {
			if err != nil {
				_ = tx.Rollback()
			}
		}()

		p, err := tx.Product.
			Query().
			Where(product.IDEQ(req.ProductID)).
			WithBooth().
			Only(c.Context())
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "product not found"})
		}

		total := p.Price * int64(req.Quantity)

		if u.Point < total {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "not enough balance"})
		}

		// User 포인트 차감
		_, err = tx.User.UpdateOne(u).
			AddPoint(-total).
			Save(c.Context())
		if err != nil {
			return err
		}

		boothID, err := tx.Product.
			Query().
			Where(product.IDEQ(p.ID)).
			QueryBooth().
			OnlyID(c.Context())
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot get booth id"})
		}

		t, err := tx.Transaction.
			Create().
			SetUserID(u.ID).
			SetProductID(p.ID).
			SetBoothID(boothID).
			SetQuantity(int64(req.Quantity)).
			SetAmount(total).
			SetStatus("SUCCESS").
			Save(c.Context())
		if err != nil {
			return err
		}

		if err = tx.Commit(); err != nil {
			return err
		}

		return c.JSON(t)
	}
}

func GetTransactionHandler(client *ent.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := c.Locals("user").(*ent.User)
		txID, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid transaction id"})
		}

		t, err := client.Transaction.
			Query().
			Where(transaction.IDEQ(txID)).
			WithUser().
			WithBooth(func(q *ent.BoothQuery) {
				q.WithUser()
			}).
			Only(c.Context())
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "transaction not found"})
		}

		if user.Role != "ADMIN" &&
			user.ID != t.Edges.User.ID &&
			!(user.Role == "HOST" && user.ID == t.Edges.Booth.Edges.User.ID) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "forbidden"})
		}

		return c.JSON(t)
	}
}

func ListTransactionsHandler(client *ent.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := c.Locals("user").(*ent.User)

		query := client.Transaction.Query().WithUser().WithBooth(func(q *ent.BoothQuery) { q.WithUser() })

		var ts []*ent.Transaction

		switch user.Role {
		case "ADMIN":
			var err error
			ts, err = query.All(c.Context())
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "query failed"})
			}

		case "HOST":
			allTx, err := query.All(c.Context())
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "query failed"})
			}

			for _, t := range allTx {
				if (t.Edges.Booth != nil && t.Edges.Booth.Edges.User.ID == user.ID) || t.Edges.User.ID == user.ID {
					ts = append(ts, t)
				}
			}

		default:
			allTx, err := query.All(c.Context())
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "query failed"})
			}

			for _, t := range allTx {
				if t.Edges.User.ID == user.ID {
					ts = append(ts, t)
				}
			}
		}

		return c.JSON(ts)
	}
}
