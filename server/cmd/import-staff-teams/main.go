package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"github.com/sitcon-tw/camp2026-game/internal/config"
	"github.com/sitcon-tw/camp2026-game/internal/mongodb"
	mongomodel "github.com/sitcon-tw/camp2026-game/internal/mongodb/model"
)

const (
	playerRoleStaff = "staff"
	expectedTeams   = 10
)

type assignmentFile struct {
	Teams []staffTeamAssignment `json:"teams"`
}

type staffTeamAssignment struct {
	TeamID string   `json:"teamId"`
	Name   string   `json:"name"`
	Staff  []string `json:"staff"`
}

type importPlan struct {
	Teams   []mongomodel.Team
	Updates []staffTeamUpdate
}

type staffTeamUpdate struct {
	PlayerID string
	Nickname string
	TeamID   string
	TeamName string
}

func main() {
	if err := run(context.Background(), os.Args[1:]); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "import staff teams failed: %v\n", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, args []string) error {
	flags := flag.NewFlagSet("import-staff-teams", flag.ContinueOnError)
	flags.SetOutput(os.Stderr)

	jsonPath := flags.String("json", "", "path to staff team assignment JSON")
	dryRun := flags.Bool("dry-run", false, "validate and print planned changes without writing to MongoDB")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if strings.TrimSpace(*jsonPath) == "" {
		return errors.New("-json is required")
	}

	file, err := os.Open(*jsonPath)
	if err != nil {
		return fmt.Errorf("open json: %w", err)
	}
	defer func() {
		_ = file.Close()
	}()

	assignments, err := readAssignments(file)
	if err != nil {
		return err
	}

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	client, err := mongodb.NewClient(ctx, cfg.MongoURI)
	if err != nil {
		return err
	}
	defer func() {
		_ = client.Disconnect(context.Background())
	}()

	db := client.Database(cfg.MongoDatabase)
	players, err := findStaffPlayers(ctx, db.Collection(mongomodel.PlayersCollection), staffNames(assignments))
	if err != nil {
		return err
	}

	plan, err := buildImportPlan(assignments, players)
	if err != nil {
		return err
	}

	if *dryRun {
		printSummary(plan, true)
		return nil
	}

	if err := upsertTeams(ctx, db.Collection(mongomodel.TeamsCollection), plan.Teams); err != nil {
		return err
	}
	if err := updateStaffTeams(ctx, db.Collection(mongomodel.PlayersCollection), plan.Updates); err != nil {
		return err
	}

	printSummary(plan, false)
	return nil
}

func readAssignments(reader io.Reader) (assignmentFile, error) {
	var assignments assignmentFile
	decoder := json.NewDecoder(reader)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&assignments); err != nil {
		return assignmentFile{}, fmt.Errorf("decode assignment json: %w", err)
	}
	assignments = normalizeAssignments(assignments)
	if err := validateAssignments(assignments); err != nil {
		return assignmentFile{}, err
	}
	return assignments, nil
}

func normalizeAssignments(assignments assignmentFile) assignmentFile {
	out := assignmentFile{Teams: make([]staffTeamAssignment, 0, len(assignments.Teams))}
	for _, team := range assignments.Teams {
		normalized := staffTeamAssignment{
			TeamID: strings.TrimSpace(team.TeamID),
			Name:   strings.TrimSpace(team.Name),
			Staff:  make([]string, 0, len(team.Staff)),
		}
		for _, staff := range team.Staff {
			normalized.Staff = append(normalized.Staff, strings.TrimSpace(staff))
		}
		out.Teams = append(out.Teams, normalized)
	}
	return out
}

func validateAssignments(assignments assignmentFile) error {
	if len(assignments.Teams) != expectedTeams {
		return fmt.Errorf("assignment json must contain %d teams, got %d", expectedTeams, len(assignments.Teams))
	}

	seenStaff := make(map[string]string)
	for index, team := range assignments.Teams {
		number := index + 1
		if team.TeamID != teamID(number) {
			return fmt.Errorf("team %d has id %q, expected %q", number, team.TeamID, teamID(number))
		}
		if team.Name != teamName(number) {
			return fmt.Errorf("team %s has name %q, expected %q", team.TeamID, team.Name, teamName(number))
		}
		if len(team.Staff) == 0 {
			return fmt.Errorf("team %s must contain at least one staff", team.TeamID)
		}
		for _, name := range team.Staff {
			if name == "" {
				return fmt.Errorf("team %s contains an empty staff nickname", team.TeamID)
			}
			if previousTeamID, ok := seenStaff[name]; ok {
				return fmt.Errorf("staff %q appears in both %s and %s", name, previousTeamID, team.TeamID)
			}
			seenStaff[name] = team.TeamID
		}
	}
	return nil
}

