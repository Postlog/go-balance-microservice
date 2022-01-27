package log

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/postlog/go-balance-microservice/pkg/logger"
	"time"
)

// BuildHandler returns handler, that logs processed requests
func BuildHandler(logger logger.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		l := logger.With(c.UserContext())

		err := c.Next()

		status := "successfully"
		if err != nil {
			status = fmt.Sprintf("with error: \"%s\"", err)
		}

		l.Infof("request [%s][%s] processed %s (duration: %dms)",
			c.Request().Header.Method(), c.Request().Header.RequestURI(), status,
			time.Now().Sub(start).Milliseconds(),
		)

		return err
	}
}
