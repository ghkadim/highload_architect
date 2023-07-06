package server

import (
	"context"
	"net/http"
	"strings"
	"time"

	openapi "github.com/ghkadim/highload_architect/generated/go_server/go"
	log "github.com/ghkadim/highload_architect/internal/logger"
	"github.com/ghkadim/highload_architect/internal/models"
)

func authorize(inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := strings.Trim(r.Header.Get("Authorization"), " ")
		if len(auth) == 0 || !strings.HasPrefix(auth, "Bearer ") {
			err := openapi.EncodeJSONResponse(nil, func(i int) *int { return &i }(http.StatusUnauthorized), w)
			if err != nil {
				log.Error("Failed to encode Json response: %v", err)
			}
			return
		}

		token := strings.TrimPrefix(auth, "Bearer ")
		r = r.WithContext(context.WithValue(r.Context(), models.BearerTokenCtxKey, token))
		inner.ServeHTTP(w, r)
	})
}

type responseWriter struct {
	http.ResponseWriter
	code int
}

func (w *responseWriter) WriteHeader(statusCode int) {
	w.code = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func logger(inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		writer := &responseWriter{ResponseWriter: w}

		inner.ServeHTTP(writer, r)

		log.Debug(
			"%s %s %d %s",
			r.Method,
			r.RequestURI,
			writer.code,
			time.Since(start),
		)
	})
}
