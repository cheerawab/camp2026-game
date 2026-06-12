package main

import (
	"strings"
	"testing"
)

const validHash = "5b168d8b83a4f506cf9be7b9b86b7ffc2b7a8dee29e5f6c82f9720764fc05f4d"

func TestBuildImportPlanCreatesRegularAndStaffPlayers(t *testing.T) {
	csv := `組別,暱稱,職稱,email (sha256),備註,Token,Staff Token
總召組,Windless,總召,` + validHash + `,,token-1,staff-token-1
`

	plan, err := buildImportPlan(strings.NewReader(csv), 6)
	if err != nil {
		t.Fatalf("build import plan: %v", err)
	}

	if len(plan.Teams) != 1 {
		t.Fatalf("expected 1 team, got %d", len(plan.Teams))
	}
	if plan.Teams[0].ID != "team-001" || plan.Teams[0].Name != "Team 001" {
		t.Fatalf("unexpected team: %#v", plan.Teams[0])
	}
	if len(plan.Players) != 2 {
		t.Fatalf("expected 2 players, got %d", len(plan.Players))
	}

	regular := plan.Players[0]
	if regular.ID != "token-1" || regular.AuthToken != "token-1" {
		t.Fatalf("unexpected regular token fields: %#v", regular)
	}
	if regular.QRCodeToken == "" {
		t.Fatalf("expected regular qr code identifier")
	}
	if regular.QRCodeToken == regular.AuthToken || regular.QRCodeToken == "staff-token-1" {
		t.Fatalf("expected qr code identifier to differ from auth tokens, got %#v", regular)
	}
	if !strings.HasPrefix(regular.QRCodeToken, "qr_") {
		t.Fatalf("expected qr code identifier prefix, got %q", regular.QRCodeToken)
	}
	if regular.Role != "" {
		t.Fatalf("expected regular player role to be empty, got %q", regular.Role)
	}
	if regular.TeamID != "team-001" {
		t.Fatalf("expected regular team team-001, got %q", regular.TeamID)
	}
	if regular.AvatarURL != "https://www.gravatar.com/avatar/"+validHash {
		t.Fatalf("unexpected regular avatar url: %q", regular.AvatarURL)
	}

	staff := plan.Players[1]
	if staff.ID != "staff-token-1" || staff.AuthToken != "staff-token-1" {
		t.Fatalf("unexpected staff token fields: %#v", staff)
	}
	if staff.QRCodeToken != "" {
		t.Fatalf("expected staff qr code identifier to be empty, got %q", staff.QRCodeToken)
	}
	if staff.Role != "staff" {
		t.Fatalf("expected staff role, got %q", staff.Role)
	}
	if staff.TeamID != "" {
		t.Fatalf("expected staff team to be empty, got %q", staff.TeamID)
	}
}

func TestBuildImportPlanAssignsSixRowsPerTeam(t *testing.T) {
	var builder strings.Builder
	builder.WriteString("組別,暱稱,email (sha256),Token,Staff Token\n")
	for i := 1; i <= 7; i++ {
		builder.WriteString("組,Player,")
		builder.WriteString(validHash)
		builder.WriteString(",token-")
		builder.WriteString(string(rune('0' + i)))
		builder.WriteString(",staff-token-")
		builder.WriteString(string(rune('0' + i)))
		builder.WriteString("\n")
	}

	plan, err := buildImportPlan(strings.NewReader(builder.String()), 6)
	if err != nil {
		t.Fatalf("build import plan: %v", err)
	}

	if len(plan.Teams) != 2 {
		t.Fatalf("expected 2 teams, got %d", len(plan.Teams))
	}
	if plan.Players[0].TeamID != "team-001" {
		t.Fatalf("expected first row in team-001, got %q", plan.Players[0].TeamID)
	}
	if plan.Players[10].TeamID != "team-001" {
		t.Fatalf("expected sixth row in team-001, got %q", plan.Players[10].TeamID)
	}
	if plan.Players[12].TeamID != "team-002" {
		t.Fatalf("expected seventh row in team-002, got %q", plan.Players[12].TeamID)
	}
}

func TestBuildImportPlanRejectsDuplicateTokens(t *testing.T) {
	csv := `組別,暱稱,email (sha256),Token,Staff Token
組,A,` + validHash + `,same-token,staff-token-1
組,B,` + validHash + `,token-2,same-token
`

	_, err := buildImportPlan(strings.NewReader(csv), 6)
	if err == nil {
		t.Fatalf("expected duplicate token error")
	}
}

func TestBuildImportPlanRejectsMissingRequiredValue(t *testing.T) {
	csv := `組別,暱稱,email (sha256),Token,Staff Token
組,A,` + validHash + `,,staff-token-1
`

	_, err := buildImportPlan(strings.NewReader(csv), 6)
	if err == nil {
		t.Fatalf("expected missing token error")
	}
}

func TestBuildImportPlanRejectsInvalidEmailHash(t *testing.T) {
	csv := `組別,暱稱,email (sha256),Token,Staff Token
組,A,not-a-hash,token-1,staff-token-1
`

	_, err := buildImportPlan(strings.NewReader(csv), 6)
	if err == nil {
		t.Fatalf("expected invalid hash error")
	}
}
