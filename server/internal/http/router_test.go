package httpserver

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRoot(t *testing.T) {
	router := NewRouter(Dependencies{APIVersion: "v1"})

	res := performRequest(router, http.MethodGet, "/api/v1/", nil)
	if res.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, res.Code)
	}

	var body struct {
		Message string `json:"message"`
	}
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body.Message == "" {
		t.Fatal("expected message")
	}
}

func TestHealth(t *testing.T) {
	router := NewRouter(Dependencies{APIVersion: "v1"})

	res := performRequest(router, http.MethodGet, "/api/v1/healthz", nil)
	if res.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, res.Code)
	}
}

func TestReady(t *testing.T) {
	router := NewRouter(Dependencies{
		APIVersion: "v1",
		ReadinessCheck: func(context.Context) error {
			return nil
		},
	})

	res := performRequest(router, http.MethodGet, "/api/v1/readyz", nil)
	if res.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, res.Code)
	}
}

func TestReadyFailure(t *testing.T) {
	router := NewRouter(Dependencies{
		APIVersion: "v1",
		ReadinessCheck: func(context.Context) error {
			return errors.New("database unavailable")
		},
	})

	res := performRequest(router, http.MethodGet, "/api/v1/readyz", nil)
	if res.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected status %d, got %d", http.StatusServiceUnavailable, res.Code)
	}
	if contentType := res.Header().Get("Content-Type"); contentType != "application/problem+json" {
		t.Fatalf("expected problem content type, got %q", contentType)
	}
}

func TestSwaggerJSON(t *testing.T) {
	router := NewRouter(Dependencies{APIVersion: "v1"})

	res := performRequest(router, http.MethodGet, "/api/v1/swagger.json", nil)
	if res.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, res.Code)
	}
	if contentType := res.Header().Get("Content-Type"); contentType != "application/json" {
		t.Fatalf("expected json content type, got %q", contentType)
	}
}

func TestScalarDocs(t *testing.T) {
	router := NewRouter(Dependencies{APIVersion: "v1"})

	res := performRequest(router, http.MethodGet, "/api/v1/docs/index.html", nil)
	if res.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, res.Code)
	}
	if contentType := res.Header().Get("Content-Type"); contentType != "text/html; charset=utf-8" {
		t.Fatalf("expected html content type, got %q", contentType)
	}

	body := res.Body.String()
	for _, want := range []string{
		"@scalar/api-reference",
		"Scalar.createApiReference",
		"../swagger.json",
	} {
		if !strings.Contains(body, want) {
			t.Fatalf("expected docs html to contain %q", want)
		}
	}
}

func TestValidationExample(t *testing.T) {
	router := NewRouter(Dependencies{APIVersion: "v1"})

	res := performRequest(router, http.MethodPost, "/api/v1/examples/validation", strings.NewReader(`{
		"players": [
			{
				"displayName": " Alice ",
				"teamNumber": 3,
				"favoritePebbleType": "engineering"
			}
		]
	}`))
	if res.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d: %s", http.StatusCreated, res.Code, res.Body.String())
	}

	var body struct {
		Message string `json:"message"`
		Players []struct {
			DisplayName        string `json:"displayName"`
			TeamNumber         int    `json:"teamNumber"`
			FavoritePebbleType string `json:"favoritePebbleType"`
		} `json:"players"`
	}
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body.Message != "validation example accepted" {
		t.Fatalf("unexpected message %q", body.Message)
	}
	if got := body.Players[0].DisplayName; got != "Alice" {
		t.Fatalf("expected trimmed displayName, got %q", got)
	}
}

func TestValidationExampleMalformedJSON(t *testing.T) {
	router := NewRouter(Dependencies{APIVersion: "v1"})

	res := performRequest(router, http.MethodPost, "/api/v1/examples/validation", strings.NewReader(`{"players": [`))
	assertProblem(t, res, http.StatusBadRequest, "")
}

func TestValidationExampleUnknownField(t *testing.T) {
	router := NewRouter(Dependencies{APIVersion: "v1"})

	res := performRequest(router, http.MethodPost, "/api/v1/examples/validation", strings.NewReader(`{
		"players": [],
		"unexpected": true
	}`))
	assertProblem(t, res, http.StatusBadRequest, "")
}

func TestValidationExampleMissingPlayers(t *testing.T) {
	router := NewRouter(Dependencies{APIVersion: "v1"})

	res := performRequest(router, http.MethodPost, "/api/v1/examples/validation", strings.NewReader(`{}`))
	problem := assertProblem(t, res, http.StatusUnprocessableEntity, "body.players")
	if got := problem.Errors[0].Message; got != "players is required" {
		t.Fatalf("unexpected validation message %q", got)
	}
}

