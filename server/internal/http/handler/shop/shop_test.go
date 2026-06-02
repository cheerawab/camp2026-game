package shop

import (
	"testing"

	"github.com/sitcon-tw/camp2026-game/internal/content"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func TestShopItemsOnlyIncludesEnabledPurchasableItems(t *testing.T) {
	store := loadTestContent(t)

	items := shopItems(store)
	if len(items) != 2 {
		t.Fatalf("expected 2 shop items, got %#v", items)
	}
	if items[0].ID != "item-crafting-fragment" || items[0].PriceOpenPower != 50 {
		t.Fatalf("unexpected first shop item: %#v", items[0])
	}
	if items[1].ID != "item-theme-ticket" || items[1].PriceOpenPower != 200 {
		t.Fatalf("unexpected second shop item: %#v", items[1])
	}

	if _, ok := shopItemByID(store, "item-memory-tag"); ok {
		t.Fatal("expected non-purchasable item to be excluded")
	}
}

func TestShopItemResponse(t *testing.T) {
	response := shopItemResponse(content.Item{
		ID:             "item-crafting-fragment",
		Name:           "合成碎片",
		Type:           "material",
		Rarity:         "common",
		Description:    "小石造型合成使用的基礎素材。",
		PriceOpenPower: 50,
	})

	if response.ID != "item-crafting-fragment" || response.PriceOpenPower != 50 {
		t.Fatalf("unexpected shop item response: %#v", response)
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

func loadTestContent(t *testing.T) *content.Store {
	t.Helper()

	store, err := content.Load("../../../../content")
	if err != nil {
		t.Fatalf("load test content: %v", err)
	}
	return store
}
