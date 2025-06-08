package conf

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	jsoniter "github.com/json-iterator/go"
	"log"
	"os"
	"time"
)

// FiberConfig adalah konfigurasi untuk Fiber framework
// jsoniter make JSON encoding/decoding lebih efisien
var FiberConfig = fiber.Config{
	AppName:                   "grab-orders-service",
	Prefork:                   GetEnv("ENABLE_PREFORK", "false") == "true", // untuk multi-core, aktifkan di production
	Immutable:                 true,
	StrictRouting:             true,
	CaseSensitive:             true,
	BodyLimit:                 4 * 1024 * 1024, // 4MB
	ETag:                      false,
	ReadTimeout:               10 * time.Second,
	WriteTimeout:              10 * time.Second,
	IdleTimeout:               60 * time.Second,
	ProxyHeader:               "X-Forwarded-For",
	EnableTrustedProxyCheck:   true,
	TrustedProxies:            []string{"10.0.0.0/8", "192.168.0.0/16"},
	ErrorHandler:              CustomErrorHandler,
	DisableStartupMessage:     IsProduction(), // Nonaktifkan pesan startup di production
	DisableDefaultContentType: true,
	DisableDefaultDate:        true,
	ReduceMemoryUsage:         true,
	JSONEncoder:               jsoniter.ConfigCompatibleWithStandardLibrary.Marshal,
	JSONDecoder:               jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal,
	RequestMethods:            []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH", "HEAD"},
}

// CustomErrorHandler adalah handler untuk menangani error di Fiber
// Ini akan menangani error yang tidak tertangani dan mengembalikan response JSON
func CustomErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	message := "Internal Server Error"

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
		if !IsProduction() {
			message = e.Message // Show detailed error hanya di development
		}
	}

	// Structured logging
	log.Printf("[ERROR] %s | %s %s -> %d | %v | IP: %s | UA: %s",
		time.Now().Format(time.RFC3339),
		c.Method(),
		c.OriginalURL(),
		code,
		err,
		c.IP(),
		c.Get("User-Agent"),
	)

	return c.Status(code).JSON(fiber.Map{
		"success":   false,
		"message":   message,
		"timestamp": time.Now().Unix(),
	})
}

// CORSConfig adalah konfigurasi untuk CORS
var CORSConfig = cors.New(cors.Config{
	Next:             nil,
	AllowOriginsFunc: nil,
	AllowOrigins: func() string {
		if IsProduction() {
			return GetEnv("ALLOWED_ORIGINS", "https://yourdomain.com")
		}
		return "http://localhost:3000" // Allow untuk development
	}(),
	AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS,PATCH,HEAD",
	AllowHeaders:     "Origin,Content-Type,Accept,Authorization,X-API-Key",
	AllowCredentials: true,
	MaxAge:           3600, // 1 jam, ini mengatur berapa lama preflight request akan disimpan di cache respons browser (OPTIONS request)
})

// APIKeyMiddleware adalah middleware untuk memeriksa API Key
func APIKeyMiddleware() fiber.Handler {
	apiKey := GetEnv("API_KEY", "")
	if apiKey == "" {
		log.Fatal("API_KEY environment variable is required")
	}

	return func(c *fiber.Ctx) error {
		key := c.Get("X-API-Key")
		if key == "" || key != apiKey {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "Unauthorized - Invalid API Key",
			})
		}
		return c.Next()
	}
}

// RecoverMiddleware adalah middleware untuk menangani panic
var RecoverMiddleware = recover.New(recover.Config{
	EnableStackTrace: false,
	StackTraceHandler: func(c *fiber.Ctx, e interface{}) {
		// Log panic ke monitoring system
		log.Printf("Panic: %v", e)
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Internal Server Error",
		})
	},
})

// SecurityConfig adalah konfigurasi untuk keamanan aplikasi
// Menggunakan helmet untuk mengatur header keamanan
var SecurityConfig = helmet.New(helmet.Config{
	XSSProtection:             "1; mode=block",
	ContentTypeNosniff:        "nosniff",
	XFrameOptions:             "DENY",
	ReferrerPolicy:            "no-referrer",
	CrossOriginEmbedderPolicy: "require-corp",
	CrossOriginOpenerPolicy:   "same-origin",
	CrossOriginResourcePolicy: "cross-origin",
	OriginAgentCluster:        "?1",
	XDNSPrefetchControl:       "off",
	XDownloadOptions:          "noopen",
	XPermittedCrossDomain:     "none",
})

// RateLimitConfig adalah konfigurasi untuk rate limiting request
var RateLimitConfig = limiter.New(limiter.Config{
	Max:        120, // Maksimum 120 request per IP dalam periode waktu yang ditentukan
	Expiration: 1 * time.Minute,
	KeyGenerator: func(c *fiber.Ctx) string {
		return c.IP()
	},
	LimitReached: func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
			"success": false,
			"message": "Rate limit exceeded",
		})
	},
	SkipFailedRequests:     true,
	SkipSuccessfulRequests: false,
})

// CompressionConfig untuk compression
// Server → Internet → Client
// |         |         |
// 100KB     30KB      100KB
// JSON   (compressed)  JSON
var CompressionConfig = compress.New(compress.Config{
	Level: compress.LevelBestSpeed,
})

// LoggerConfig untuk logging
var LoggerConfig = logger.New(logger.Config{
	Format:     "${time} | ${status} | ${latency} | ${ip} | ${method} | ${path} | ${error}\n",
	TimeFormat: "2006-01-02 15:04:05",
	TimeZone:   "Local",
	Output:     os.Stdout,
	Done: func(c *fiber.Ctx, logString []byte) {
		// Custom log processing jika diperlukan
		// Bisa dikirim ke logging service seperti ELK, Fluentd, etc.
	},
})
