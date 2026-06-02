package authctx

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"github.com/sitcon-tw/camp2026-game/internal/http/httpx"
	mongomodel "github.com/sitcon-tw/camp2026-game/internal/mongodb/model"
)

func TestPlayerContext(t *testing.T) {
	player := mongomodel.Player{ID: "7H9K2Q"}
	ctx := WithPlayer(t.Context(), player)

	got, ok := PlayerFromContext(ctx)
	if !ok {
		t.Fatal("expected player in context")
	}
	if got.ID != player.ID {
		t.Fatalf("expected player id %q, got %q", player.ID, got.ID)
	}
}

func TestRequirePlayerRequiresDatabase(t *testing.T) {
	handler := RequirePlayer(nil)(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		t.Fatal("next handler should not run")
	}))

	req := httptest.NewRequest(http.MethodGet, "/api/me/status", nil)
	req.AddCookie(&http.Cookie{Name: CookieName, Value: "auth_token_123456"})
	res := httptest.NewRecorder()
	handler.ServeHTTP(res, req)

	problem := assertProblem(t, res, http.StatusServiceUnavailable)
	if problem.Detail != "database is unavailable" {
		t.Fatalf("expected database unavailable detail, got %q", problem.Detail)
	}
}

func TestRequirePlayerRequiresCookie(t *testing.T) {
	handler := RequirePlayer(fakeDatabase(t))(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		t.Fatal("next handler should not run")
	}))

	req := httptest.NewRequest(http.MethodGet, "/api/me/status", nil)
	res := httptest.NewRecorder()
	handler.ServeHTTP(res, req)

	problem := assertProblem(t, res, http.StatusUnauthorized)
	if problem.Detail != "authentication required" {
		t.Fatalf("expected authentication required detail, got %q", problem.Detail)
	}
}

func TestRequirePlayerRequiresNonEmptyCookie(t *testing.T) {
	handler := RequirePlayer(fakeDatabase(t))(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		t.Fatal("next handler should not run")
	}))

	req := httptest.NewRequest(http.MethodGet, "/api/me/status", nil)
	req.AddCookie(&http.Cookie{Name: CookieName, Value: strings.Repeat(" ", 1)})
	res := httptest.NewRecorder()
	handler.ServeHTTP(res, req)

	assertProblem(t, res, http.StatusUnauthorized)
}

func fakeDatabase(t *testing.T) *mongo.Database {
	t.Helper()

	client, err := mongo.Connect(options.Client().
		ApplyURI("mongodb://127.0.0.1:1").
		SetServerSelectionTimeout(time.Millisecond))
	if err != nil {
		t.Fatalf("connect mongo client: %v", err)
	}
	t.Cleanup(func() {
		_ = client.Disconnect(t.Context())
	})
	return client.Database("camp2026_game_test")
}

func assertProblem(t *testing.T, res *httptest.ResponseRecorder, status int) httpx.ProblemDetails {
	t.Helper()

	if res.Code != status {
		t.Fatalf("expected status %d, got %d: %s", status, res.Code, res.Body.String())
	}
	if contentType := res.Header().Get("Content-Type"); contentType != "application/problem+json" {
		t.Fatalf("expected problem content type, got %q", contentType)
	}

	var problem httpx.ProblemDetails
	if err := json.NewDecoder(res.Body).Decode(&problem); err != nil {
		t.Fatalf("decode problem: %v", err)
	}
	if problem.Status != status {
		t.Fatalf("expected problem status %d, got %d", status, problem.Status)
	}
	return problem
}