func TestValidationExampleEmptyPlayers(t *testing.T) {
	router := NewRouter(Dependencies{APIVersion: "v1"})

	res := performRequest(router, http.MethodPost, "/api/v1/examples/validation", strings.NewReader(`{"players":[]}`))
	problem := assertProblem(t, res, http.StatusUnprocessableEntity, "body.players")
	if got := problem.Errors[0].Message; got != "players must be at least 1" {
		t.Fatalf("unexpected validation message %q", got)
	}
}

func TestValidationExampleTooManyPlayers(t *testing.T) {
	router := NewRouter(Dependencies{APIVersion: "v1"})

	res := performRequest(router, http.MethodPost, "/api/v1/examples/validation", strings.NewReader(validationExampleBodyWithPlayers(21)))
	problem := assertProblem(t, res, http.StatusUnprocessableEntity, "body.players")
	if got := problem.Errors[0].Message; got != "players must be at most 20" {
		t.Fatalf("unexpected validation message %q", got)
	}
}

func TestValidationExampleNestedDisplayName(t *testing.T) {
	router := NewRouter(Dependencies{APIVersion: "v1"})

	res := performRequest(router, http.MethodPost, "/api/v1/examples/validation", strings.NewReader(`{
		"players": [
			{
				"displayName": "   ",
				"teamNumber": 3,
				"favoritePebbleType": "engineering"
			}
		]
	}`))
	problem := assertProblem(t, res, http.StatusUnprocessableEntity, "body.players[0].displayName")
	if got := problem.Errors[0].Message; got != "displayName is required" {
		t.Fatalf("unexpected validation message %q", got)
	}
}

func TestValidationExampleNestedTeamNumber(t *testing.T) {
	router := NewRouter(Dependencies{APIVersion: "v1"})

	res := performRequest(router, http.MethodPost, "/api/v1/examples/validation", strings.NewReader(`{
		"players": [
			{
				"displayName": "Alice",
				"teamNumber": 3,
				"favoritePebbleType": "engineering"
			},
			{
				"displayName": "Bob",
				"teamNumber": 21,
				"favoritePebbleType": "exploration"
			}
		]
	}`))
	problem := assertProblem(t, res, http.StatusUnprocessableEntity, "body.players[1].teamNumber")
	if got := problem.Errors[0].Message; got != "teamNumber must be at most 20" {
		t.Fatalf("unexpected validation message %q", got)
	}
}

func TestValidationExampleFavoritePebbleType(t *testing.T) {
	router := NewRouter(Dependencies{APIVersion: "v1"})

	res := performRequest(router, http.MethodPost, "/api/v1/examples/validation", strings.NewReader(`{
		"players": [
			{
				"displayName": "Alice",
				"teamNumber": 3,
				"favoritePebbleType": "unknown"
			}
		]
	}`))
	problem := assertProblem(t, res, http.StatusUnprocessableEntity, "body.players[0].favoritePebbleType")
	if !strings.Contains(problem.Errors[0].Message, "favoritePebbleType must be one of:") {
		t.Fatalf("unexpected validation message %q", problem.Errors[0].Message)
	}
}

func TestNotFound(t *testing.T) {
	router := NewRouter(Dependencies{APIVersion: "v1"})

	res := performRequest(router, http.MethodGet, "/missing", nil)
	if res.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, res.Code)
	}
}

type problemResponse struct {
	Status int `json:"status"`
	Errors []struct {
		Location string `json:"location"`
		Message  string `json:"message"`
	} `json:"errors"`
}

func assertProblem(t *testing.T, res *httptest.ResponseRecorder, status int, location string) problemResponse {
	t.Helper()

	if res.Code != status {
		t.Fatalf("expected status %d, got %d: %s", status, res.Code, res.Body.String())
	}
	if contentType := res.Header().Get("Content-Type"); contentType != "application/problem+json" {
		t.Fatalf("expected problem content type, got %q", contentType)
	}

	var problem problemResponse
	if err := json.NewDecoder(res.Body).Decode(&problem); err != nil {
		t.Fatalf("decode problem: %v", err)
	}
	if problem.Status != status {
		t.Fatalf("expected problem status %d, got %d", status, problem.Status)
	}
	if location != "" {
		if len(problem.Errors) == 0 {
			t.Fatalf("expected validation errors")
		}
		if got := problem.Errors[0].Location; got != location {
			t.Fatalf("expected error location %q, got %q", location, got)
		}
	}
	return problem
}

func validationExampleBodyWithPlayers(count int) string {
	var builder strings.Builder
	builder.WriteString(`{"players":[`)
	for i := range count {
		if i > 0 {
			builder.WriteString(",")
		}
		builder.WriteString(`{"displayName":"Player","teamNumber":1,"favoritePebbleType":"engineering"}`)
	}
	builder.WriteString(`]}`)
	return builder.String()
}

func performRequest(handler http.Handler, method, path string, body *strings.Reader) *httptest.ResponseRecorder {
	if body == nil {
		body = strings.NewReader("")
	}
	req := httptest.NewRequest(method, path, body)
	res := httptest.NewRecorder()
	handler.ServeHTTP(res, req)
	return res
}
