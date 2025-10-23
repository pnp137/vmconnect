package health

import (
	"os"

	"github.com/gofiber/fiber/v2"
)

func Health(c *fiber.Ctx) error {
	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	services := map[string]interface{}{}
	response := map[string]interface{}{
		"service":   "vmconnect-api",
		"component": "api",
		"hostname":  hostname,
		"services":  services,
	}
	return c.JSON(response)
}
