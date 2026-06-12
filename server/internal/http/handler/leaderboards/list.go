package leaderboards

import (
	"context"
	"net/http"
	"sort"
	"strings"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"github.com/sitcon-tw/camp2026-game/internal/http/httpx"
	mongomodel "github.com/sitcon-tw/camp2026-game/internal/mongodb/model"
)

const (
	TypeOpenPower = "open_power"
	TypeSitones   = "sitones"
	TypeMatches   = "matches"
)

// List godoc
// @Summary List team leaderboard
// @Description Lists team rankings by open power, sitone count, or completed match score.
// @Tags leaderboards
// @Produce json
// @Security AuthCookieAuth
// @Param type query string false "Leaderboard type" Enums(open_power,sitones,matches)
// @Success 200 {object} ListResponse
// @Failure 401 {object} httpx.ProblemDetails
// @Failure 422 {object} httpx.ProblemDetails
// @Failure 500 {object} httpx.ProblemDetails
// @Failure 503 {object} httpx.ProblemDetails
// @Router /leaderboards [get]
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	player, ok := currentPlayer(w, r)
	if !ok || !h.requireDatabase(w, r) {
		return
	}

	leaderboardType := normalizeType(r.URL.Query().Get("type"))
	if !isValidType(leaderboardType) {
		httpx.WriteProblem(w, r, httpx.NewError(http.StatusUnprocessableEntity, "unsupported leaderboard type"))
		return
	}

	response, err := h.leaderboard(r.Context(), leaderboardType, player.TeamID)
	if err != nil {
		httpx.WriteProblem(w, r, httpx.InternalServerError("leaderboard unavailable", "leaderboard_lookup_failed", err))
		return
	}
	httpx.WriteJSON(w, http.StatusOK, response)
}

func normalizeType(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return TypeOpenPower
	}
	return value
}

func isValidType(value string) bool {
	switch value {
	case TypeOpenPower, TypeSitones, TypeMatches:
		return true
	default:
		return false
	}
}

func (h *Handler) leaderboard(ctx context.Context, leaderboardType, currentTeamID string) (ListResponse, error) {
	teams, err := h.findTeams(ctx)
	if err != nil {
		return ListResponse{}, err
	}
	scores, err := h.scoresByTeam(ctx, leaderboardType)
	if err != nil {
		return ListResponse{}, err
	}

	rows := make([]TeamRankResponse, 0, len(teams))
	for _, team := range teams {
		rows = append(rows, TeamRankResponse{
			TeamID:  team.ID,
			Name:    team.Name,
			Score:   scores[team.ID],
			Metric:  metricForType(leaderboardType),
			Current: team.ID == currentTeamID,
		})
	}
	sort.Slice(rows, func(i, j int) bool {
		if rows[i].Score == rows[j].Score {
			return rows[i].Name < rows[j].Name
		}
		return rows[i].Score > rows[j].Score
	})

	var currentTeam *TeamRankResponse
	gapToPrevious := 0
	for i := range rows {
		rows[i].Rank = i + 1
		if rows[i].Current {
			if i > 0 {
				gapToPrevious = rows[i-1].Score - rows[i].Score
			}
			currentTeam = &rows[i]
		}
	}

	return ListResponse{
		Type:          leaderboardType,
		Teams:         rows,
		CurrentTeam:   currentTeam,
		GapToPrevious: gapToPrevious,
	}, nil
}

func (h *Handler) findTeams(ctx context.Context) ([]mongomodel.Team, error) {
	cursor, err := h.db.Collection(mongomodel.TeamsCollection).Find(
		ctx,
		bson.M{},
		options.Find().SetSort(bson.D{{Key: "name", Value: 1}}),
	)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = cursor.Close(ctx)
	}()

	var teams []mongomodel.Team
	if err := cursor.All(ctx, &teams); err != nil {
		return nil, err
	}
	return teams, nil
}

func (h *Handler) scoresByTeam(ctx context.Context, leaderboardType string) (map[string]int, error) {
	var collection string
	var pipeline mongo.Pipeline
	switch leaderboardType {
	case TypeOpenPower:
		collection = mongomodel.OpenPowerRecordsCollection
		pipeline = openPowerScoresByTeamPipeline()
	case TypeSitones:
		collection = mongomodel.PlayerSitonesCollection
		pipeline = inventoryScoresByTeamPipeline()
	case TypeMatches:
		collection = mongomodel.MatchesCollection
		pipeline = matchScoresByTeamPipeline()
	default:
		return map[string]int{}, nil
	}

	cursor, err := h.db.Collection(collection).Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = cursor.Close(ctx)
	}()

	return scoreMapFromCursor(ctx, cursor)
}

func openPowerScoresByTeamPipeline() mongo.Pipeline {
	return mongo.Pipeline{
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$player_id"},
			{Key: "score", Value: bson.D{{Key: "$sum", Value: "$amount"}}},
		}}},
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: mongomodel.PlayersCollection},
			{Key: "localField", Value: "_id"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "player"},
		}}},
		{{Key: "$unwind", Value: "$player"}},
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$player.team_id"},
			{Key: "score", Value: bson.D{{Key: "$sum", Value: "$score"}}},
		}}},
	}
}

func inventoryScoresByTeamPipeline() mongo.Pipeline {
	return mongo.Pipeline{
		{{Key: "$match", Value: bson.D{{Key: "quantity", Value: bson.D{{Key: "$gt", Value: 0}}}}}},
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$player_id"},
			{Key: "score", Value: bson.D{{Key: "$sum", Value: "$quantity"}}},
		}}},
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: mongomodel.PlayersCollection},
			{Key: "localField", Value: "_id"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "player"},
		}}},
		{{Key: "$unwind", Value: "$player"}},
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$player.team_id"},
			{Key: "score", Value: bson.D{{Key: "$sum", Value: "$score"}}},
		}}},
	}
}

func matchScoresByTeamPipeline() mongo.Pipeline {
	return mongo.Pipeline{
		{{Key: "$match", Value: bson.D{{Key: "status", Value: mongomodel.MatchStatusCompleted}}}},
		{{Key: "$unwind", Value: "$players"}},
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$players.player_id"},
			{Key: "score", Value: bson.D{{Key: "$sum", Value: "$players.score"}}},
		}}},
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: mongomodel.PlayersCollection},
			{Key: "localField", Value: "_id"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "player"},
		}}},
		{{Key: "$unwind", Value: "$player"}},
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$player.team_id"},
			{Key: "score", Value: bson.D{{Key: "$sum", Value: "$score"}}},
		}}},
	}
}

func scoreMapFromCursor(ctx context.Context, cursor *mongo.Cursor) (map[string]int, error) {
	var rows []struct {
		ID    string `bson:"_id"`
		Score int    `bson:"score"`
	}
	if err := cursor.All(ctx, &rows); err != nil {
		return nil, err
	}

	out := make(map[string]int, len(rows))
	for _, row := range rows {
		out[row.ID] = row.Score
	}
	return out, nil
}

func metricForType(leaderboardType string) string {
	switch leaderboardType {
	case TypeSitones:
		return "小石"
	case TypeMatches:
		return "分"
	default:
		return "OP"
	}
}
