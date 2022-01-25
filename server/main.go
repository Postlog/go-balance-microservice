package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/postlog/go-balance-microservice/internal/config"
	"github.com/postlog/go-balance-microservice/internal/database"
	"github.com/postlog/go-balance-microservice/internal/logger"
	"github.com/postlog/go-balance-microservice/internal/transaction"
)

func main() {
	configPath := flag.String("config", "config/dev.json", "path to the config file")
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		panic(fmt.Errorf("unable to load config: %v", err))
	}

	l, err := logger.New(cfg)
	if err != nil {
		panic(fmt.Errorf("unable to configure l: %v", err))
	}
	defer l.Flush()

	db, err := database.New(cfg.DatabaseAddress)
	if err != nil {
		panic(fmt.Errorf("unable to establish connection to database: %v", err))
	}
	ctx := context.Background()
	r := transaction.NewRepository(db)
	s := transaction.NewService(r, l)
	ts, err := s.GetTransactions(ctx, "11111111-3a7a-4d5e-8a6c-febc8c5b3f13", -1, 20, "amount", "desc")
	if err != nil {
		panic(err)
	}

	for _, t := range ts {
		data, err := json.Marshal(&t)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(string(data))
	}
}
