package matches

import (
	"context"
	"fmt"
	"net/http"
	"sort"
	"strings"

	"go.mongodb.org/mongo-driver/v2/bson"

	"github.com/sitcon-tw/camp2026-game/internal/http/httpx"
	mongomodel "github.com/sitcon-tw/camp2026-game/internal/mongodb/model"
)

const (
	minSitoneLoadoutSize = 1
	maxSitoneLoadoutSize = 5
)

func normalizeSitoneLoadout(sitoneIDs []string) ([]string, error) {
	normalized := make([]string, 0, len(sitoneIDs))

	for _, sitoneID := range sitoneIDs {
		sitoneID = strings.TrimSpace(sitoneID)
		if sitoneID == "" {
			continue
		}
		normalized = append(normalized, sitoneID)
	}

	if len(normalized) < minSitoneLoadoutSize {
		return nil, httpx.NewError(http.StatusUnprocessableEntity, "select at least one sitone")
	}
	if len(normalized) > maxSitoneLoadoutSize {
		return nil, httpx.NewError(http.StatusUnprocessableEntity, fmt.Sprintf("select at most %d sitones", maxSitoneLoadoutSize))
	}
	return normalized, nil
}

func (h *Handler) validateOwnedSitoneLoadout(ctx context.Context, playerID string, sitoneIDs []string) ([]string, error) {
	normalized, err := normalizeSitoneLoadout(sitoneIDs)
	if err != nil {
		return nil, err
	}

	owned, err := h.ownedSitoneCounts(ctx, playerID)
	if err != nil {
		return nil, httpx.InternalServerError("sitone loadout unavailable", "match_sitone_loadout_inventory_lookup_failed", err)
	}

	used := make(map[string]int, len(normalized))
	for _, sitoneID := range normalized {
		if _, ok := h.content.GetSitone(sitoneID); !ok {
			return nil, httpx.NewError(http.StatusUnprocessableEntity, "sitone loadout contains unknown sitone")
		}
		quantity := owned[sitoneID]
		if quantity <= 0 {
			return nil, httpx.NewError(http.StatusUnprocessableEntity, "sitone loadout contains unavailable sitone")
		}
		used[sitoneID]++
		if used[sitoneID] > quantity {
			return nil, httpx.NewError(http.StatusUnprocessableEntity, "sitone loadout exceeds owned quantity")
		}
	}

	return normalized, nil
}

func (h *Handler) defaultSitoneLoadout(ctx context.Context, player mongomodel.Player) ([]string, error) {
	if len(player.DefaultSitoneIDs) > 0 {
		if loadout, err := h.validateOwnedSitoneLoadout(ctx, player.ID, player.DefaultSitoneIDs); err == nil {
			return loadout, nil
		}
	}

	owned, err := h.ownedSitoneCounts(ctx, player.ID)
	if err != nil {
		return nil, err
	}

	if len(owned) == 0 {
		return nil, nil
	}

	ids := make([]string, 0, len(owned))
	for sitoneID := range owned {
		if _, ok := h.content.GetSitone(sitoneID); ok {
			ids = append(ids, sitoneID)
		}
	}
	sort.Strings(ids)

	loadout := make([]string, 0, maxSitoneLoadoutSize)
	for _, sitoneID := range ids {
		for i := 0; i < owned[sitoneID]; i++ {
			if len(loadout) >= maxSitoneLoadoutSize {
				return loadout, nil
			}
			loadout = append(loadout, sitoneID)
		}
	}
	return loadout, nil
}

func (h *Handler) ownedSitoneCounts(ctx context.Context, playerID string) (map[string]int, error) {
	cursor, err := h.db.Collection(mongomodel.PlayerSitonesCollection).Find(
		ctx,
		bson.M{"player_id": playerID, "quantity": bson.M{"$gt": 0}},
	)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = cursor.Close(ctx)
	}()

	owned := make(map[string]int)
	for cursor.Next(ctx) {
		var record mongomodel.PlayerSitone
		if err := cursor.Decode(&record); err != nil {
			return nil, err
		}
		if record.SitoneID != "" {
			owned[record.SitoneID] += record.Quantity
		}
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return owned, nil
}

func (h *Handler) saveDefaultSitoneLoadout(ctx context.Context, playerID string, sitoneIDs []string) error {
	_, err := h.db.Collection(mongomodel.PlayersCollection).UpdateOne(
		ctx,
		bson.M{"_id": playerID},
		bson.M{"$set": bson.M{"default_sitone_ids": sitoneIDs}},
	)
	return err
}

func cloneStrings(values []string) []string {
	if len(values) == 0 {
		return nil
	}
	out := make([]string, len(values))
	copy(out, values)
	return out
}
