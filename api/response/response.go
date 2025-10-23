package response

import (
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"linksupply.io/vmconnect/logger"
)

type HTTPResponse struct {
	Code    int         `json:"code,omitempty"`
	Message string      `json:"message,omitempty"`
	Content interface{} `json:"content"`
}

type HTTPResponseContent struct {
	Count    int         `json:"count"`
	Previous *string     `json:"prev"`
	Next     *string     `json:"next"`
	Results  interface{} `json:"results"`
}

// Error details will be returned by a service function to the handler
type ErrorDetails struct {
	Code    int
	Message string
	Error   error
}

// StatusCode refer to Http Code
func WriteHTTPResponse(c *fiber.Ctx, statusCode int, responseBody *HTTPResponse) error {
	if statusCode < 100 || statusCode > 600 {
		return errors.New(fmt.Sprintf("Invalid status code for HTTP response: %v", statusCode))
	}
	c.Status(statusCode)
	err := c.JSON(responseBody)
	return err
}

// Code refer to Application Code
func GetErrorHTTPResponseBody(code int, message string, content ...map[string]interface{}) *HTTPResponse {

	response := &HTTPResponse{
		Code:    code,
		Message: message,
		Content: map[string]interface{}{}, // Initialize with an empty map
	}

	if len(content) > 0 && content[0] != nil {
		response.Content = content[0]
	}
	return response
}

func DefaultErrorHandler(c *fiber.Ctx, err error) error {
	log := logger.NewLogger()
	log.Error(err, "Error thrown from DefaultErrorHandler")
	errorBody := GetErrorHTTPResponseBody(500, "Internal Server Error")
	return WriteHTTPResponse(c, 500, errorBody)
}
