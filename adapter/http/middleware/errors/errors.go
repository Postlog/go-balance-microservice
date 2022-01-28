package errors

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/postlog/go-balance-microservice/adapter/http/protocol"
	"github.com/postlog/go-balance-microservice/pkg/errors"
	"github.com/postlog/go-balance-microservice/pkg/logger"
	"runtime/debug"
)

// BuildHandler returns handler, that handle errors
//
// Handler newer returns an error
func BuildHandler(logger logger.Logger) fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		l := logger.With(c.UserContext())

		defer func() {
			if e := recover(); e != nil {
				var ok bool
				if err, ok = e.(error); !ok {
					err = fmt.Errorf("%v", e)
				}
				l.Errorf("recovered from panic (%s)\n%s", err, debug.Stack())
			}

			if err == nil {
				return
			}
			l.Errorf("unexpected error (%s)\n%s", err, debug.Stack())
			resp, code := buildErrorResponse(err)
			if err = c.Status(code).JSON(resp); err != nil {
				l.Errorf("failed to send error response: %s", err)
			}

		}()

		return c.Next()
	}
}

// buildErrorResponse builds protocol.Response object, that contains information about provided error
func buildErrorResponse(err error) (protocol.Response, int) {
	var (
		msg        string
		statusCode int
	)

	switch typedErr := err.(type) {
	case errors.ArgumentError:
		msg = typedErr.Error()
		statusCode = fiber.StatusBadRequest
		break
	case *fiber.Error:
		msg = typedErr.Message
		statusCode = typedErr.Code
		break
	default:
		msg = "unexpected internal server error"
		statusCode = fiber.StatusInternalServerError
		break
	}

	return protocol.Response{ErrorMessage: &msg, Payload: nil}, statusCode
}
