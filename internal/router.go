package internal

import "net/http"

func Router(h *recordHandler) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /records", h.GetRecordsHandler)
	mux.HandleFunc("GET /records/{id}", h.GetRecordHandler)
	mux.HandleFunc("POST /records", h.PostRecordHandler)
	return mux
}
