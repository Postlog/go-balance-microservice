package errors

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/postlog/go-balance-microservice/adapter/http/protocol"
	"github.com/postlog/go-balance-microservice/pkg/errors"
	"github.com/postlog/go-balance-microservice/pkg/logger"
	"runtime/debug"
)

// BuildHandler returns handler, that handle error
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

			resp, code := buildErrorResponse(err)
			if code == fiber.StatusInternalServerError {
				l.Errorf("unexpected error (%s)\n%s", err, debug.Stack())
			}
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
		errorResp  protocol.Error
		statusCode int
	)

	switch typedErr := err.(type) {
	case errors.ArgumentError:
		errorResp.Message = typedErr.Error()
		statusCode = fiber.StatusBadRequest
		break
	case errors.ServiceError:
		errorResp.Message = typedErr.Error()
		code := typedErr.GetCode()
		errorResp.Code = &code
		statusCode = fiber.StatusBadRequest
		break
	case *fiber.Error:
		errorResp.Message = typedErr.Error()
		statusCode = typedErr.Code
		break
	default:
		errorResp.Message = "unexpected internal server error"
		statusCode = fiber.StatusInternalServerError
		break
	}

	return protocol.Response{Error: &errorResp, Payload: nil}, statusCode
}
