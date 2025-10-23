package health

import (
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
)

func TestHealth(t *testing.T) {
	app := fiber.New()

	app.Get("/health", func(ctx *fiber.Ctx) error {
		return Health(ctx)
	})

	req := httptest.NewRequest("GET", "/health", nil)
	resp, _ := app.Test(req)

	utils.AssertEqual(t, 200, resp.StatusCode)

	var apiResponse map[string]interface{}
	respBodyBytes, _ := io.ReadAll(resp.Body)

	_ = json.Unmarshal(respBodyBytes, &apiResponse)
	utils.AssertEqual(t, "vmconnect-api", apiResponse["service"])
	utils.AssertEqual(t, "api", apiResponse["component"])
	utils.AssertEqual(t, true, apiResponse["hostname"] != "")
	utils.AssertEqual(t, map[string]interface{}{}, apiResponse["services"])
}
