package config

import (
	"strings"
	"testing"
)

func TestParseGroupTeamMap(t *testing.T) {
	got, err := ParseGroupTeamMap("-1001=team-001, -1002 = team-002")
	if err != nil {
		t.Fatalf("ParseGroupTeamMap returned error: %v", err)
	}
	if len(got) != 2 || got[-1001] != "team-001" || got[-1002] != "team-002" {
		t.Fatalf("unexpected group team map: %#v", got)
	}
}

func TestParseGroupTeamMapRejectsDuplicateChatID(t *testing.T) {
	if _, err := ParseGroupTeamMap("-1001=team-001,-1001=team-002"); err == nil {
		t.Fatalf("expected duplicate chat id error")
	}
}

func TestParseGroupTeamMapRejectsDuplicateTeamID(t *testing.T) {
	if _, err := ParseGroupTeamMap("-1001=team-001,-1002=team-001"); err == nil {
		t.Fatalf("expected duplicate team id error")
	}
}

func TestParseGroupTeamMapRequiresEntries(t *testing.T) {
	if _, err := ParseGroupTeamMap(" "); err == nil {
		t.Fatalf("expected missing map error")
	}
}

func TestListValueTrimsAndDeduplicates(t *testing.T) {
	t.Setenv("INITIAL_SITONE_IDS", "stone-a, stone-b, stone-a,,")

	got := listValue("INITIAL_SITONE_IDS", []string{"fallback"})
	if strings.Join(got, ",") != "stone-a,stone-b" {
		t.Fatalf("unexpected list value: %#v", got)
	}
}
