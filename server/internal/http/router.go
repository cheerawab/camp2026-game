package httpserver

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"

	examplehandler "github.com/sitcon-tw/camp2026-game/internal/http/handler/example"
	systemhandler "github.com/sitcon-tw/camp2026-game/internal/http/handler/system"
	"github.com/sitcon-tw/camp2026-game/internal/http/httpx"
	"github.com/sitcon-tw/camp2026-game/internal/postgres/sqlc"
)

type Dependencies struct {
	Log            *slog.Logger
	APIVersion     string
	RequestTimeout time.Duration
	ReadinessCheck func(context.Context) error
	Queries        *sqlc.Queries
}

func NewRouter(dep Dependencies) http.Handler {
	if dep.Log == nil {
		dep.Log = slog.New(slog.NewTextHandler(io.Discard, nil))
	}
	if dep.APIVersion == "" {
		dep.APIVersion = "v1"
	}
	if dep.RequestTimeout <= 0 {
		dep.RequestTimeout = 10 * time.Second
	}

	r := chi.NewRouter()
	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)
	r.Use(chimw.StripSlashes)
	r.Use(chimw.Timeout(dep.RequestTimeout))
	r.Use(recoverer(dep.Log))
	r.Use(requestLogger(dep.Log))

	r.Route("/api/"+dep.APIVersion, func(api chi.Router) {
		registerRoutes(api, dep, dep.APIVersion)
	})

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		httpx.WriteProblem(w, r, httpx.NotFound("not found"))
	})
	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		httpx.WriteProblem(w, r, httpx.NewError(http.StatusMethodNotAllowed, "method not allowed"))
	})

	return r
}

func registerRoutes(api chi.Router, dep Dependencies, apiVersion string) {
	systemhandler.New(systemhandler.Dependencies{
		ReadinessCheck: dep.ReadinessCheck,
		Queries:        dep.Queries,
	}).RegisterRoutes(api)
	examplehandler.New().RegisterRoutes(api)

	registerSwaggerRoutes(api, apiVersion)
}
