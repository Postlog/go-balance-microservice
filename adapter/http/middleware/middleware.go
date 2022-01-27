package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/postlog/go-balance-microservice/adapter/http/middleware/context"
	"github.com/postlog/go-balance-microservice/adapter/http/middleware/errors"
	"github.com/postlog/go-balance-microservice/adapter/http/middleware/log"
	"github.com/postlog/go-balance-microservice/pkg/logger"
)

func Register(router fiber.Router, logger logger.Logger) {
	router.Use(
		context.BuildHandler(),
		errors.BuildHandler(logger),
		log.BuildHandler(logger),
	)
}
