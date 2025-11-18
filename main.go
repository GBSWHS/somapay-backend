package main

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"log"
	"somapay-backend/ent"
	"somapay-backend/handler"
	"somapay-backend/middleware"
	"somapay-backend/storage"
)

func main() {
	app := fiber.New()
	client := storage.GetClient()
	sessionStore := storage.GetSessionStore()

	ctx := context.Background()

	setupCors(app)
	setupRoutes(app, client, sessionStore, ctx)

	if err := app.Listen("0.0.0.0:8443"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func setupCors(app *fiber.App) {
	app.Use(cors.New(cors.Config{
		AllowOrigins: "https://pay.gbsw.hs.kr",
		AllowMethods: "GET,POST,PUT,PATCH,DELETE,OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))
}

func setupRoutes(app *fiber.App, client *ent.Client, sessionStore *storage.SessionStore, ctx context.Context) {

	// Index (Ping)
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	// Login Route
	app.Post("/login", handler.LoginHandler(client, sessionStore))

	// Auth Middleware
	auth := middleware.AuthMiddleware(client, sessionStore)

	// User Routes
	userGroup := app.Group("/users", auth)
	userGroup.Post("/", handler.CreateUserHandler(client))
	userGroup.Get("/:id", handler.GetUserHandler(client))
	userGroup.Patch("/:id", handler.UpdateUserHandler(client))

	// Booth Routes
	boothGroup := app.Group("/booths", auth)
	boothGroup.Post("/", handler.CreateBoothHandler(client))
	boothGroup.Get("/", handler.ListBoothsHandler(client))
	boothGroup.Get("/:id", handler.GetBoothHandler(client))
	boothGroup.Patch("/:id", handler.UpdateBoothHandler(client))
	boothGroup.Delete("/:id", handler.DeleteBoothHandler(client))

	// Product Routes
	productGroup := app.Group("/products", auth)
	productGroup.Get("/", handler.ListAllProductsHandler(client))
	productGroup.Get("/booth/:booth_id", handler.ListProductsByBoothHandler(client))
	productGroup.Get("/:id", handler.GetProductHandler(client))
	productGroup.Post("/", handler.CreateProductHandler(client))
	productGroup.Patch("/:id", handler.UpdateProductHandler(client))
	productGroup.Delete("/:id", handler.DeleteProductHandler(client))

	// Charge Request Routes
	chargeGroup := app.Group("/charge-requests", auth)
	chargeGroup.Post("/", handler.CreateChargeRequestHandler(client))
	chargeGroup.Get("/", handler.ListChargeRequestsHandler(client))
	chargeGroup.Get("/:id", handler.GetChargeRequestHandler(client))
	chargeGroup.Patch("/:id", handler.UpdateChargeRequestHandler(client))

	// Transaction Routes
	transactionGroup := app.Group("/transactions", auth)
	transactionGroup.Post("/", handler.CreateTransactionHandler(client))
	transactionGroup.Get("/", handler.ListTransactionsHandler(client))
	transactionGroup.Get("/:id", handler.GetTransactionHandler(client))
}
