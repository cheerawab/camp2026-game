package httpserver

import (
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/sitcon-tw/camp2026-game/internal/http/httpx"
)

var (
	globalRateLimit    = rateLimitConfig{Requests: 300, Window: time.Minute}
	authLoginRateLimit = rateLimitConfig{Requests: 10, Window: time.Minute}
)

type rateLimitConfig struct {
	Requests int
	Window   time.Duration
}

type ipRateLimiter struct {
	mu           sync.Mutex
	requests     int
	window       time.Duration
	now          func() time.Time
	clients      map[string]rateLimitBucket
	nextCleanup  time.Time
	cleanupEvery time.Duration
}

type rateLimitBucket struct {
	count   int
	resetAt time.Time
}

func newIPRateLimiter(config rateLimitConfig) *ipRateLimiter {
	if config.Requests <= 0 {
		panic("rate limit requests must be positive")
	}
	if config.Window <= 0 {
		panic("rate limit window must be positive")
	}

	cleanupEvery := config.Window
	if cleanupEvery > time.Minute {
		cleanupEvery = time.Minute
	}

	return &ipRateLimiter{
		requests:     config.Requests,
		window:       config.Window,
		now:          time.Now,
		clients:      make(map[string]rateLimitBucket),
		cleanupEvery: cleanupEvery,
	}
}

func (l *ipRateLimiter) allow(ip string) (bool, time.Duration) {
	now := l.now()

	l.mu.Lock()
	defer l.mu.Unlock()

	if l.nextCleanup.IsZero() || !now.Before(l.nextCleanup) {
		l.cleanup(now)
		l.nextCleanup = now.Add(l.cleanupEvery)
	}

	bucket := l.clients[ip]
	if bucket.resetAt.IsZero() || !now.Before(bucket.resetAt) {
		bucket = rateLimitBucket{resetAt: now.Add(l.window)}
	}
	if bucket.count >= l.requests {
		l.clients[ip] = bucket
		return false, bucket.resetAt.Sub(now)
	}

	bucket.count++
	l.clients[ip] = bucket
	return true, 0
}

func (l *ipRateLimiter) cleanup(now time.Time) {
	for ip, bucket := range l.clients {
		if !now.Before(bucket.resetAt) {
			delete(l.clients, ip)
		}
	}
}

func rateLimitByIP(limiter *ipRateLimiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if allowed, retryAfter := limiter.allow(clientIP(r)); !allowed {
				writeRateLimitExceeded(w, r, retryAfter)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func rateLimitRoute(method string, path string, limiter *ipRateLimiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == method && r.URL.Path == path {
				if allowed, retryAfter := limiter.allow(clientIP(r)); !allowed {
					writeRateLimitExceeded(w, r, retryAfter)
					return
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}

func writeRateLimitExceeded(w http.ResponseWriter, r *http.Request, retryAfter time.Duration) {
	w.Header().Set("Retry-After", retryAfterSeconds(retryAfter))
	httpx.WriteProblem(w, r, httpx.NewError(http.StatusTooManyRequests, "rate limit exceeded"))
}

func retryAfterSeconds(duration time.Duration) string {
	if duration <= 0 {
		return "1"
	}

	seconds := int((duration + time.Second - time.Nanosecond) / time.Second)
	if seconds < 1 {
		seconds = 1
	}
	return strconv.Itoa(seconds)
}

func clientIP(r *http.Request) string {
	remoteAddr := strings.TrimSpace(r.RemoteAddr)
	if remoteAddr == "" {
		return "unknown"
	}
	if ip := net.ParseIP(remoteAddr); ip != nil {
		return ip.String()
	}

	host, _, err := net.SplitHostPort(remoteAddr)
	if err == nil {
		host = strings.TrimSpace(host)
		if ip := net.ParseIP(host); ip != nil {
			return ip.String()
		}
		if host != "" {
			return host
		}
	}

	if ip := net.ParseIP(strings.Trim(remoteAddr, "[]")); ip != nil {
		return ip.String()
	}
	return remoteAddr
}
