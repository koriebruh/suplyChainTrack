package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/koriebruh/suplyChainTrack/conf"
)

func main() {
	config := conf.LoadConfig() // Load configuration from .env or environment variables
	app := fiber.New()

	app.Use(conf.RecoverMiddleware) //Handle panic recovery
	app.Use(conf.LoggerConfig)      // Log requests
	app.Use(conf.SecurityConfig)    // Security headers
	app.Use(conf.CompressionConfig) // Compress responses for more fast delivery
	app.Use(conf.CORSConfig)        // Enable CORS for all routes
	app.Use(conf.RateLimitConfig)   // Rate limiting to prevent abuse
	app.Use(etag.New())             // ETag middleware for caching

	api := app.Group("/api/v1")
	api.Use(conf.APIKeyMiddleware())

	api.Get("/hi", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"success": true, "message": "Welcome to the API!"})
	})

	// Tambahkan rute untuk /health dan /metrics jika perlu
	api.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"success": true,
			"status":  "healthy",
		})
	})

	if err := app.Listen(fmt.Sprintf(":%v", config.AppConfig.Port)); err != nil {
		panic(err)
	}
}
