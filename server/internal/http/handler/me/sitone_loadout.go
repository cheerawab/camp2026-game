package me

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

// SitoneLoadout godoc
// @Summary Get current player sitone loadout
// @Description Returns the authenticated player's default sitone loadout.
// @Tags me
// @Produce json
// @Security AuthCookieAuth
// @Success 200 {object} SitoneLoadoutResponse
// @Failure 401 {object} httpx.ProblemDetails
// @Failure 500 {object} httpx.ProblemDetails
// @Failure 503 {object} httpx.ProblemDetails
// @Router /me/sitone-loadout [get]
func (h *Handler) SitoneLoadout(w http.ResponseWriter, r *http.Request) {
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

	sitoneIDs, err := h.defaultSitoneLoadout(r.Context(), player)
	if err != nil {
		httpx.WriteProblem(w, r, httpx.InternalServerError("sitone loadout unavailable", "me_sitone_loadout_lookup_failed", err))
		return
	}
	httpx.WriteJSON(w, http.StatusOK, SitoneLoadoutResponse{SitoneIDs: sitoneIDs})
}

// UpdateSitoneLoadout godoc
// @Summary Update current player sitone loadout
// @Description Updates the authenticated player's default sitone loadout.
// @Tags me
// @Accept json
// @Produce json
// @Security AuthCookieAuth
// @Param request body SitoneLoadoutRequest true "Sitone loadout request"
// @Success 200 {object} SitoneLoadoutResponse
// @Failure 400 {object} httpx.ProblemDetails
// @Failure 401 {object} httpx.ProblemDetails
// @Failure 422 {object} httpx.ProblemDetails
// @Failure 500 {object} httpx.ProblemDetails
// @Failure 503 {object} httpx.ProblemDetails
// @Router /me/sitone-loadout [put]
func (h *Handler) UpdateSitoneLoadout(w http.ResponseWriter, r *http.Request) {
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

	var body SitoneLoadoutRequest
	if err := httpx.DecodeJSON(r, &body); err != nil {
		httpx.WriteProblem(w, r, err)
		return
	}
	if err := httpx.ValidateStruct(body); err != nil {
		httpx.WriteProblem(w, r, err)
		return
	}
	sitoneIDs, err := h.validateOwnedSitoneLoadout(r.Context(), player.ID, body.SitoneIDs)
	if err != nil {
		httpx.WriteProblem(w, r, err)
		return
	}

	if err := h.saveDefaultSitoneLoadout(r.Context(), player.ID, sitoneIDs); err != nil {
		httpx.WriteProblem(w, r, httpx.InternalServerError("sitone loadout update failed", "me_sitone_loadout_save_failed", err))
		return
	}
	httpx.WriteJSON(w, http.StatusOK, SitoneLoadoutResponse{SitoneIDs: sitoneIDs})
}

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

func (h *Handler) validateOwnedSitoneLoadout(rctx context.Context, playerID string, sitoneIDs []string) ([]string, error) {
	normalized, err := normalizeSitoneLoadout(sitoneIDs)
	if err != nil {
		return nil, err
	}

	owned, err := h.ownedSitoneCounts(rctx, playerID)
	if err != nil {
		return nil, httpx.InternalServerError("sitone loadout unavailable", "me_sitone_loadout_inventory_lookup_failed", err)
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

func (h *Handler) defaultSitoneLoadout(rctx context.Context, player mongomodel.Player) ([]string, error) {
	if len(player.DefaultSitoneIDs) > 0 {
		if loadout, err := h.validateOwnedSitoneLoadout(rctx, player.ID, player.DefaultSitoneIDs); err == nil {
			return loadout, nil
		}
	}

	owned, err := h.ownedSitoneCounts(rctx, player.ID)
	if err != nil {
		return nil, err
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

func (h *Handler) ownedSitoneCounts(rctx context.Context, playerID string) (map[string]int, error) {
	cursor, err := h.db.Collection(mongomodel.PlayerSitonesCollection).Find(
		rctx,
		bson.M{"player_id": playerID, "quantity": bson.M{"$gt": 0}},
	)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = cursor.Close(rctx)
	}()

	owned := make(map[string]int)
	for cursor.Next(rctx) {
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

func (h *Handler) saveDefaultSitoneLoadout(rctx context.Context, playerID string, sitoneIDs []string) error {
	_, err := h.db.Collection(mongomodel.PlayersCollection).UpdateOne(
		rctx,
		bson.M{"_id": playerID},
		bson.M{"$set": bson.M{"default_sitone_ids": sitoneIDs}},
	)
	return err
}
