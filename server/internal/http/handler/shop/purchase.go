package shop

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"github.com/sitcon-tw/camp2026-game/internal/content"
	"github.com/sitcon-tw/camp2026-game/internal/http/httpx"
	mongomodel "github.com/sitcon-tw/camp2026-game/internal/mongodb/model"
)

const purchaseQuantity = 1

var errInsufficientOpenPower = errors.New("insufficient open power")

// Purchase godoc
// @Summary Purchase shop item
// @Description Purchases one enabled shop item, deducts open power, and adds it to the current player's item bag.
// @Tags shop
// @Accept json
// @Produce json
// @Security AuthCookieAuth
// @Param request body PurchaseRequest true "Purchase request"
// @Success 201 {object} PurchaseResponse
// @Failure 400 {object} httpx.ProblemDetails
// @Failure 401 {object} httpx.ProblemDetails
// @Failure 404 {object} httpx.ProblemDetails
// @Failure 409 {object} httpx.ProblemDetails
// @Failure 422 {object} httpx.ProblemDetails
// @Failure 500 {object} httpx.ProblemDetails
// @Failure 503 {object} httpx.ProblemDetails
// @Router /shop/purchases [post]
func (h *Handler) Purchase(w http.ResponseWriter, r *http.Request) {
	player, ok := currentPlayer(w, r)
	if !ok || !h.requireDatabase(w, r) || !h.requireMongoClient(w, r) || !h.requireContent(w, r) {
		return
	}

	var body PurchaseRequest
	if err := httpx.DecodeJSON(r, &body); err != nil {
		httpx.WriteProblem(w, r, err)
		return
	}
	body.ItemID = strings.TrimSpace(body.ItemID)
	if err := httpx.ValidateStruct(body); err != nil {
		httpx.WriteProblem(w, r, err)
		return
	}

	item, found := shopItemByID(h.content, body.ItemID)
	if !found {
		httpx.WriteProblem(w, r, httpx.NotFound("shop item not found"))
		return
	}

	result, err := h.purchaseItem(r.Context(), player.ID, item)
	if err != nil {
		if errors.Is(err, errInsufficientOpenPower) {
			httpx.WriteProblem(w, r, httpx.NewError(http.StatusConflict, "insufficient open power"))
			return
		}
		httpx.WriteProblem(w, r, httpx.NewError(http.StatusInternalServerError, "purchase failed"))
		return
	}

	httpx.WriteJSON(w, http.StatusCreated, PurchaseResponse{
		PurchaseID:     result.purchaseID,
		ItemID:         item.ID,
		Quantity:       purchaseQuantity,
		PriceOpenPower: item.PriceOpenPower,
		OpenPower:      result.openPower,
		Item:           shopItemResponse(item, true),
	})
}

func (h *Handler) requireMongoClient(w http.ResponseWriter, r *http.Request) bool {
	if h.client == nil {
		httpx.WriteProblem(w, r, httpx.ServiceUnavailable("database is unavailable"))
		return false
	}
	return true
}

type purchaseResult struct {
	purchaseID string
	openPower  int
}

func (h *Handler) purchaseItem(ctx context.Context, playerID string, item content.Item) (purchaseResult, error) {
	session, err := h.client.StartSession()
	if err != nil {
		return purchaseResult{}, err
	}
	defer session.EndSession(ctx)

	var result purchaseResult
	err = mongo.WithSession(ctx, session, func(ctx context.Context) error {
		if err := session.StartTransaction(); err != nil {
			return err
		}
		committed := false
		defer func() {
			if !committed {
				_ = session.AbortTransaction(context.Background())
			}
		}()

		openPower, err := h.sumOpenPower(ctx, playerID)
		if err != nil {
			return err
		}
		if openPower < item.PriceOpenPower {
			return errInsufficientOpenPower
		}

		now := time.Now()
		purchaseID := newID("purchase")
		if err := h.insertPurchase(ctx, purchaseID, playerID, item, now); err != nil {
			return err
		}
		if err := h.insertOpenPowerDeduction(ctx, purchaseID, playerID, item.PriceOpenPower, now); err != nil {
			return err
		}
		if err := h.incrementPlayerItem(ctx, playerID, item.ID); err != nil {
			return err
		}

		if err := session.CommitTransaction(ctx); err != nil {
			return err
		}
		committed = true
		result = purchaseResult{
			purchaseID: purchaseID,
			openPower:  openPower - item.PriceOpenPower,
		}
		return nil
	})
	if err != nil {
		return purchaseResult{}, err
	}
	return result, nil
}

func (h *Handler) insertPurchase(ctx context.Context, purchaseID string, playerID string, item content.Item, createdAt time.Time) error {
	_, err := h.db.Collection(mongomodel.ShopPurchasesCollection).InsertOne(ctx, mongomodel.ShopPurchase{
		ID:             purchaseID,
		PlayerID:       playerID,
		ItemID:         item.ID,
		Quantity:       purchaseQuantity,
		PriceOpenPower: item.PriceOpenPower,
		CreatedAt:      createdAt,
	})
	return err
}

func (h *Handler) insertOpenPowerDeduction(ctx context.Context, purchaseID string, playerID string, price int, createdAt time.Time) error {
	_, err := h.db.Collection(mongomodel.OpenPowerRecordsCollection).InsertOne(ctx, mongomodel.OpenPowerRecord{
		ID:        newID("open_power"),
		PlayerID:  playerID,
		Amount:    -price,
		Reason:    "shop_purchase",
		Source:    purchaseID,
		CreatedAt: createdAt,
	})
	return err
}

func (h *Handler) incrementPlayerItem(ctx context.Context, playerID string, itemID string) error {
	_, err := h.db.Collection(mongomodel.PlayerItemsCollection).UpdateOne(
		ctx,
		bson.M{
			"player_id": playerID,
			"item_id":   itemID,
		},
		bson.M{
			"$setOnInsert": bson.M{
				"_id":       newID("player_item"),
				"player_id": playerID,
				"item_id":   itemID,
			},
			"$inc": bson.M{"quantity": purchaseQuantity},
		},
		options.UpdateOne().SetUpsert(true),
	)
	return err
}

func (h *Handler) sumOpenPower(ctx context.Context, playerID string) (int, error) {
	cursor, err := h.db.Collection(mongomodel.OpenPowerRecordsCollection).
		Aggregate(ctx, openPowerTotalPipeline(playerID))
	if err != nil {
		return 0, err
	}
	defer func() {
		_ = cursor.Close(ctx)
	}()

	return openPowerTotalFromCursor(ctx, cursor)
}

func openPowerTotalFromCursor(ctx context.Context, cursor *mongo.Cursor) (int, error) {
	var totals []struct {
		Total int `bson:"total"`
	}
	if err := cursor.All(ctx, &totals); err != nil {
		return 0, err
	}
	if len(totals) == 0 {
		return 0, nil
	}
	return totals[0].Total, nil
}

func openPowerTotalPipeline(playerID string) mongo.Pipeline {
	return mongo.Pipeline{
		{{Key: "$match", Value: bson.D{{Key: "player_id", Value: playerID}}}},
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: nil},
			{Key: "total", Value: bson.D{{Key: "$sum", Value: "$amount"}}},
		}}},
	}
}

func newID(prefix string) string {
	return prefix + "_" + bson.NewObjectID().Hex()
}
