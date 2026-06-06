package fusions

import (
	"testing"

	"github.com/sitcon-tw/camp2026-game/internal/content"
)

func TestRecipeAvailable(t *testing.T) {
	recipe := content.FusionRecipe{
		Inputs: []content.FusionComponent{
			{Kind: content.FusionKindSitone, ID: "sitone-engineering", Quantity: 1},
			{Kind: content.FusionKindItem, ID: "item-crafting-fragment", Quantity: 3},
		},
	}
	inventory := inventoryCounts{
		sitones: map[string]int{"sitone-engineering": 1},
		items:   map[string]int{"item-crafting-fragment": 3},
	}
	if !recipeAvailable(recipe, inventory) {
		t.Fatal("expected recipe to be available")
	}

	inventory.items["item-crafting-fragment"] = 2
	if recipeAvailable(recipe, inventory) {
		t.Fatal("expected recipe to be unavailable with insufficient items")
	}
}

func TestInventoryCollection(t *testing.T) {
	collection, field, err := inventoryCollection(content.FusionKindItem)
	if err != nil {
		t.Fatalf("item collection: %v", err)
	}
	if collection != "player_items" || field != "item_id" {
		t.Fatalf("unexpected item collection mapping: %s %s", collection, field)
	}

	collection, field, err = inventoryCollection(content.FusionKindSitone)
	if err != nil {
		t.Fatalf("sitone collection: %v", err)
	}
	if collection != "player_sitones" || field != "sitone_id" {
		t.Fatalf("unexpected sitone collection mapping: %s %s", collection, field)
	}
}

func TestModelComponents(t *testing.T) {
	components := modelComponents([]content.FusionComponent{
		{Kind: content.FusionKindItem, ID: "item-crafting-fragment", Quantity: 3},
	})
	if len(components) != 1 {
		t.Fatalf("expected 1 component, got %#v", components)
	}
	if components[0].RefID != "item-crafting-fragment" || components[0].Quantity != 3 {
		t.Fatalf("unexpected model component: %#v", components[0])
	}
}
