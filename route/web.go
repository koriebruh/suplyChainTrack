package route

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/swagger"
	"github.com/koriebruh/suplyChainTrack/conf"
	"github.com/koriebruh/suplyChainTrack/internal/handler"
	"github.com/koriebruh/suplyChainTrack/internal/metirc"
	"github.com/koriebruh/suplyChainTrack/pkg"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func RunApplicationContext() {
	pkg.NewSlog()               // Initialize logging
	config := conf.LoadConfig() // Load configuration from .env or environment variables
	metricsExporter := metirc.NewAppMetricsExporter()

	/* APPLICATION SETTING */
	app := fiber.New()
	app.Use(metricsExporter.FiberMetricMiddleware()) // Middleware for collecting metrics
	app.Use(conf.RecoverMiddleware)                  //Handle panic recovery
	app.Use(conf.LoggerConfig)                       // Log requests
	app.Use(conf.SecurityConfig)                     // Security headers
	app.Use(conf.CompressionConfig)                  // Compress responses for more fast delivery
	app.Use(conf.CORSConfig)                         // Enable CORS for all routes
	app.Use(conf.RateLimitConfig)                    // Rate limiting to prevent abuse
	app.Use(etag.New())                              // ETag middleware for caching

	/*ROUTE FOR METRIC EXPORTER*/
	app.Get("/metrics", adaptor.HTTPHandler(promhttp.Handler())) // Endpoint to expose metrics

	/* ROUTES */
	api := app.Group("/api/v1")
	MetricRoute(api, config)
	api.Use(conf.APIKeyMiddleware())

	if err := app.Listen(fmt.Sprintf(":%v", config.AppConfig.Port)); err != nil {
		panic(err)
	}

	// Graceful shutdown
}

func MetricRoute(r fiber.Router, config *conf.Config) {
	metric := handler.NewMetricHandlerImpl(*config)
	r.Get("/docs/*", swagger.HandlerDefault)
	r.Get("/health", metric.Health)
}

func ProductsRoute(r fiber.Router, config *conf.Config) {}

func BlockchainTxRoute(r fiber.Router, config *conf.Config) {}

func SupplyChainRoute(r fiber.Router, config *conf.Config) {}

func StakeHolderRoute(r fiber.Router, config *conf.Config) {}
