package middleware

import (
	"context"
	"net/http"

	"github.com/teagan42/snidemind/internal/llm"
)

var LLMContextKey ContextKey = "llm"

func AddLLMMiddleware(llm llm.LLM) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctx = context.WithValue(ctx, LLMContextKey, llm)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}

func GetLLMFromContext(ctx context.Context) (*llm.LLM, bool) {
	if llm, ok := ctx.Value(LLMContextKey).(llm.LLM); ok {
		return &llm, true
	}
	return nil, false
}
