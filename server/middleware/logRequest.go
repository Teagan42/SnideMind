package middleware

import (
	"net/http"

	"go.uber.org/zap"
)

func LogRequestMiddleware() func(http.Handler) http.Handler {
	log := zap.L().Named("RequestLogger")
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			body, err := PeekBody(log, r)
			if err != nil {
				log.Error("Error reading request body", zap.Error(err))
				http.Error(w, "Error reading request body", http.StatusInternalServerError)
				return
			}
			log.Info(
				"Request received",
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.String("query", r.URL.RawQuery),
				zap.String("body", string(body)),
			)

			next.ServeHTTP(w, r)
		})
	}
}
