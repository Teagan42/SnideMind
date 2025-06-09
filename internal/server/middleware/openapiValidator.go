package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/getkin/kin-openapi/routers"
)

type ContextKey string

const (
	BodyKey        ContextKey = "body"
	QueryKey       ContextKey = "query"
	RouteKey       ContextKey = "route"
	RouteParamsKey ContextKey = "routeParams"
)

func OpenAPIValidationMiddleware(router routers.Router) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			route, routeParams, err := router.FindRoute(r)
			if err != nil {
				http.Error(w, "Route not found: "+err.Error(), http.StatusNotFound)
				return
			}

			input := &openapi3filter.RequestValidationInput{
				Request:     r,
				PathParams:  routeParams,
				QueryParams: r.URL.Query(),
				Route:       route,
				Options: &openapi3filter.Options{
					IncludeResponseStatus: true,
				},
			}
			ctx := r.Context()
			if r.Body != nil {
				var raw map[string]any
				if err := json.NewDecoder(r.Body).Decode(&raw); err != nil {
					http.Error(w, "Invalid JSON body: "+err.Error(), http.StatusBadRequest)
					return
				}
				ctx = context.WithValue(ctx, BodyKey, raw)
			}
			if r.URL.Query() != nil {
				rawQuery := make(map[string]any)
				for key, values := range r.URL.Query() {
					if len(values) > 0 {
						rawQuery[key] = values[0] // Use the first value for simplicity
					}
				}
				ctx = context.WithValue(ctx, QueryKey, rawQuery)
			}
			ctx = context.WithValue(ctx, RouteKey, route)
			ctx = context.WithValue(ctx, RouteParamsKey, routeParams)
			r = r.WithContext(ctx)

			if err := openapi3filter.ValidateRequest(r.Context(), input); err != nil {
				http.Error(w, "Request validation failed: "+err.Error(), http.StatusBadRequest)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func GetValidatedBody[T any](r *http.Request) (T, error) {
	var zero T
	val := r.Context().Value(BodyKey)
	if val == nil {
		return zero, fmt.Errorf("no validated body in context")
	}

	rawMap, ok := val.(map[string]any)
	if !ok {
		return zero, fmt.Errorf("unexpected type in context")
	}

	rawJSON, err := json.Marshal(rawMap)
	if err != nil {
		return zero, err
	}

	var result T
	if err := json.Unmarshal(rawJSON, &result); err != nil {
		return zero, err
	}

	return result, nil
}

func GetQueryParams[T any](r *http.Request) (T, error) {
	var zero T
	val := r.Context().Value(QueryKey)
	if val == nil {
		return zero, fmt.Errorf("no query params in context")
	}

	rawMap, ok := val.(map[string]any)
	if !ok {
		return zero, fmt.Errorf("unexpected query context type")
	}

	rawJSON, err := json.Marshal(rawMap)
	if err != nil {
		return zero, err
	}

	var result T
	if err := json.Unmarshal(rawJSON, &result); err != nil {
		return zero, err
	}

	return result, nil
}

func GetRoute[T any](r *http.Request) (T, error) {
	var zero T
	val := r.Context().Value(RouteKey)
	if val == nil {
		return zero, fmt.Errorf("no route in context")
	}

	route, ok := val.(T)
	if !ok {
		return zero, fmt.Errorf("unexpected route context type")
	}

	return route, nil
}

func GetRouteParams[T any](r *http.Request) (T, error) {
	var zero T
	val := r.Context().Value(RouteParamsKey)
	if val == nil {
		return zero, fmt.Errorf("no route params in context")
	}

	routeParams, ok := val.(T)
	if !ok {
		return zero, fmt.Errorf("unexpected route params context type")
	}

	return routeParams, nil
}
