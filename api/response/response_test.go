package response

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/valyala/fasthttp"
)

func TestGetErrorHTTPResponseBody(t *testing.T) {
	response := GetErrorHTTPResponseBody(800, "Custom Error Code")
	utils.AssertEqual(t, response.Code, 800)
	utils.AssertEqual(t, response.Message, "Custom Error Code")
}

func TestWriteHTTPResponseSuccess(t *testing.T) {
	app := fiber.New()
	ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
	defer app.ReleaseCtx(ctx)

	body := HTTPResponse{
		Code:    200,
		Message: "success",
		Content: map[string]interface{}{},
	}
	err := WriteHTTPResponse(ctx, 200, &body)
	if err != nil {
		t.Error("Write http response should not return error")
	}

	var response HTTPResponse
	_ = json.Unmarshal(ctx.Response().Body(), &response)

	utils.AssertEqual(t, response.Code, 200)
	utils.AssertEqual(t, response.Message, "success")

}

func TestWriteHTTPResponseError(t *testing.T) {
	app := fiber.New()
	ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
	defer app.ReleaseCtx(ctx)

	body := HTTPResponse{
		Code:    700,
		Message: "Custom Error Code",
		Content: map[string]interface{}{},
	}
	err := WriteHTTPResponse(ctx, 700, &body)
	if err == nil {
		t.Error("Write http response should return error")
	}

	body = HTTPResponse{
		Code:    70,
		Message: "Custom Error Code",
		Content: map[string]interface{}{},
	}
	err = WriteHTTPResponse(ctx, 70, &body)
	if err == nil {
		t.Error("Write http response should return error")
	}

}

func TestDefaultErrorHandler(t *testing.T) {
	app := fiber.New()
	ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
	defer app.ReleaseCtx(ctx)

	err := DefaultErrorHandler(ctx, errors.New(""))
	if err != nil {
		t.Error("Write http response should not return error")
	}

	var response HTTPResponse
	_ = json.Unmarshal(ctx.Response().Body(), &response)

	utils.AssertEqual(t, response.Code, 500)
	utils.AssertEqual(t, response.Message, "Internal Server Error")

}
