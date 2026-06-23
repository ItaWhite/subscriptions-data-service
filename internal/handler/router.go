package handler

import (
	"net/http"
	_ "subscriptions-data-service/docs"

	httpSwagger "github.com/swaggo/http-swagger"
)

func Router(h *recordHandler) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /records", h.GetRecordsHandler)
	mux.HandleFunc("GET /records/{id}", h.GetRecordHandler)
	mux.HandleFunc("POST /records", h.PostRecordHandler)
	mux.HandleFunc("PUT /records/{id}", h.PutRecordHandler)
	mux.HandleFunc("DELETE /records/{id}", h.DeleteRecordHandler)
	mux.HandleFunc("GET /records/total", h.GetTotalPriceHandler)

	mux.Handle("/swagger/", httpSwagger.WrapHandler)
	return mux
}
