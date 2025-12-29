package main

import (
	"context"
	"log"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"github.com/szuryanailham/expense-tracker/internal/env"
)

func main() {
	err := godotenv.Load()
	 if err != nil {
    	log.Fatal("Error loading .env file")
  	}
	ctx := context.Background()
	cfg := config{
		addr:":8080",
		db: dbConfig{
			dsn:env.GetString("GOOSE_DBSTRING","host=127.0.0.1 port=5432 user=root password=root dbname=expense_tracker sslmode=disable"),
		},
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout,nil))
	slog.SetDefault(logger)

	// connect to database
	conn, err := pgx.Connect(ctx, cfg.db.dsn)
	if err != nil {
		panic(err)
	}
	defer conn.Close(ctx)
	logger.Info("database connected to database", "dsn", cfg.addr)
	
	api:= application{
		config: cfg,
		db:conn,
	}

	if err := api.run(api.mount()); err != nil {
		slog.Error("server failed to start","errors", err)
		os.Exit(1)
	}

}