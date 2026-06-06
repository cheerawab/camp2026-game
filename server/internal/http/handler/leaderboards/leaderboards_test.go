package leaderboards

import "testing"

func TestNormalizeTypeDefaultsToOpenPower(t *testing.T) {
	if got := normalizeType(""); got != TypeOpenPower {
		t.Fatalf("expected default type %q, got %q", TypeOpenPower, got)
	}
}

func TestIsValidType(t *testing.T) {
	for _, value := range []string{TypeOpenPower, TypeSitones, TypeMatches} {
		if !isValidType(value) {
			t.Fatalf("expected %q to be valid", value)
		}
	}
	if isValidType("unknown") {
		t.Fatal("expected unknown type to be invalid")
	}
}

func TestMetricForType(t *testing.T) {
	if got := metricForType(TypeOpenPower); got != "OP" {
		t.Fatalf("expected OP metric, got %q", got)
	}
	if got := metricForType(TypeSitones); got != "小石" {
		t.Fatalf("expected sitone metric, got %q", got)
	}
	if got := metricForType(TypeMatches); got != "分" {
		t.Fatalf("expected match metric, got %q", got)
	}
}

func TestLeaderboardPipelines(t *testing.T) {
	if got := len(openPowerScoresByTeamPipeline()); got != 4 {
		t.Fatalf("expected 4 open power stages, got %d", got)
	}
	if got := len(inventoryScoresByTeamPipeline()); got != 5 {
		t.Fatalf("expected 5 inventory stages, got %d", got)
	}
	if got := len(matchScoresByTeamPipeline()); got != 6 {
		t.Fatalf("expected 6 match stages, got %d", got)
	}
}
