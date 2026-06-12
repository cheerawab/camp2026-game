package authctx

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"github.com/sitcon-tw/camp2026-game/internal/http/httpx"
	mongomodel "github.com/sitcon-tw/camp2026-game/internal/mongodb/model"
)

const CookieName = "camp2026_auth"
const PlayerRoleStaff = "staff"

type playerContextKey struct{}

func WithPlayer(ctx context.Context, player mongomodel.Player) context.Context {
	return context.WithValue(ctx, playerContextKey{}, player)
}

func PlayerFromContext(ctx context.Context) (mongomodel.Player, bool) {
	player, ok := ctx.Value(playerContextKey{}).(mongomodel.Player)
	return player, ok
}

func RequirePlayer(db *mongo.Database) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if db == nil {
				httpx.WriteProblem(w, r, httpx.ServiceUnavailable("database is unavailable"))
				return
			}

			cookie, err := r.Cookie(CookieName)
			if err != nil || strings.TrimSpace(cookie.Value) == "" {
				httpx.WriteProblem(w, r, httpx.NewError(http.StatusUnauthorized, "authentication required"))
				return
			}

			var player mongomodel.Player
			err = db.Collection(mongomodel.PlayersCollection).
				FindOne(r.Context(), bson.M{"auth_token": strings.TrimSpace(cookie.Value)}).
				Decode(&player)
			if errors.Is(err, mongo.ErrNoDocuments) {
				httpx.WriteProblem(w, r, httpx.NewError(http.StatusUnauthorized, "authentication required"))
				return
			}
			if err != nil {
				httpx.WriteProblem(w, r, httpx.InternalServerError("authentication failed", "authentication_lookup_failed", err))
				return
			}
			if player.ID == "" {
				httpx.WriteProblem(w, r, httpx.InternalServerError("authentication failed", "authentication_player_invalid", errors.New("authenticated player has empty id")))
				return
			}

			next.ServeHTTP(w, r.WithContext(WithPlayer(r.Context(), player)))
		})
	}
}

func RequireStaff(db *mongo.Database) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return RequirePlayer(db)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			player, ok := PlayerFromContext(r.Context())
			if !ok || player.Role != PlayerRoleStaff {
				httpx.WriteProblem(w, r, httpx.NewError(http.StatusForbidden, "staff access required"))
				return
			}

			next.ServeHTTP(w, r)
		}))
	}
}
