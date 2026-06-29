package httpserver

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSameOriginUnsafeRequestGuard(t *testing.T) {
	for _, tt := range []struct {
		name       string
		method     string
		target     string
		headers    map[string]string
		wantStatus int
	}{
		{
			name:       "allows same-origin unsafe request",
			method:     http.MethodPost,
			target:     "https://game.example/api/matches",
			headers:    map[string]string{"Origin": "https://game.example"},
			wantStatus: http.StatusNoContent,
		},
		{
			name:       "allows same-origin unsafe request with default port",
			method:     http.MethodPost,
			target:     "https://game.example/api/matches",
			headers:    map[string]string{"Origin": "https://game.example:443"},
			wantStatus: http.StatusNoContent,
		},
		{
			name:       "rejects cross-origin unsafe request",
			method:     http.MethodPost,
			target:     "https://game.example/api/matches",
			headers:    map[string]string{"Origin": "https://evil.example"},
			wantStatus: http.StatusForbidden,
		},
		{
			name:       "allows same-origin referer fallback",
			method:     http.MethodPut,
			target:     "https://game.example/api/me/sitone-loadout",
			headers:    map[string]string{"Referer": "https://game.example/stones"},
			wantStatus: http.StatusNoContent,
		},
		{
			name:       "rejects cross-site fetch metadata",
			method:     http.MethodPost,
			target:     "https://game.example/api/matches",
			headers:    map[string]string{"Sec-Fetch-Site": "cross-site"},
			wantStatus: http.StatusForbidden,
		},
		{
			name:       "allows safe method with cross-origin header",
			method:     http.MethodGet,
			target:     "https://game.example/api/me/status",
			headers:    map[string]string{"Origin": "https://evil.example"},
			wantStatus: http.StatusNoContent,
		},
		{
			name:       "allows non-browser clients without origin headers",
			method:     http.MethodPost,
			target:     "https://game.example/api/matches",
			wantStatus: http.StatusNoContent,
		},
		{
			name:   "allows production proxy forwarded origin",
			method: http.MethodPost,
			target: "http://backend:8080/api/matches",
			headers: map[string]string{
				"Origin":            "https://game.example",
				"X-Forwarded-Host":  "game.example",
				"X-Forwarded-Proto": "https",
			},
			wantStatus: http.StatusNoContent,
		},
		{
			name:       "allows local development proxy ports",
			method:     http.MethodPost,
			target:     "http://localhost:8080/api/matches",
			headers:    map[string]string{"Origin": "http://localhost:3000"},
			wantStatus: http.StatusNoContent,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			handler := sameOriginUnsafeRequestGuard()(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusNoContent)
			}))
			req := httptest.NewRequest(tt.method, tt.target, strings.NewReader(""))
			for key, value := range tt.headers {
				req.Header.Set(key, value)
			}
			res := httptest.NewRecorder()

			handler.ServeHTTP(res, req)

			if res.Code != tt.wantStatus {
				t.Fatalf("expected status %d, got %d: %s", tt.wantStatus, res.Code, res.Body.String())
			}
		})
	}
}

func TestAPIRouterRejectsCrossOriginUnsafeRequestsBeforeAuth(t *testing.T) {
	router := NewRouter(Dependencies{})
	req := httptest.NewRequest(http.MethodPost, "https://game.example/api/auth/logout", strings.NewReader(""))
	req.Header.Set("Origin", "https://evil.example")
	req.AddCookie(&http.Cookie{Name: "camp2026_auth", Value: "auth_token_123456"})
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Code != http.StatusForbidden {
		t.Fatalf("expected status %d, got %d: %s", http.StatusForbidden, res.Code, res.Body.String())
	}
	if contentType := res.Header().Get("Content-Type"); contentType != "application/problem+json" {
		t.Fatalf("expected problem content type, got %q", contentType)
	}
}

func TestAPIRouterAllowsSameOriginUnsafeRequests(t *testing.T) {
	router := NewRouter(Dependencies{})
	req := httptest.NewRequest(http.MethodPost, "https://game.example/api/auth/logout", strings.NewReader(""))
	req.Header.Set("Origin", "https://game.example")
	req.AddCookie(&http.Cookie{Name: "camp2026_auth", Value: "auth_token_123456"})
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Code == http.StatusForbidden {
		t.Fatalf("expected same-origin request to pass CSRF guard, got %d: %s", res.Code, res.Body.String())
	}
}
