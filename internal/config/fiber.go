package config

import "github.com/gofiber/fiber/v2"

// GetFiberConfig returns fiber configuration based on app config
func GetFiberConfig(cfg *Config) fiber.Config {
	return fiber.Config{
		AppName: cfg.App.Name,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"error":   err.Error(),
				"success": false,
			})
		},
		// Enable prefork in production for better performance
		Prefork: cfg.App.Environment == "production",
	}
}
