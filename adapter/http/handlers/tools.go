package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/postlog/go-balance-microservice/adapter/http/protocol"
)

// ParseBody parses JSON from the body of provided through context request
func ParseBody(c *fiber.Ctx, req protocol.ValidatableRequest) error {
	if err := c.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, buildJSONErrorMessage(err))
	}
	if err := req.Validate(); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return nil
}

// buildJSONErrorMessage converts the JSON error to the suitable message
func buildJSONErrorMessage(err error) string {
	msg := "unexpected JSON schema"

	switch err.(type) {
	case *json.UnmarshalTypeError:
		typedError := err.(*json.UnmarshalTypeError)
		field := typedError.Field
		provided := typedError.Value
		expected := typedError.Type.String()

		msg += " " + fmt.Sprintf("%s expected, but %s provided in field %s", expected, provided, field)
		break
	}

	return msg
}
