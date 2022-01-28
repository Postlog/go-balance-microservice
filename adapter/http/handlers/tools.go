package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/postlog/go-balance-microservice/adapter/http/protocol"
)

// ParseBody parses JSON from the body of provided through context request
func ParseBody(c *fiber.Ctx, req protocol.ValidatableRequest) error {
	if err := c.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "unexpected JSON schema")
	}
	if err := req.Validate(); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return nil
}
