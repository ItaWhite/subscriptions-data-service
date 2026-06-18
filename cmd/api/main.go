package main

import (
	"context"
	"fmt"
	"log"
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

	err := godotenv.Load("cmd/api/.env")
	if err != nil {
		log.Fatal("error loading .env file", "error", err)
	}

	db, err := internal.ConnectDb(ctx, os.Getenv("POSTGRES_URL"))
	if err != nil {
		log.Fatal("error connecting to database")
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

	fmt.Println("server started")
	go func() {
		log.Fatal(server.ListenAndServe())
	}()

	<-ctx.Done()

	ctx, stop := context.WithTimeout(context.Background(), 5*time.Second)
	defer stop()
	log.Fatal(server.Shutdown(ctx))
}
