package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"github.com/sitcon-tw/camp2026-game/internal/config"
	"github.com/sitcon-tw/camp2026-game/internal/content"
	"github.com/sitcon-tw/camp2026-game/internal/mongodb"
	mongomodel "github.com/sitcon-tw/camp2026-game/internal/mongodb/model"
)

const (
	defaultSitoneLimit = 5
	playerRoleStaff    = "staff"
)

var defaultInitialSitoneIDs = []string{
	"stone_explorer_base",
	"stone_inspiration_base",
	"stone_resonance_base",
	"stone_engineering_base",
	"stone_entertainment_base",
}

type grant struct {
	Player     mongomodel.Player
	OwnedTotal int
	SitoneIDs  []string
}

type grantPlan struct {
	Players        []mongomodel.Player
	Grants         []grant
	SkippedAtLimit int
	SkippedNoGrant int
}

func main() {
	if err := run(context.Background(), os.Args[1:]); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "grant initial sitones failed: %v\n", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, args []string) error {
	flags := flag.NewFlagSet("grant-initial-sitones", flag.ContinueOnError)
	flags.SetOutput(os.Stderr)

	dryRun := flags.Bool("dry-run", false, "compute grants without writing to MongoDB")
	includeStaff := flags.Bool("include-staff", false, "include staff players")
	limit := flags.Int("limit", defaultSitoneLimit, "maximum total sitone quantity a player can have after the grant")
	sitoneIDsValue := flags.String("sitones", strings.Join(defaultInitialSitoneIDs, ","), "comma-separated sitone IDs to grant in order")
	verbose := flags.Bool("verbose", false, "print per-player grant details")
	if err := flags.Parse(args); err != nil {
		return err
	}

	sitoneIDs, err := parseSitoneIDs(*sitoneIDsValue)
	if err != nil {
		return err
	}
	if *limit <= 0 {
		return errors.New("-limit must be positive")
	}
	if len(sitoneIDs) > *limit {
		return fmt.Errorf("-sitones contains %d IDs but -limit is %d", len(sitoneIDs), *limit)
	}

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}
	store, err := content.Load(cfg.ContentDir)
	if err != nil {
		return fmt.Errorf("load content: %w", err)
	}
	if err := validateSitoneIDs(store, sitoneIDs); err != nil {
		return err
	}

	client, err := mongodb.NewClient(ctx, cfg.MongoURI)
	if err != nil {
		return err
	}
	defer func() {
		_ = client.Disconnect(context.Background())
	}()

	db := client.Database(cfg.MongoDatabase)
	players, err := findPlayers(ctx, db.Collection(mongomodel.PlayersCollection), *includeStaff)
	if err != nil {
		return err
	}
	owned, err := findOwnedSitones(ctx, db.Collection(mongomodel.PlayerSitonesCollection))
	if err != nil {
		return err
	}

	plan := buildGrantPlan(players, owned, sitoneIDs, *limit)
	if *verbose {
		printGrantDetails(plan)
	}
	if *dryRun {
		printSummary(plan, true)
		return nil
	}

	if err := applyGrants(ctx, db.Collection(mongomodel.PlayerSitonesCollection), plan.Grants); err != nil {
		return err
	}
	printSummary(plan, false)
	return nil
}

func parseSitoneIDs(value string) ([]string, error) {
	parts := strings.Split(value, ",")
	ids := make([]string, 0, len(parts))
	seen := make(map[string]struct{}, len(parts))
	for _, part := range parts {
		id := strings.TrimSpace(part)
		if id == "" {
			continue
		}
		if _, ok := seen[id]; ok {
			return nil, fmt.Errorf("duplicate sitone id %q", id)
		}
		seen[id] = struct{}{}
		ids = append(ids, id)
	}
	if len(ids) == 0 {
		return nil, errors.New("-sitones must include at least one sitone id")
	}
	return ids, nil
}

func validateSitoneIDs(store *content.Store, sitoneIDs []string) error {
	for _, sitoneID := range sitoneIDs {
		if _, ok := store.GetSitone(sitoneID); !ok {
			return fmt.Errorf("unknown sitone id %q", sitoneID)
		}
	}
	return nil
}

func findPlayers(ctx context.Context, collection *mongo.Collection, includeStaff bool) ([]mongomodel.Player, error) {
	filter := bson.M{}
	if !includeStaff {
		filter["role"] = bson.M{"$ne": playerRoleStaff}
	}

	cursor, err := collection.Find(ctx, filter, options.Find().SetSort(bson.D{{Key: "_id", Value: 1}}))
	if err != nil {
		return nil, fmt.Errorf("find players: %w", err)
	}
	defer func() {
		_ = cursor.Close(ctx)
	}()

	var players []mongomodel.Player
	if err := cursor.All(ctx, &players); err != nil {
		return nil, fmt.Errorf("decode players: %w", err)
	}
	return players, nil
}

