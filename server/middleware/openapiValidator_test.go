package middleware

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/routers"
	"github.com/getkin/kin-openapi/routers/legacy"
	"go.uber.org/zap"
)

func TestPeekBody(t *testing.T) {
	body := []byte(`{"foo":"bar"}`)
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	logger := zap.NewNop()

	got, err := PeekBody(logger, req)
	if err != nil {
		t.Fatalf("PeekBody returned error: %v", err)
	}
	if !bytes.Equal(got, body) {
		t.Errorf("PeekBody returned %s, want %s", got, body)
	}

	// Body should be readable again
	got2, err := io.ReadAll(req.Body)
	if err != nil {
		t.Fatalf("Reading body again failed: %v", err)
	}
	if !bytes.Equal(got2, body) {
		t.Errorf("Body after PeekBody = %s, want %s", got2, body)
	}
}

func TestPeekResponse(t *testing.T) {
	body := []byte(`{"foo":"bar"}`)
	resp := &http.Response{
		Body: io.NopCloser(bytes.NewReader(body)),
	}
	logger := zap.NewNop()

	got, err := PeekResponse(logger, resp)
	if err != nil {
		t.Fatalf("PeekResponse returned error: %v", err)
	}
	if !bytes.Equal(got, body) {
		t.Errorf("PeekResponse returned %s, want %s", got, body)
	}

	// Body should be readable again
	got2, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Reading body again failed: %v", err)
	}
	if !bytes.Equal(got2, body) {
		t.Errorf("Body after PeekResponse = %s, want %s", got2, body)
	}
}

type dummyRouter struct {
	route  routers.Route
	params map[string]string
	err    error
}

func (d *dummyRouter) FindRoute(r *http.Request) (*routers.Route, map[string]string, error) {
	return &d.route, d.params, d.err
}

func TestOpenAPIValidationMiddleware_RouteNotFound(t *testing.T) {
	router := &dummyRouter{err: errors.New("not found")}
	mw := OpenAPIValidationMiddleware(router)
	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Handler should not be called on route not found")
	}))
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/notfound", nil)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusNotFound {
		t.Errorf("Expected 404, got %d", rr.Code)
	}
}

func TestOpenAPIValidationMiddleware_SetsContextValues(t *testing.T) {
	// Minimal OpenAPI spec for validation
	spec := `{
		"openapi":"3.0.0",
		"info":{
			"title":"Test API",
			"version":"1.0.0"
		},
		"paths":{
			"/":{
				"post":{
					"requestBody":{
						"content":{
							"application/json":{
								"schema":{
									"type":"object",
									"properties":{
										"foo":{
											"type":"string"
										}
									},
									"required":["foo"]
								}
							}
						}
					},
					"responses": {
        				"200": {
          					"description": "OK",
          					"content": {
								"application/json": {
									"schema": {
										"type": "object"
									}
								}
							}
						}
					}
				}
			}
		}
	}`
	if doc, err := openapi3.NewLoader().LoadFromData([]byte(spec)); err != nil {
		t.Fatalf("Failed to load OpenAPI spec: %v", err)
		return
	} else if err := doc.Validate(openapi3.NewLoader().Context); err != nil {
		t.Fatalf("Failed to validate OpenAPI spec: %v", err)
		return
	} else {
		router, _ := legacy.NewRouter(doc)
		mw := OpenAPIValidationMiddleware(router)
		called := false
		handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Printf("Handler called with method: %s, URL: %s\n", r.Method, r.URL)
			called = true
			// Body
			err := r.Context().Err()
			if err != nil {
				t.Errorf("Unexpected context error: %v", err)
			}
			b, err := GetValidatedBody[map[string]any](r)
			if err != nil {
				t.Errorf("GetValidatedBody error: %v", err)
			}
			if b["foo"] != "bar" {
				t.Errorf("Expected foo=bar, got %v", b["foo"])
			}
			// Query
			q, err := GetQueryParams[map[string]any](r)
			if err != nil {
				t.Errorf("GetQueryParams error: %v", err)
			}
			if q["baz"] != "qux" {
				t.Errorf("Expected baz=qux, got %v", q["baz"])
			}
			// Route
			route, err := GetRoute[routers.Route](r)
			if err != nil {
				t.Errorf("GetRoute error: %v", err)
			}
			if route.Operation == nil {
				t.Errorf("Expected route.Operation to be set")
			}
			// RouteParams
			params, err := GetRouteParams[map[string]string](r)
			if err != nil {
				t.Errorf("GetRouteParams error: %v", err)
			}
			if len(params) != 0 {
				t.Errorf("Expected empty route params, got %v", params)
			}
		}))
		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/?baz=qux", bytes.NewReader([]byte(`{"foo":"bar"}`)))
		req.Header.Set("Content-Type", "application/json")
		handler.ServeHTTP(rr, req)
		if !called {
			t.Error("Handler was not called")
		}
	}
}

func TestGetValidatedBody_NoBody(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	_, err := GetValidatedBody[map[string]any](req)
	if err == nil {
		t.Error("Expected error when no validated body in context")
	}
}

func TestGetValidatedBody_TypeMismatch(t *testing.T) {
	ctx := context.WithValue(context.Background(), BodyKey, "not a map")
	req := httptest.NewRequest(http.MethodGet, "/", nil).WithContext(ctx)
	_, err := GetValidatedBody[map[string]any](req)
	if err == nil {
		t.Error("Expected error on type mismatch")
	}
}

func TestGetQueryParams_NoQuery(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	_, err := GetQueryParams[map[string]any](req)
	if err == nil {
		t.Error("Expected error when no query params in context")
	}
}

func TestGetQueryParams_TypeMismatch(t *testing.T) {
	ctx := context.WithValue(context.Background(), QueryKey, "not a map")
	req := httptest.NewRequest(http.MethodGet, "/", nil).WithContext(ctx)
	_, err := GetQueryParams[map[string]any](req)
	if err == nil {
		t.Error("Expected error on query type mismatch")
	}
}

func TestGetRoute_NoRoute(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	_, err := GetRoute[routers.Route](req)
	if err == nil {
		t.Error("Expected error when no route in context")
	}
}

func TestGetRoute_TypeMismatch(t *testing.T) {
	ctx := context.WithValue(context.Background(), RouteKey, "not a route")
	req := httptest.NewRequest(http.MethodGet, "/", nil).WithContext(ctx)
	_, err := GetRoute[routers.Route](req)
	if err == nil {
		t.Error("Expected error on route type mismatch")
	}
}

func TestGetRouteParams_NoParams(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	_, err := GetRouteParams[map[string]string](req)
	if err == nil {
		t.Error("Expected error when no route params in context")
	}
}

func TestGetRouteParams_TypeMismatch(t *testing.T) {
	ctx := context.WithValue(context.Background(), RouteParamsKey, "not a map")
	req := httptest.NewRequest(http.MethodGet, "/", nil).WithContext(ctx)
	_, err := GetRouteParams[map[string]string](req)
	if err == nil {
		t.Error("Expected error on route params type mismatch")
	}
}
