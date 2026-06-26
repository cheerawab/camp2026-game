package me

import (
	"context"
	"net/http"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"github.com/sitcon-tw/camp2026-game/internal/content"
	"github.com/sitcon-tw/camp2026-game/internal/http/httpx"
	mongomodel "github.com/sitcon-tw/camp2026-game/internal/mongodb/model"
)

// ListItems godoc
// @Summary List current player items
// @Description Returns items owned by the authenticated player with catalog definitions.
// @Tags me
// @Produce json
// @Security AuthCookieAuth
// @Success 200 {object} ItemListResponse
// @Failure 401 {object} httpx.ProblemDetails
// @Failure 500 {object} httpx.ProblemDetails
// @Failure 503 {object} httpx.ProblemDetails
// @Router /me/items [get]
func (h *Handler) ListItems(w http.ResponseWriter, r *http.Request) {
	player, ok := currentPlayer(w, r)
	if !ok {
		return
	}
	if h.db == nil {
		httpx.WriteProblem(w, r, httpx.ServiceUnavailable("database is unavailable"))
		return
	}
	if h.content == nil {
		httpx.WriteProblem(w, r, httpx.ServiceUnavailable("content store is unavailable"))
		return
	}

	records, err := h.findPlayerItems(r.Context(), player.ID)
	if err != nil {
		httpx.WriteProblem(w, r, httpx.InternalServerError("items unavailable", "me_items_lookup_failed", err))
		return
	}

	items, err := mapPlayerItems(h.content, records)
	if err != nil {
		httpx.WriteProblem(w, r, httpx.InternalServerError("item inventory is inconsistent", "me_items_response_failed", err))
		return
	}

	httpx.WriteJSON(w, http.StatusOK, ItemListResponse{
		Items: items,
	})
}

func (h *Handler) findPlayerItems(ctx context.Context, playerID string) ([]mongomodel.PlayerItem, error) {
	cursor, err := h.db.Collection(mongomodel.PlayerItemsCollection).Find(
		ctx,
		bson.M{"player_id": playerID, "quantity": bson.M{"$gt": 0}},
		options.Find().SetSort(bson.D{{Key: "item_id", Value: 1}}),
	)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = cursor.Close(ctx)
	}()

	var records []mongomodel.PlayerItem
	if err := cursor.All(ctx, &records); err != nil {
		return nil, err
	}
	return records, nil
}

func mapPlayerItems(store *content.Store, records []mongomodel.PlayerItem) ([]PlayerItemResponse, error) {
	out := make([]PlayerItemResponse, 0, len(records))
	for _, record := range records {
		item, ok := store.GetItem(record.ItemID)
		if !ok {
			continue
		}
		out = append(out, PlayerItemResponse{
			ID:       record.ID,
			ItemID:   record.ItemID,
			Quantity: record.Quantity,
			Item: ItemResponse{
				ID:          item.ID,
				Name:        item.Name,
				Type:        item.Type,
				Rarity:      item.Rarity,
				Description: item.Description,
				IconPath:    item.IconPath,
				Source:      item.Source,
			},
		})
	}
	return out, nil
}