func findOwnedSitones(ctx context.Context, collection *mongo.Collection) (map[string]map[string]int, error) {
	cursor, err := collection.Find(ctx, bson.M{"quantity": bson.M{"$gt": 0}})
	if err != nil {
		return nil, fmt.Errorf("find player sitones: %w", err)
	}
	defer func() {
		_ = cursor.Close(ctx)
	}()

	owned := make(map[string]map[string]int)
	for cursor.Next(ctx) {
		var record mongomodel.PlayerSitone
		if err := cursor.Decode(&record); err != nil {
			return nil, fmt.Errorf("decode player sitone: %w", err)
		}
		if record.PlayerID == "" || record.SitoneID == "" || record.Quantity <= 0 {
			continue
		}
		if owned[record.PlayerID] == nil {
			owned[record.PlayerID] = make(map[string]int)
		}
		owned[record.PlayerID][record.SitoneID] += record.Quantity
	}
	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("iterate player sitones: %w", err)
	}
	return owned, nil
}

func buildGrantPlan(
	players []mongomodel.Player,
	owned map[string]map[string]int,
	initialSitoneIDs []string,
	limit int,
) grantPlan {
	plan := grantPlan{Players: players}
	for _, player := range players {
		ownedBySitone := owned[player.ID]
		total := totalOwnedSitones(ownedBySitone)
		if total >= limit {
			plan.SkippedAtLimit++
			continue
		}

		remaining := limit - total
		sitoneIDs := make([]string, 0, remaining)
		for _, sitoneID := range initialSitoneIDs {
			if remaining == 0 {
				break
			}
			if ownedBySitone[sitoneID] > 0 {
				continue
			}
			sitoneIDs = append(sitoneIDs, sitoneID)
			remaining--
		}
		if len(sitoneIDs) == 0 {
			plan.SkippedNoGrant++
			continue
		}
		plan.Grants = append(plan.Grants, grant{
			Player:     player,
			OwnedTotal: total,
			SitoneIDs:  sitoneIDs,
		})
	}
	return plan
}

func totalOwnedSitones(owned map[string]int) int {
	total := 0
	for _, quantity := range owned {
		if quantity > 0 {
			total += quantity
		}
	}
	return total
}

func applyGrants(ctx context.Context, collection *mongo.Collection, grants []grant) error {
	for _, grant := range grants {
		for _, sitoneID := range grant.SitoneIDs {
			_, err := collection.UpdateOne(
				ctx,
				bson.M{
					"player_id": grant.Player.ID,
					"sitone_id": sitoneID,
				},
				bson.M{
					"$setOnInsert": bson.M{
						"_id":       newID("player_sitone"),
						"player_id": grant.Player.ID,
						"sitone_id": sitoneID,
					},
					"$inc": bson.M{"quantity": 1},
				},
				options.UpdateOne().SetUpsert(true),
			)
			if err != nil {
				return fmt.Errorf("grant %s to %s: %w", sitoneID, grant.Player.ID, err)
			}
		}
	}
	return nil
}

func newID(prefix string) string {
	return prefix + "_" + bson.NewObjectID().Hex()
}

func printGrantDetails(plan grantPlan) {
	grants := append([]grant(nil), plan.Grants...)
	sort.Slice(grants, func(i, j int) bool {
		return grants[i].Player.ID < grants[j].Player.ID
	})
	for _, grant := range grants {
		fmt.Printf("grant player=%s nickname=%q owned=%d sitones=%s\n",
			grant.Player.ID,
			grant.Player.Nickname,
			grant.OwnedTotal,
			strings.Join(grant.SitoneIDs, ","),
		)
	}
}

func printSummary(plan grantPlan, dryRun bool) {
	prefix := "grant complete"
	if dryRun {
		prefix = "dry run complete"
	}
	fmt.Printf("%s: players=%d grant_players=%d grant_sitones=%d skipped_at_limit=%d skipped_no_grant=%d\n",
		prefix,
		len(plan.Players),
		len(plan.Grants),
		countGrantedSitones(plan.Grants),
		plan.SkippedAtLimit,
		plan.SkippedNoGrant,
	)
}

func countGrantedSitones(grants []grant) int {
	total := 0
	for _, grant := range grants {
		total += len(grant.SitoneIDs)
	}
	return total
}
