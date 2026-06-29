package httpserver

import (
	"net/http"
	"testing"
	"time"
)

func TestIPRateLimiterAllowsUntilLimitThenResets(t *testing.T) {
	now := time.Date(2026, 6, 29, 12, 0, 0, 0, time.UTC)
	limiter := newIPRateLimiter(rateLimitConfig{Requests: 2, Window: time.Minute})
	limiter.now = func() time.Time { return now }

	if allowed, _ := limiter.allow("203.0.113.10"); !allowed {
		t.Fatal("expected first request to be allowed")
	}
	if allowed, _ := limiter.allow("203.0.113.10"); !allowed {
		t.Fatal("expected second request to be allowed")
	}
	allowed, retryAfter := limiter.allow("203.0.113.10")
	if allowed {
		t.Fatal("expected third request to be rate limited")
	}
	if retryAfter != time.Minute {
		t.Fatalf("expected retry after %s, got %s", time.Minute, retryAfter)
	}

	now = now.Add(time.Minute)
	if allowed, _ := limiter.allow("203.0.113.10"); !allowed {
		t.Fatal("expected request to be allowed after reset")
	}
}

func TestIPRateLimiterTracksIPsIndependently(t *testing.T) {
	limiter := newIPRateLimiter(rateLimitConfig{Requests: 1, Window: time.Minute})

	if allowed, _ := limiter.allow("203.0.113.10"); !allowed {
		t.Fatal("expected first IP to be allowed")
	}
	if allowed, _ := limiter.allow("203.0.113.10"); allowed {
		t.Fatal("expected first IP to be rate limited")
	}
	if allowed, _ := limiter.allow("203.0.113.11"); !allowed {
		t.Fatal("expected second IP to be allowed")
	}
}

func TestClientIPParsesRemoteAddr(t *testing.T) {
	tests := []struct {
		name       string
		remoteAddr string
		want       string
	}{
		{name: "host port", remoteAddr: "203.0.113.10:49152", want: "203.0.113.10"},
		{name: "bare IPv4", remoteAddr: "203.0.113.10", want: "203.0.113.10"},
		{name: "IPv6 host port", remoteAddr: "[2001:db8::1]:49152", want: "2001:db8::1"},
		{name: "bare IPv6", remoteAddr: "2001:db8::1", want: "2001:db8::1"},
		{name: "empty", remoteAddr: "", want: "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &http.Request{RemoteAddr: tt.remoteAddr}
			if got := clientIP(req); got != tt.want {
				t.Fatalf("expected client IP %q, got %q", tt.want, got)
			}
		})
	}
}
