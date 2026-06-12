package main

import (
	"reflect"
	"testing"

	mongomodel "github.com/sitcon-tw/camp2026-game/internal/mongodb/model"
)

var testInitialSitoneIDs = []string{
	"stone_explorer_base",
	"stone_inspiration_base",
	"stone_resonance_base",
	"stone_engineering_base",
	"stone_entertainment_base",
}

func TestBuildGrantPlanGrantsFiveSitonesToEmptyPlayer(t *testing.T) {
	plan := buildGrantPlan(
		[]mongomodel.Player{{ID: "player-a", Nickname: "Alice"}},
		map[string]map[string]int{},
		testInitialSitoneIDs,
		5,
	)

	if len(plan.Grants) != 1 {
		t.Fatalf("expected 1 grant, got %d", len(plan.Grants))
	}
	if !reflect.DeepEqual(plan.Grants[0].SitoneIDs, testInitialSitoneIDs) {
		t.Fatalf("unexpected sitones: %#v", plan.Grants[0].SitoneIDs)
	}
}

func TestBuildGrantPlanSkipsPlayersAlreadyAtLimit(t *testing.T) {
	plan := buildGrantPlan(
		[]mongomodel.Player{{ID: "player-a"}},
		map[string]map[string]int{
			"player-a": {
				"stone_explorer_base":      1,
				"stone_inspiration_base":   1,
				"stone_resonance_base":     1,
				"stone_engineering_base":   1,
				"stone_entertainment_base": 1,
			},
		},
		testInitialSitoneIDs,
		5,
	)

	if len(plan.Grants) != 0 {
		t.Fatalf("expected no grants, got %#v", plan.Grants)
	}
	if plan.SkippedAtLimit != 1 {
		t.Fatalf("expected skipped at limit 1, got %d", plan.SkippedAtLimit)
	}
}

func TestBuildGrantPlanFillsOnlyMissingSitonesUpToLimit(t *testing.T) {
	plan := buildGrantPlan(
		[]mongomodel.Player{{ID: "player-a"}},
		map[string]map[string]int{
			"player-a": {
				"stone_explorer_base":  2,
				"stone_resonance_base": 1,
			},
		},
		testInitialSitoneIDs,
		5,
	)

	if len(plan.Grants) != 1 {
		t.Fatalf("expected 1 grant, got %d", len(plan.Grants))
	}
	want := []string{"stone_inspiration_base", "stone_engineering_base"}
	if !reflect.DeepEqual(plan.Grants[0].SitoneIDs, want) {
		t.Fatalf("unexpected sitones: got %#v want %#v", plan.Grants[0].SitoneIDs, want)
	}
}

func TestParseSitoneIDsRejectsDuplicates(t *testing.T) {
	_, err := parseSitoneIDs("stone_explorer_base,stone_explorer_base")
	if err == nil {
		t.Fatalf("expected duplicate error")
	}
}