func buildImportPlan(assignments assignmentFile, players []mongomodel.Player) (importPlan, error) {
	if err := validateAssignments(assignments); err != nil {
		return importPlan{}, err
	}

	requiredStaff := staffNameSet(assignments)
	playersByNickname := make(map[string][]mongomodel.Player, len(requiredStaff))
	for _, player := range players {
		if _, ok := requiredStaff[player.Nickname]; !ok {
			continue
		}
		if player.Role != playerRoleStaff {
			continue
		}
		if player.ID == "" {
			return importPlan{}, fmt.Errorf("staff %q has empty player id", player.Nickname)
		}
		playersByNickname[player.Nickname] = append(playersByNickname[player.Nickname], player)
	}

	plan := importPlan{
		Teams:   make([]mongomodel.Team, 0, len(assignments.Teams)),
		Updates: make([]staffTeamUpdate, 0, len(requiredStaff)),
	}
	for _, team := range assignments.Teams {
		plan.Teams = append(plan.Teams, mongomodel.Team{ID: team.TeamID, Name: team.Name})
		for _, nickname := range team.Staff {
			matches := playersByNickname[nickname]
			switch len(matches) {
			case 0:
				return importPlan{}, fmt.Errorf("missing existing staff player with nickname %q", nickname)
			case 1:
				plan.Updates = append(plan.Updates, staffTeamUpdate{
					PlayerID: matches[0].ID,
					Nickname: nickname,
					TeamID:   team.TeamID,
					TeamName: team.Name,
				})
			default:
				return importPlan{}, fmt.Errorf("multiple existing staff players with nickname %q", nickname)
			}
		}
	}

	sort.Slice(plan.Updates, func(i, j int) bool {
		if plan.Updates[i].TeamID != plan.Updates[j].TeamID {
			return plan.Updates[i].TeamID < plan.Updates[j].TeamID
		}
		return plan.Updates[i].Nickname < plan.Updates[j].Nickname
	})
	return plan, nil
}

func staffNames(assignments assignmentFile) []string {
	names := make([]string, 0)
	for _, team := range assignments.Teams {
		names = append(names, team.Staff...)
	}
	sort.Strings(names)
	return names
}

func staffNameSet(assignments assignmentFile) map[string]struct{} {
	names := make(map[string]struct{})
	for _, name := range staffNames(assignments) {
		names[name] = struct{}{}
	}
	return names
}

func findStaffPlayers(ctx context.Context, collection *mongo.Collection, nicknames []string) ([]mongomodel.Player, error) {
	cursor, err := collection.Find(
		ctx,
		bson.M{
			"nickname": bson.M{"$in": nicknames},
			"role":     playerRoleStaff,
		},
		options.Find().
			SetProjection(bson.D{
				{Key: "auth_token", Value: 0},
				{Key: "qrcode_token", Value: 0},
				{Key: "default_sitone_ids", Value: 0},
			}).
			SetSort(bson.D{
				{Key: "nickname", Value: 1},
				{Key: "_id", Value: 1},
			}),
	)
	if err != nil {
		return nil, fmt.Errorf("find staff players: %w", err)
	}
	defer func() {
		_ = cursor.Close(ctx)
	}()

	var players []mongomodel.Player
	if err := cursor.All(ctx, &players); err != nil {
		return nil, fmt.Errorf("decode staff players: %w", err)
	}
	if players == nil {
		return []mongomodel.Player{}, nil
	}
	return players, nil
}

func upsertTeams(ctx context.Context, collection *mongo.Collection, teams []mongomodel.Team) error {
	for _, team := range teams {
		_, err := collection.UpdateOne(
			ctx,
			bson.M{"_id": team.ID},
			bson.M{"$set": bson.M{"name": team.Name}},
			options.UpdateOne().SetUpsert(true),
		)
		if err != nil {
			return fmt.Errorf("upsert %s/%s: %w", collection.Name(), team.ID, err)
		}
	}
	return nil
}

func updateStaffTeams(ctx context.Context, collection *mongo.Collection, updates []staffTeamUpdate) error {
	for _, update := range updates {
		result, err := collection.UpdateOne(
			ctx,
			bson.M{
				"_id":  update.PlayerID,
				"role": playerRoleStaff,
			},
			bson.M{"$set": bson.M{
				"role":    playerRoleStaff,
				"team_id": update.TeamID,
			}},
		)
		if err != nil {
			return fmt.Errorf("update %s staff %q: %w", collection.Name(), update.Nickname, err)
		}
		if result.MatchedCount != 1 {
			return fmt.Errorf("update %s staff %q: matched %d staff players", collection.Name(), update.Nickname, result.MatchedCount)
		}
	}
	return nil
}

func teamID(number int) string {
	return fmt.Sprintf("team-%03d", number)
}

func teamName(number int) string {
	return fmt.Sprintf("Team %03d", number)
}

func printSummary(plan importPlan, dryRun bool) {
	prefix := "import complete"
	if dryRun {
		prefix = "dry run complete"
	}
	fmt.Printf("%s: teams=%d staff_updates=%d\n", prefix, len(plan.Teams), len(plan.Updates))
	for _, update := range plan.Updates {
		fmt.Printf("%s -> %s\n", update.Nickname, update.TeamID)
	}
}
