package me

import (
	"context"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"github.com/sitcon-tw/camp2026-game/internal/http/httpx"
	mongomodel "github.com/sitcon-tw/camp2026-game/internal/mongodb/model"
)

// ListCompletedMatches godoc
// @Summary List current player completed matches
// @Description Returns completed quiz matches joined by the authenticated player. Waiting and active matches are not returned.
// @Tags me
// @Produce json
// @Security AuthCookieAuth
// @Success 200 {object} CompletedMatchListResponse
// @Failure 401 {object} httpx.ProblemDetails
// @Failure 500 {object} httpx.ProblemDetails
// @Failure 503 {object} httpx.ProblemDetails
// @Router /me/matches [get]
func (h *Handler) ListCompletedMatches(w http.ResponseWriter, r *http.Request) {
	player, ok := currentPlayer(w, r)
	if !ok {
		return
	}
	if h.db == nil {
		httpx.WriteProblem(w, r, httpx.ServiceUnavailable("database is unavailable"))
		return
	}

	matches, err := h.findCompletedMatches(r.Context(), player.ID)
	if err != nil {
		httpx.WriteProblem(w, r, httpx.NewError(http.StatusInternalServerError, "matches unavailable"))
		return
	}

	httpx.WriteJSON(w, http.StatusOK, CompletedMatchListResponse{
		Matches: mapCompletedMatches(matches),
	})
}

func (h *Handler) findCompletedMatches(ctx context.Context, playerID string) ([]mongomodel.Match, error) {
	cursor, err := h.db.Collection(mongomodel.MatchesCollection).Find(
		ctx,
		completedMatchesFilter(playerID),
		options.Find().SetSort(bson.D{
			{Key: "completed_at", Value: -1},
			{Key: "created_at", Value: -1},
		}),
	)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = cursor.Close(ctx)
	}()

	var matches []mongomodel.Match
	if err := cursor.All(ctx, &matches); err != nil {
		return nil, err
	}
	return matches, nil
}

func completedMatchesFilter(playerID string) bson.D {
	return bson.D{
		{Key: "status", Value: mongomodel.MatchStatusCompleted},
		{Key: "players.player_id", Value: playerID},
	}
}

func mapCompletedMatches(records []mongomodel.Match) []CompletedMatchResponse {
	matches := make([]CompletedMatchResponse, 0, len(records))
	for _, record := range records {
		matches = append(matches, CompletedMatchResponse{
			MatchID:       record.ID,
			Status:        record.Status,
			HostPlayerID:  record.HostPlayerID,
			Players:       mapCompletedMatchPlayers(record.Players),
			QuestionCount: len(record.QuestionIDs),
			CreatedAt:     record.CreatedAt,
			StartedAt:     optionalTime(record.StartedAt),
			CompletedAt:   optionalTime(record.CompletedAt),
		})
	}
	return matches
}

func mapCompletedMatchPlayers(records []mongomodel.MatchPlayer) []CompletedMatchPlayerResponse {
	players := make([]CompletedMatchPlayerResponse, 0, len(records))
	for _, record := range records {
		players = append(players, CompletedMatchPlayerResponse{
			PlayerID:  record.PlayerID,
			Nickname:  record.Nickname,
			SitoneIDs: cloneStringSlice(record.SitoneIDs),
			Score:     record.Score,
		})
	}
	return players
}

func optionalTime(value time.Time) *time.Time {
	if value.IsZero() {
		return nil
	}
	return &value
}

func cloneStringSlice(input []string) []string {
	if input == nil {
		return nil
	}
	return append([]string(nil), input...)
}
