package main

import (
	"context"
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"github.com/sitcon-tw/camp2026-game/internal/config"
	"github.com/sitcon-tw/camp2026-game/internal/mongodb"
	mongomodel "github.com/sitcon-tw/camp2026-game/internal/mongodb/model"
)

const (
	defaultTeamSize = 6
	playerRoleStaff = "staff"

	headerGroup      = "組別"
	headerNickname   = "暱稱"
	headerEmailHash  = "email (sha256)"
	headerToken      = "Token"
	headerStaffToken = "Staff Token"
)

type csvRow struct {
	RowNumber  int
	Group      string
	Nickname   string
	EmailHash  string
	Token      string
	StaffToken string
}

type importPlan struct {
	Rows    []csvRow
	Teams   []mongomodel.Team
	Players []mongomodel.Player
}

func main() {
	if err := run(context.Background(), os.Args[1:]); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "import staff failed: %v\n", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, args []string) error {
	flags := flag.NewFlagSet("import-staff", flag.ContinueOnError)
	flags.SetOutput(os.Stderr)

	csvPath := flags.String("csv", "", "path to staff CSV")
	teamSize := flags.Int("team-size", defaultTeamSize, "number of CSV rows assigned to each team")
	dryRun := flags.Bool("dry-run", false, "parse and validate CSV without writing to MongoDB")
	if err := flags.Parse(args); err != nil {
		return err
	}

	if strings.TrimSpace(*csvPath) == "" {
		return errors.New("-csv is required")
	}
	if *teamSize <= 0 {
		return errors.New("-team-size must be positive")
	}

	file, err := os.Open(*csvPath)
	if err != nil {
		return fmt.Errorf("open csv: %w", err)
	}
	defer func() {
		_ = file.Close()
	}()

	plan, err := buildImportPlan(file, *teamSize)
	if err != nil {
		return err
	}

	if *dryRun {
		printSummary(plan, true)
		return nil
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
	if err := upsertTeams(ctx, db.Collection(mongomodel.TeamsCollection), plan.Teams); err != nil {
		return err
	}
	if err := upsertPlayers(ctx, db.Collection(mongomodel.PlayersCollection), plan.Players); err != nil {
		return err
	}

	printSummary(plan, false)
	return nil
}

func buildImportPlan(reader io.Reader, teamSize int) (importPlan, error) {
	if teamSize <= 0 {
		return importPlan{}, errors.New("team size must be positive")
	}

	rows, err := readCSVRows(reader)
	if err != nil {
		return importPlan{}, err
	}

	teamsByID := make(map[string]mongomodel.Team)
	teams := make([]mongomodel.Team, 0, (len(rows)+teamSize-1)/teamSize)
	players := make([]mongomodel.Player, 0, len(rows)*2)
	seenTokens := make(map[string]int, len(rows)*2)

	for index, row := range rows {
		teamNumber := index/teamSize + 1
		team := mongomodel.Team{
			ID:   teamID(teamNumber),
			Name: teamName(teamNumber),
		}
		if _, ok := teamsByID[team.ID]; !ok {
			teamsByID[team.ID] = team
			teams = append(teams, team)
		}

		avatarURL := gravatarURL(row.EmailHash)
		regular := mongomodel.Player{
			ID:          row.Token,
			AuthToken:   row.Token,
			QRCodeToken: row.Token,
			Nickname:    row.Nickname,
			TeamID:      team.ID,
			AvatarURL:   avatarURL,
		}
		staff := mongomodel.Player{
			ID:        row.StaffToken,
			AuthToken: row.StaffToken,
			Nickname:  row.Nickname,
			TeamID:    team.ID,
			AvatarURL: avatarURL,
			Role:      playerRoleStaff,
		}

		for _, player := range []mongomodel.Player{regular, staff} {
			if previousRow, ok := seenTokens[player.ID]; ok {
				return importPlan{}, fmt.Errorf("duplicate token %q on CSV rows %d and %d", player.ID, previousRow, row.RowNumber)
			}
			seenTokens[player.ID] = row.RowNumber
			players = append(players, player)
		}
	}

	return importPlan{
		Rows:    rows,
		Teams:   teams,
		Players: players,
	}, nil
}

func readCSVRows(reader io.Reader) ([]csvRow, error) {
	csvReader := csv.NewReader(reader)
	csvReader.FieldsPerRecord = -1
	csvReader.TrimLeadingSpace = true

	header, err := csvReader.Read()
	if errors.Is(err, io.EOF) {
		return nil, errors.New("csv is empty")
	}
	if err != nil {
		return nil, fmt.Errorf("read csv header: %w", err)
	}

	columns, err := requiredColumns(header)
	if err != nil {
		return nil, err
	}

	var rows []csvRow
	for {
		record, err := csvReader.Read()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("read csv row: %w", err)
		}
		rowNumber := len(rows) + 2
		row, err := parseRow(record, columns, rowNumber)
		if err != nil {
			return nil, err
		}
		rows = append(rows, row)
	}

	if len(rows) == 0 {
		return nil, errors.New("csv contains no data rows")
	}
	return rows, nil
}

