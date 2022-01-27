package context

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/postlog/go-balance-microservice/pkg/utils"
)

// BuildHandler returns handler, that sets the user context to the fiber context and adds the request id to it
func BuildHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := context.Background()
		c.SetUserContext(utils.WithRequestID(ctx))
		return c.Next()
	}
}
