package matches

import (
	"context"
	"crypto/rand"
	"fmt"
	"io"
	"math/big"
	"sort"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"github.com/sitcon-tw/camp2026-game/internal/content"
	mongomodel "github.com/sitcon-tw/camp2026-game/internal/mongodb/model"
)

const (
	winnerMaterialDropBaseRate = 25
	loserMaterialDropBaseRate  = 15
	matchItemDropQuantity      = 1
)

var baseSitoneDropIDs = []string{
	"stone_explorer_base",
	"stone_inspiration_base",
	"stone_resonance_base",
	"stone_engineering_base",
	"stone_entertainment_base",
}

func (h *Handler) writeMatchItemDrop(ctx context.Context, match mongomodel.Match, player mongomodel.MatchPlayer, effects battleEffects, createdAt time.Time) error {
	record, err := h.matchItemDrop(ctx, match, player, effects, createdAt)
	if err != nil {
		return err
	}
	_, err = h.db.Collection(mongomodel.MatchItemDropsCollection).UpdateOne(
		ctx,
		bson.M{"_id": record.ID},
		bson.M{"$setOnInsert": record},
		options.UpdateOne().SetUpsert(true),
	)
	if err != nil {
		return err
	}

	var stored mongomodel.MatchItemDrop
	if err := h.db.Collection(mongomodel.MatchItemDropsCollection).
		FindOne(ctx, bson.M{"_id": record.ID}).
		Decode(&stored); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil
		}
		return err
	}
	if !stored.Dropped || stored.SitoneID == "" || stored.Granted {
		return nil
	}
	quantity := stored.Quantity
	if quantity <= 0 {
		quantity = matchItemDropQuantity
	}

	if err := h.incrementPlayerSitoneForDrop(ctx, player.PlayerID, stored, quantity); err != nil {
		return err
	}
	_, err = h.db.Collection(mongomodel.MatchItemDropsCollection).UpdateOne(
		ctx,
		bson.M{"_id": stored.ID},
		bson.M{"$set": bson.M{"granted": true}},
	)
	return err
}

func (h *Handler) incrementPlayerSitoneForDrop(ctx context.Context, playerID string, drop mongomodel.MatchItemDrop, quantity int) error {
	if updated, err := h.incrementExistingPlayerSitoneForDrop(ctx, playerID, drop, quantity); err != nil || updated {
		return err
	}

	_, err := h.db.Collection(mongomodel.PlayerSitonesCollection).InsertOne(ctx, bson.M{
		"_id":                dropPlayerSitoneRecordID(playerID, drop.SitoneID),
		"player_id":          playerID,
		"sitone_id":          drop.SitoneID,
		"quantity":           quantity,
		"drop_grant_sources": []string{drop.Source},
	})
	if err == nil {
		return nil
	}
	if !mongo.IsDuplicateKeyError(err) {
		return err
	}

	_, err = h.incrementExistingPlayerSitoneForDrop(ctx, playerID, drop, quantity)
	return err
}

func (h *Handler) incrementExistingPlayerSitoneForDrop(ctx context.Context, playerID string, drop mongomodel.MatchItemDrop, quantity int) (bool, error) {
	result, err := h.db.Collection(mongomodel.PlayerSitonesCollection).UpdateOne(
		ctx,
		bson.M{
			"player_id":          playerID,
			"sitone_id":          drop.SitoneID,
			"drop_grant_sources": bson.M{"$ne": drop.Source},
		},
		bson.M{
			"$set": bson.M{
				"player_id": playerID,
				"sitone_id": drop.SitoneID,
			},
			"$inc":      bson.M{"quantity": quantity},
			"$addToSet": bson.M{"drop_grant_sources": drop.Source},
		},
	)
	if err != nil {
		return false, err
	}
	return result.MatchedCount > 0, nil
}

func (h *Handler) matchItemDrop(ctx context.Context, match mongomodel.Match, player mongomodel.MatchPlayer, effects battleEffects, createdAt time.Time) (mongomodel.MatchItemDrop, error) {
	dropRate := matchMaterialDropRate(match, player, effects)
	record := mongomodel.MatchItemDrop{
		ID:        matchItemDropRecordID(match.ID, player.PlayerID),
		MatchID:   match.ID,
		PlayerID:  player.PlayerID,
		DropRate:  dropRate,
		Dropped:   false,
		Source:    matchItemDropSource(match.ID, player.PlayerID),
		CreatedAt: createdAt,
	}
	roll, err := secureRandomPercent()
	if err != nil {
		return mongomodel.MatchItemDrop{}, err
	}
	if roll >= dropRate {
		return record, nil
	}

	pool := h.baseSitoneDropPool()
	if len(pool) == 0 {
		return record, nil
	}
	pool, err = h.preferredBaseSitoneDropPool(ctx, player.PlayerID, pool)
	if err != nil {
		return mongomodel.MatchItemDrop{}, err
	}
	index, err := secureRandomIndex(len(pool))
	if err != nil {
		return mongomodel.MatchItemDrop{}, err
	}
	record.SitoneID = pool[index].ID
	record.Quantity = matchItemDropQuantity
	record.Dropped = true
	return record, nil
}

