package httpx

import (
	"bytes"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	chimw "github.com/go-chi/chi/v5/middleware"
)

func TestWriteProblemLogsInternalCauseWithoutExposingIt(t *testing.T) {
	var logs bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&logs, nil))
	handler := chimw.RequestID(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := WithLogger(r.Context(), logger)
		WriteProblem(w, r.WithContext(ctx), InternalServerError("login failed", "login_player_lookup_failed", errors.New("mongo connection refused")))
	}))

	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", nil)
	res := httptest.NewRecorder()
	handler.ServeHTTP(res, req)

	if res.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, res.Code)
	}
	var problem ProblemDetails
	if err := json.NewDecoder(res.Body).Decode(&problem); err != nil {
		t.Fatalf("decode problem: %v", err)
	}
	if problem.Detail != "login failed" {
		t.Fatalf("expected client detail to stay generic, got %q", problem.Detail)
	}
	if problem.Code != "login_player_lookup_failed" {
		t.Fatalf("expected error code, got %q", problem.Code)
	}
	if problem.RequestID == "" {
		t.Fatal("expected request id in problem response")
	}
	if strings.Contains(res.Body.String(), "mongo connection refused") {
		t.Fatalf("expected client problem to omit cause, got %s", res.Body.String())
	}

	logOutput := logs.String()
	for _, want := range []string{
		`"msg":"http problem"`,
		`"code":"login_player_lookup_failed"`,
		`"request_id":"` + problem.RequestID + `"`,
		`"error":"mongo connection refused"`,
	} {
		if !strings.Contains(logOutput, want) {
			t.Fatalf("expected log to contain %s, got %s", want, logOutput)
		}
	}
}

func TestWriteProblemSkipsBareInternalHTTPErrorLog(t *testing.T) {
	var logs bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&logs, nil))
	req := httptest.NewRequest(http.MethodGet, "/api/panic", nil)
	req = req.WithContext(WithLogger(req.Context(), logger))
	res := httptest.NewRecorder()

	WriteProblem(res, req, NewError(http.StatusInternalServerError, "internal server error"))

	var problem ProblemDetails
	if err := json.NewDecoder(res.Body).Decode(&problem); err != nil {
		t.Fatalf("decode problem: %v", err)
	}
	if problem.Code != "internal_server_error" {
		t.Fatalf("expected default 5xx code, got %q", problem.Code)
	}
	if logs.Len() != 0 {
		t.Fatalf("expected bare internal http error to skip duplicate log, got %s", logs.String())
	}
}

func TestWriteProblemOmitsTrackingFieldsForClientError(t *testing.T) {
	handler := chimw.RequestID(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		WriteProblem(w, r, NewError(http.StatusUnauthorized, "authentication required"))
	}))

	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", nil)
	res := httptest.NewRecorder()
	handler.ServeHTTP(res, req)

	var problem ProblemDetails
	if err := json.NewDecoder(res.Body).Decode(&problem); err != nil {
		t.Fatalf("decode problem: %v", err)
	}
	if problem.Code != "" {
		t.Fatalf("expected client error to omit code, got %q", problem.Code)
	}
	if problem.RequestID != "" {
		t.Fatalf("expected client error to omit request id, got %q", problem.RequestID)
	}
}
