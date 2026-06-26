package shop

import (
	"fmt"
	"testing"

	"github.com/sitcon-tw/camp2026-game/internal/content"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func TestShopItemsIncludesAllEnabledPurchasableContentItems(t *testing.T) {
	store := loadTestContent(t)

	items := shopItems(store)
	if len(items) != 33 {
		t.Fatalf("expected 33 shop items, got %#v", items)
	}
	if items[0].ID != "item_adventure_backpack" || items[0].PriceOpenPower != 50 {
		t.Fatalf("unexpected first shop item: %#v", items[0])
	}
	if _, ok := shopItemByID(store, "item_polaroid_film"); ok {
		t.Fatal("expected drop-only polaroid film item not to be purchasable")
	}
	if _, ok := shopItemByID(store, "item_shared_notes_link"); !ok {
		t.Fatal("expected dual-source shared notes item to be purchasable")
	}
	if _, ok := shopItemByID(store, "item_charm_debug"); !ok {
		t.Fatal("expected charm item to be purchasable")
	}
}

func TestShopItemResponse(t *testing.T) {
	response := shopItemResponse(content.Item{
		ID:             "item_adventure_backpack",
		Name:           "冒險背包",
		Type:           "material",
		Rarity:         "common",
		Description:    "冒險背包，可用於小石合成。",
		Source:         "shop",
		PriceOpenPower: 50,
	}, true)

	if response.ID != "item_adventure_backpack" || response.Source != "shop" || response.PriceOpenPower != 50 || !response.Redeemed {
		t.Fatalf("unexpected shop item response: %#v", response)
	}
}

func TestShopItemResponsesIncludesRedeemedState(t *testing.T) {
	responses := shopItemResponses([]content.Item{
		{ID: "item-a", Name: "A", Type: "material", Rarity: "common", PriceOpenPower: 10},
		{ID: "item-b", Name: "B", Type: "material", Rarity: "common", PriceOpenPower: 20},
	}, map[string]struct{}{"item-b": {}})

	if len(responses) != 2 {
		t.Fatalf("expected 2 responses, got %#v", responses)
	}
	if responses[0].Redeemed {
		t.Fatalf("expected first item not redeemed: %#v", responses[0])
	}
	if !responses[1].Redeemed {
		t.Fatalf("expected second item redeemed: %#v", responses[1])
	}
}

func TestOpenPowerTotalPipeline(t *testing.T) {
	pipeline := openPowerTotalPipeline("player-a")
	if len(pipeline) != 2 {
		t.Fatalf("expected 2 pipeline stages, got %#v", pipeline)
	}

	matchStage, ok := pipeline[0][0].Value.(bson.D)
	if !ok {
		t.Fatalf("expected match stage document, got %#v", pipeline[0][0].Value)
	}
	var got any
	for _, element := range matchStage {
		if element.Key == "player_id" {
			got = element.Value
			break
		}
	}
	if got != "player-a" {
		t.Fatalf("expected player id match, got %#v", got)
	}
}

func TestTransactionUnsupported(t *testing.T) {
	err := fmt.Errorf("wrapped: %w", mongo.CommandError{
		Code:    20,
		Message: "Transaction numbers are only allowed on a replica set member or mongos",
	})
	if !transactionUnsupported(err) {
		t.Fatal("expected standalone transaction error to be unsupported")
	}

	for _, err := range []error{
		mongo.CommandError{Code: 19, Message: "Transaction numbers are only allowed"},
		mongo.CommandError{Code: 20, Message: "not a transaction error"},
		fmt.Errorf("plain error"),
	} {
		if transactionUnsupported(err) {
			t.Fatalf("expected %v not to be unsupported", err)
		}
	}
}

func loadTestContent(t *testing.T) *content.Store {
	t.Helper()

	store, err := content.Load("../../../../content")
	if err != nil {
		t.Fatalf("load test content: %v", err)
	}
	return store
}
