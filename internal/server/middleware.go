package server

import (
	"context"
	"net/http"
	"strings"
	"time"

	log "github.com/ghkadim/highload_architect/internal/logger"
	"github.com/ghkadim/highload_architect/internal/models"
)

func authorize(inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := strings.Trim(r.Header.Get("Authorization"), " ")
		if len(auth) == 0 || !strings.HasPrefix(auth, "Bearer ") {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(auth, "Bearer ")
		r = r.WithContext(context.WithValue(r.Context(), models.BearerTokenCtxKey, token))
		inner.ServeHTTP(w, r)
	})
}

type responseWriter struct {
	http.ResponseWriter
	http.Hijacker
	code int
}

func (w *responseWriter) WriteHeader(statusCode int) {
	w.code = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func logger(inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		writer := &responseWriter{ResponseWriter: w, Hijacker: w.(http.Hijacker)}

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
