package httpserver

import (
	"net"
	"net/http"
	"net/url"
	"strings"

	"github.com/sitcon-tw/camp2026-game/internal/http/httpx"
)

type requestOrigin struct {
	scheme string
	host   string
	port   string
}

func sameOriginUnsafeRequestGuard() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if isSafeMethod(r.Method) {
				next.ServeHTTP(w, r)
				return
			}

			if strings.EqualFold(strings.TrimSpace(r.Header.Get("Sec-Fetch-Site")), "cross-site") {
				rejectCrossOriginRequest(w, r)
				return
			}

			if origin := strings.TrimSpace(r.Header.Get("Origin")); origin != "" {
				if !requestMatchesOrigin(r, origin) {
					rejectCrossOriginRequest(w, r)
					return
				}
				next.ServeHTTP(w, r)
				return
			}

			if referer := strings.TrimSpace(r.Header.Get("Referer")); referer != "" {
				if !requestMatchesOrigin(r, referer) {
					rejectCrossOriginRequest(w, r)
					return
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}

func isSafeMethod(method string) bool {
	switch method {
	case http.MethodGet, http.MethodHead, http.MethodOptions, http.MethodTrace:
		return true
	default:
		return false
	}
}

func rejectCrossOriginRequest(w http.ResponseWriter, r *http.Request) {
	httpx.WriteProblem(w, r, httpx.NewError(http.StatusForbidden, "cross-site request forbidden"))
}

func requestMatchesOrigin(r *http.Request, value string) bool {
	candidate, ok := parseHeaderOrigin(value)
	if !ok {
		return false
	}

	for _, expected := range requestOrigins(r) {
		if originsMatch(candidate, expected) || localDevelopmentOriginsMatch(candidate, expected) {
			return true
		}
	}
	return false
}

func parseHeaderOrigin(value string) (requestOrigin, bool) {
	parsed, err := url.Parse(value)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" || parsed.User != nil {
		return requestOrigin{}, false
	}
	return normalizeOrigin(parsed.Scheme, parsed.Host)
}

func requestOrigins(r *http.Request) []requestOrigin {
	origins := make([]requestOrigin, 0, 2)

	// The edge proxy must overwrite these headers when the public origin differs from the backend origin.
	if forwardedHost := firstForwardedValue(r.Header.Get("X-Forwarded-Host")); forwardedHost != "" {
		if origin, ok := normalizeOrigin(requestScheme(r, firstForwardedValue(r.Header.Get("X-Forwarded-Proto"))), forwardedHost); ok {
			origins = append(origins, origin)
		}
	}

	if directHost := strings.TrimSpace(r.Host); directHost != "" {
		if origin, ok := normalizeOrigin(requestScheme(r, ""), directHost); ok {
			origins = append(origins, origin)
		}
	}

	return origins
}

func requestScheme(r *http.Request, forwardedProto string) string {
	if forwardedProto != "" {
		return forwardedProto
	}
	if r.TLS != nil {
		return "https"
	}
	if r.URL != nil && r.URL.Scheme != "" {
		return r.URL.Scheme
	}
	return "http"
}

func firstForwardedValue(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}
	if beforeComma, _, found := strings.Cut(value, ","); found {
		return strings.TrimSpace(beforeComma)
	}
	return value
}

func normalizeOrigin(scheme string, hostport string) (requestOrigin, bool) {
	scheme = strings.ToLower(strings.TrimSpace(scheme))
	hostport = strings.ToLower(strings.TrimSpace(hostport))
	if scheme == "" || hostport == "" {
		return requestOrigin{}, false
	}

	host, port := splitHostPort(hostport)
	if host == "" {
		return requestOrigin{}, false
	}
	if isDefaultPort(scheme, port) {
		port = ""
	}

	return requestOrigin{scheme: scheme, host: host, port: port}, true
}

func splitHostPort(hostport string) (string, string) {
	host, port, err := net.SplitHostPort(hostport)
	if err == nil {
		return strings.Trim(host, "[]"), port
	}
	return strings.Trim(hostport, "[]"), ""
}

func isDefaultPort(scheme string, port string) bool {
	return (scheme == "http" && port == "80") || (scheme == "https" && port == "443")
}

func originsMatch(a requestOrigin, b requestOrigin) bool {
	return a.scheme == b.scheme && a.host == b.host && a.port == b.port
}

func localDevelopmentOriginsMatch(a requestOrigin, b requestOrigin) bool {
	return a.scheme == b.scheme && isLoopbackHost(a.host) && isLoopbackHost(b.host)
}

func isLoopbackHost(host string) bool {
	if host == "localhost" || strings.HasSuffix(host, ".localhost") {
		return true
	}
	ip := net.ParseIP(host)
	return ip != nil && ip.IsLoopback()
}
