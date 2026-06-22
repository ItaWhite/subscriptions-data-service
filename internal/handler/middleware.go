package handler

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
)

const (
	loggerKey = "logger"
)

func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get("X-Request-ID")
		if requestID == "" {
			requestID = uuid.NewString()
		}
		r.Header.Set("X-Request-ID", requestID)
		w.Header().Set("X-Request-ID", requestID)

		next.ServeHTTP(w, r)
	})
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get("X-Request-ID")
		logger := slog.With("request_id", requestID, "method", r.Method, "url", r.URL.String())
		ctx := context.WithValue(r.Context(), loggerKey, logger)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func Panic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger := ctx.Value(loggerKey).(*slog.Logger)

		defer func() {
			p := recover()
			if p != nil {
				logger.Error("panic recovered", "error", p)
				http.Error(w, "внутренняя ошибка", http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
