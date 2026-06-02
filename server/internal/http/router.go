package httpserver

import (
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"github.com/sitcon-tw/camp2026-game/internal/content"
	"github.com/sitcon-tw/camp2026-game/internal/http/authctx"
	authhandler "github.com/sitcon-tw/camp2026-game/internal/http/handler/auth"
	cataloghandler "github.com/sitcon-tw/camp2026-game/internal/http/handler/catalog"
	matcheshandler "github.com/sitcon-tw/camp2026-game/internal/http/handler/matches"
	mehandler "github.com/sitcon-tw/camp2026-game/internal/http/handler/me"
	systemhandler "github.com/sitcon-tw/camp2026-game/internal/http/handler/system"
	"github.com/sitcon-tw/camp2026-game/internal/http/httpx"
)

type Dependencies struct {
	Log            *slog.Logger
	RequestTimeout time.Duration
	Content        *content.Store
	MongoClient    *mongo.Client
	MongoDB        *mongo.Database
}

func NewRouter(dep Dependencies) http.Handler {
	if dep.Log == nil {
		dep.Log = slog.New(slog.NewTextHandler(io.Discard, nil))
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

	r.Route("/api", func(api chi.Router) {
		registerRoutes(api, dep)
	})

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		httpx.WriteProblem(w, r, httpx.NotFound("not found"))
	})
	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		httpx.WriteProblem(w, r, httpx.NewError(http.StatusMethodNotAllowed, "method not allowed"))
	})

	return r
}

func registerRoutes(api chi.Router, dep Dependencies) {
	systemhandler.New(systemhandler.Dependencies{
		MongoClient: dep.MongoClient,
	}).RegisterRoutes(api)
	authhandler.New(authhandler.Dependencies{
		MongoDB: dep.MongoDB,
	}).RegisterRoutes(api)
	cataloghandler.New(cataloghandler.Dependencies{
		Content: dep.Content,
	}).RegisterRoutes(api)
	mehandler.New(mehandler.Dependencies{
		Content: dep.Content,
		MongoDB: dep.MongoDB,
	}).RegisterRoutes(api.With(authctx.RequirePlayer(dep.MongoDB)))
	matcheshandler.New(matcheshandler.Dependencies{
		Content: dep.Content,
		MongoDB: dep.MongoDB,
	}).RegisterRoutes(api.With(authctx.RequirePlayer(dep.MongoDB)))

	registerSwaggerRoutes(api)
}