func matchMaterialDropRate(match mongomodel.Match, player mongomodel.MatchPlayer, effects battleEffects) int {
	base := loserMaterialDropBaseRate
	if winner, ok := clearMatchWinner(match); ok && winner.PlayerID == player.PlayerID {
		base = winnerMaterialDropBaseRate
	}
	return base + effects.MaterialDropBonusPercent
}

func (h *Handler) baseSitoneDropPool() []content.Sitone {
	out := make([]content.Sitone, 0, len(baseSitoneDropIDs))
	for _, sitoneID := range baseSitoneDropIDs {
		sitone, ok := h.content.GetSitone(sitoneID)
		if !ok {
			continue
		}
		out = append(out, sitone)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].ID < out[j].ID
	})
	return out
}

func (h *Handler) preferredBaseSitoneDropPool(ctx context.Context, playerID string, pool []content.Sitone) ([]content.Sitone, error) {
	if h.db == nil || len(pool) == 0 {
		return pool, nil
	}

	owned, err := h.ownedSitoneQuantities(ctx, playerID, sitoneDropPoolIDs(pool))
	if err != nil {
		return nil, err
	}
	return leastOwnedSitoneDropPool(pool, owned), nil
}

func (h *Handler) ownedSitoneQuantities(ctx context.Context, playerID string, sitoneIDs []string) (map[string]int, error) {
	cursor, err := h.db.Collection(mongomodel.PlayerSitonesCollection).Find(
		ctx,
		bson.M{
			"player_id": playerID,
			"sitone_id": bson.M{"$in": sitoneIDs},
			"quantity":  bson.M{"$gt": 0},
		},
	)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = cursor.Close(ctx)
	}()

	owned := make(map[string]int)
	for cursor.Next(ctx) {
		var sitone mongomodel.PlayerSitone
		if err := cursor.Decode(&sitone); err != nil {
			return nil, err
		}
		if sitone.SitoneID != "" {
			owned[sitone.SitoneID] += sitone.Quantity
		}
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return owned, nil
}

func sitoneDropPoolIDs(pool []content.Sitone) []string {
	ids := make([]string, 0, len(pool))
	for _, sitone := range pool {
		ids = append(ids, sitone.ID)
	}
	return ids
}

func leastOwnedSitoneDropPool(pool []content.Sitone, owned map[string]int) []content.Sitone {
	if len(pool) == 0 {
		return pool
	}

	minQuantity := owned[pool[0].ID]
	for _, sitone := range pool[1:] {
		if owned[sitone.ID] < minQuantity {
			minQuantity = owned[sitone.ID]
		}
	}

	out := make([]content.Sitone, 0, len(pool))
	for _, sitone := range pool {
		if owned[sitone.ID] != minQuantity {
			continue
		}
		out = append(out, sitone)
	}
	return out
}

func secureRandomPercent() (int, error) {
	return secureRandomInt(rand.Reader, 100)
}

func secureRandomIndex(length int) (int, error) {
	return secureRandomInt(rand.Reader, length)
}

func secureRandomInt(random io.Reader, max int) (int, error) {
	if max <= 0 {
		return 0, fmt.Errorf("random max must be positive")
	}
	value, err := rand.Int(random, big.NewInt(int64(max)))
	if err != nil {
		return 0, fmt.Errorf("secure random int: %w", err)
	}
	return int(value.Int64()), nil
}

func (h *Handler) findMatchItemDrops(ctx context.Context, matchID string) ([]mongomodel.MatchItemDrop, error) {
	cursor, err := h.db.Collection(mongomodel.MatchItemDropsCollection).Find(
		ctx,
		bson.M{"match_id": matchID},
		options.Find().SetSort(bson.D{{Key: "player_id", Value: 1}}),
	)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = cursor.Close(ctx)
	}()

	var drops []mongomodel.MatchItemDrop
	if err := cursor.All(ctx, &drops); err != nil {
		return nil, err
	}
	return drops, nil
}

func (h *Handler) matchItemDropResponse(drop mongomodel.MatchItemDrop) *MatchItemDropResponse {
	response := &MatchItemDropResponse{
		Dropped:  drop.Dropped,
		ItemID:   drop.ItemID,
		SitoneID: drop.SitoneID,
		Quantity: drop.Quantity,
		DropRate: drop.DropRate,
	}
	if drop.ItemID != "" {
		if item, ok := h.content.GetItem(drop.ItemID); ok {
			response.ItemName = item.Name
		}
	}
	if drop.SitoneID != "" {
		if sitone, ok := h.content.GetSitone(drop.SitoneID); ok {
			response.SitoneName = sitone.Name
		}
	}
	return response
}

func matchItemDropRecordID(matchID, playerID string) string {
	return "match_item_drop_" + matchID + "_" + playerID
}

func matchItemDropSource(matchID, playerID string) string {
	return "quiz_match_sitone_drop:" + matchID + ":player:" + playerID
}

func dropPlayerSitoneRecordID(playerID, sitoneID string) string {
	return "player_sitone_drop_" + playerID + "_" + sitoneID
}
