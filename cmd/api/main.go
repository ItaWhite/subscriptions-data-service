package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"subscriptions-data-service/internal"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	err := godotenv.Load()
	if err != nil {
		slog.Warn("error loading .env file", "error", err)
	}

	db, err := internal.ConnectDb(ctx, os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_HOST"), os.Getenv("POSTGRES_DB"))
	if err != nil {
		log.Fatal("error connecting to database: ", err)
	}
	defer db.Close()

	r := internal.NewRecordRepository(db)
	s := internal.NewRecordService(r)
	h := internal.NewRecordHandler(s)
	mux := internal.Router(h)

	server := http.Server{
		Addr:    fmt.Sprintf(":%s", os.Getenv("SERVER_PORT")),
		Handler: mux,
	}

	slog.Info("server started", "port", os.Getenv("SERVER_PORT"))
	go func() {
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("server failed", "error", err)
			cancel()
		}
	}()

	<-ctx.Done()

	ctx, stop := context.WithTimeout(context.Background(), 5*time.Second)
	defer stop()
	err = server.Shutdown(ctx)
	if err != nil {
		slog.Error("shutdown error", "error", err)
	} else {
		slog.Info("server stopped")
	}
}
