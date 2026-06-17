package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"subscriptions-data-service/internal"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("cmd/api/.env")
	if err != nil {
		log.Fatal("error loading .env file", "error", err)
	}

	db, err := internal.ConnectDb(os.Getenv("POSTGRES_URL"))
	if err != nil {
		log.Fatal("error connecting to database")
	}
	defer db.Close()
	r := internal.NewRecordRepository(db)

	err = r.DropSchema()
	if err != nil {
		log.Fatal(err)
	}
	err = r.InitSchema()
	if err != nil {
		log.Fatal(err)
	}

	s := internal.NewRecordService(r)
	h := internal.NewRecordHandler(s)

	mux := internal.Router(h)

	server := http.Server{
		Addr:    fmt.Sprintf(":%s", os.Getenv("SERVER_PORT")),
		Handler: mux,
	}

	fmt.Println("server started")
	log.Fatal(server.ListenAndServe())
}
