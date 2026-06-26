package matches

import (
	"context"
	"sort"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"github.com/sitcon-tw/camp2026-game/internal/content"
	mongomodel "github.com/sitcon-tw/camp2026-game/internal/mongodb/model"
)

const (
	winnerMaterialDropBaseRate = 45
	loserMaterialDropBaseRate  = 25
	matchItemDropQuantity      = 1
)

func (h *Handler) writeMatchItemDrop(ctx context.Context, match mongomodel.Match, player mongomodel.MatchPlayer, effects battleEffects, createdAt time.Time) error {
	record := h.matchItemDrop(match, player, effects, createdAt)
	_, err := h.db.Collection(mongomodel.MatchItemDropsCollection).UpdateOne(
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
	if !stored.Dropped || stored.ItemID == "" || stored.Granted {
		return nil
	}
	quantity := stored.Quantity
	if quantity <= 0 {
		quantity = matchItemDropQuantity
	}

	if err := h.incrementPlayerItemForDrop(ctx, player.PlayerID, stored, quantity); err != nil {
		return err
	}
	_, err = h.db.Collection(mongomodel.MatchItemDropsCollection).UpdateOne(
		ctx,
		bson.M{"_id": stored.ID},
		bson.M{"$set": bson.M{"granted": true}},
	)
	return err
}

func (h *Handler) incrementPlayerItemForDrop(ctx context.Context, playerID string, drop mongomodel.MatchItemDrop, quantity int) error {
	if updated, err := h.incrementExistingPlayerItemForDrop(ctx, playerID, drop, quantity); err != nil || updated {
		return err
	}

	_, err := h.db.Collection(mongomodel.PlayerItemsCollection).InsertOne(ctx, bson.M{
		"_id":                dropPlayerItemRecordID(playerID, drop.ItemID),
		"player_id":          playerID,
		"item_id":            drop.ItemID,
		"quantity":           quantity,
		"drop_grant_sources": []string{drop.Source},
	})
	if err == nil {
		return nil
	}
	if !mongo.IsDuplicateKeyError(err) {
		return err
	}

	_, err = h.incrementExistingPlayerItemForDrop(ctx, playerID, drop, quantity)
	return err
}

func (h *Handler) incrementExistingPlayerItemForDrop(ctx context.Context, playerID string, drop mongomodel.MatchItemDrop, quantity int) (bool, error) {
	result, err := h.db.Collection(mongomodel.PlayerItemsCollection).UpdateOne(
		ctx,
		bson.M{
			"player_id":          playerID,
			"item_id":            drop.ItemID,
			"drop_grant_sources": bson.M{"$ne": drop.Source},
		},
		bson.M{
			"$set": bson.M{
				"player_id": playerID,
				"item_id":   drop.ItemID,
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

func (h *Handler) matchItemDrop(match mongomodel.Match, player mongomodel.MatchPlayer, effects battleEffects, createdAt time.Time) mongomodel.MatchItemDrop {
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
	if deterministicPercent("material-drop", match.ID, player.PlayerID) >= dropRate {
		return record
	}

	pool := h.materialDropPool()
	if len(pool) == 0 {
		return record
	}
	index := int(deterministicUint64("material-drop-item", match.ID, player.PlayerID) % uint64(len(pool)))
	record.ItemID = pool[index].ID
	record.Quantity = matchItemDropQuantity
	record.Dropped = true
	return record
}

func matchMaterialDropRate(match mongomodel.Match, player mongomodel.MatchPlayer, effects battleEffects) int {
	base := loserMaterialDropBaseRate
	if winner, ok := clearMatchWinner(match); ok && winner.PlayerID == player.PlayerID {
		base = winnerMaterialDropBaseRate
	}
	return base + effects.MaterialDropBonusPercent
}

func (h *Handler) materialDropPool() []content.Item {
	items := h.content.ListItems()
	out := make([]content.Item, 0, len(items))
	for _, item := range items {
		if item.Type != "material" || !item.Enabled {
			continue
		}
		if item.Source == "drop" || item.Source == "both" {
			out = append(out, item)
		}
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].ID < out[j].ID
	})
	return out
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
		Quantity: drop.Quantity,
		DropRate: drop.DropRate,
	}
	if drop.ItemID != "" {
		if item, ok := h.content.GetItem(drop.ItemID); ok {
			response.ItemName = item.Name
		}
	}
	return response
}

func matchItemDropRecordID(matchID, playerID string) string {
	return "match_item_drop_" + matchID + "_" + playerID
}

func matchItemDropSource(matchID, playerID string) string {
	return "quiz_match_item_drop:" + matchID + ":player:" + playerID
}

func dropPlayerItemRecordID(playerID, itemID string) string {
	return "player_item_drop_" + playerID + "_" + itemID
}
