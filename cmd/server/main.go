package main

import (
	"flag"
	"fmt"
	"github.com/gofiber/fiber/v2"
	balanceHandlers "github.com/postlog/go-balance-microservice/adapter/http/handlers/balance"
	transactionHandlers "github.com/postlog/go-balance-microservice/adapter/http/handlers/transaction"
	"github.com/postlog/go-balance-microservice/adapter/http/middleware"
	"github.com/postlog/go-balance-microservice/config"
	balanceRepository "github.com/postlog/go-balance-microservice/dataservice/balance"
	currencyRepository "github.com/postlog/go-balance-microservice/dataservice/currency"
	transactionRepository "github.com/postlog/go-balance-microservice/dataservice/transaction"
	"github.com/postlog/go-balance-microservice/pkg/database"
	"github.com/postlog/go-balance-microservice/pkg/logger"
	"github.com/postlog/go-balance-microservice/service/balance"
	"github.com/postlog/go-balance-microservice/service/currency"
	"github.com/postlog/go-balance-microservice/service/transaction"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	configPath := flag.String("config", "config/dev.json", "path to the config file")
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		panic(fmt.Errorf("unable to load config: %s", err))
	}

	l, err := logger.New(cfg.Logger)
	if err != nil {
		panic(fmt.Errorf("unable to configure logger: %s", err))
	}
	defer l.Flush()

	db, err := database.New(cfg.DSN, l)
	if err != nil {
		panic(fmt.Errorf("unable to create new database: %s", err))
	}
	defer db.Close()

	app := setupApplication(cfg, l, db)

	gracefulShutdown(app, l)
	if err = app.Listen(":" + cfg.Port); err != nil {
		l.Errorf("unexpected error during serving connections: %s", err)
		return
	}
}

func gracefulShutdown(app *fiber.App, logger logger.Logger)  {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		<-signals
		if err := app.Shutdown(); err != nil {
			logger.Errorf("unexpected error during shutting down the server: %s", err)
		}
	}()
}

func setupApplication(cfg *config.Config, logger logger.Logger, db database.Database) *fiber.App {
	app := fiber.New()
	root := app.Group("")
	balanceAPIRouter := root.Group("balance")
	transactionAPIRouter := root.Group("transactions")

	middleware.Register(root, logger)

	balanceService := balance.NewService(balanceRepository.NewRepository(db))
	transactionService := transaction.NewService(transactionRepository.NewRepository(db))
	currencyService := currency.NewService(currencyRepository.NewClient(cfg.ExchangeRatesAPIKey, cfg.ApiRequestTimeout))

	balanceHandlers.Register(
		balanceAPIRouter,
		balanceService,
		transactionService,
		currencyService,
		db.GetTransactionWrapper(),
		cfg.BaseCurrency,
	)
	transactionHandlers.Register(transactionAPIRouter, transactionService)

	return app
}