func requiredColumns(header []string) (map[string]int, error) {
	indexes := make(map[string]int, len(header))
	for index, value := range header {
		name := strings.TrimPrefix(strings.TrimSpace(value), "\ufeff")
		indexes[name] = index
	}

	for _, required := range []string{headerGroup, headerNickname, headerEmailHash, headerToken, headerStaffToken} {
		if _, ok := indexes[required]; !ok {
			return nil, fmt.Errorf("csv header is missing %q", required)
		}
	}
	return indexes, nil
}

func parseRow(record []string, columns map[string]int, rowNumber int) (csvRow, error) {
	value := func(header string) string {
		index := columns[header]
		if index >= len(record) {
			return ""
		}
		return strings.TrimSpace(record[index])
	}

	row := csvRow{
		RowNumber:  rowNumber,
		Group:      value(headerGroup),
		Nickname:   value(headerNickname),
		EmailHash:  strings.ToLower(value(headerEmailHash)),
		Token:      value(headerToken),
		StaffToken: value(headerStaffToken),
	}
	if row.Group == "" {
		return csvRow{}, fmt.Errorf("csv row %d is missing %q", rowNumber, headerGroup)
	}
	if row.Nickname == "" {
		return csvRow{}, fmt.Errorf("csv row %d is missing %q", rowNumber, headerNickname)
	}
	if row.EmailHash == "" {
		return csvRow{}, fmt.Errorf("csv row %d is missing %q", rowNumber, headerEmailHash)
	}
	if row.Token == "" {
		return csvRow{}, fmt.Errorf("csv row %d is missing %q", rowNumber, headerToken)
	}
	if row.StaffToken == "" {
		return csvRow{}, fmt.Errorf("csv row %d is missing %q", rowNumber, headerStaffToken)
	}
	if !isHexSHA256(row.EmailHash) {
		return csvRow{}, fmt.Errorf("csv row %d has invalid %q", rowNumber, headerEmailHash)
	}
	return row, nil
}

func isHexSHA256(value string) bool {
	if len(value) != 64 {
		return false
	}
	for _, char := range value {
		if !((char >= '0' && char <= '9') || (char >= 'a' && char <= 'f')) {
			return false
		}
	}
	return true
}

func teamID(number int) string {
	return fmt.Sprintf("team-%03d", number)
}

func teamName(number int) string {
	return fmt.Sprintf("Team %03d", number)
}

func gravatarURL(hash string) string {
	return "https://www.gravatar.com/avatar/" + hash
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

func upsertPlayers(ctx context.Context, collection *mongo.Collection, players []mongomodel.Player) error {
	for _, player := range players {
		set := bson.M{
			"auth_token": player.AuthToken,
			"nickname":   player.Nickname,
			"team_id":    player.TeamID,
			"avatar_url": player.AvatarURL,
		}
		unset := bson.M{}
		if player.QRCodeToken == "" {
			unset["qrcode_token"] = ""
		} else {
			set["qrcode_token"] = player.QRCodeToken
		}
		if player.Role == "" {
			unset["role"] = ""
		} else {
			set["role"] = player.Role
		}

		update := bson.M{"$set": set}
		if len(unset) > 0 {
			update["$unset"] = unset
		}

		_, err := collection.UpdateOne(
			ctx,
			bson.M{"_id": player.ID},
			update,
			options.UpdateOne().SetUpsert(true),
		)
		if err != nil {
			return fmt.Errorf("upsert %s/%s: %w", collection.Name(), player.ID, err)
		}
	}
	return nil
}

func printSummary(plan importPlan, dryRun bool) {
	prefix := "import complete"
	if dryRun {
		prefix = "dry run complete"
	}
	fmt.Printf("%s: rows=%d teams=%d players=%d staff=%d\n",
		prefix,
		len(plan.Rows),
		len(plan.Teams),
		len(plan.Rows),
		len(plan.Rows),
	)
}
