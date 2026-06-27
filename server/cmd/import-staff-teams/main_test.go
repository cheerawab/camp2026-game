package main

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	mongomodel "github.com/sitcon-tw/camp2026-game/internal/mongodb/model"
)

func TestBundledAssignmentsMatchSheetDecisions(t *testing.T) {
	file, err := os.Open(filepath.Join("..", "..", "content", "staff_team_assignments.json"))
	if err != nil {
		t.Fatalf("open bundled assignments: %v", err)
	}
	defer func() {
		_ = file.Close()
	}()

	assignments, err := readAssignments(file)
	if err != nil {
		t.Fatalf("read assignments: %v", err)
	}

	if got := len(staffNames(assignments)); got != 55 {
		t.Fatalf("expected 55 staff assignments, got %d", got)
	}
	if !teamHasStaff(assignments, "team-004", "Tang Yu") {
		t.Fatalf("expected team-004 to include Tang Yu")
	}
	if teamHasStaff(assignments, "team-004", "Tang") {
		t.Fatalf("did not expect abbreviated Tang assignment")
	}
	if teamHasStaff(assignments, "team-010", "CC") {
		t.Fatalf("did not expect duplicate CC in team-010")
	}
	if !teamHasStaff(assignments, "team-006", "CC") {
		t.Fatalf("expected CC in team-006")
	}
}

func TestReadAssignmentsRejectsDuplicateStaff(t *testing.T) {
	_, err := readAssignments(strings.NewReader(`{
		"teams": [
			{"teamId":"team-001","name":"Team 001","staff":["A"]},
			{"teamId":"team-002","name":"Team 002","staff":["A"]},
			{"teamId":"team-003","name":"Team 003","staff":["C"]},
			{"teamId":"team-004","name":"Team 004","staff":["D"]},
			{"teamId":"team-005","name":"Team 005","staff":["E"]},
			{"teamId":"team-006","name":"Team 006","staff":["F"]},
			{"teamId":"team-007","name":"Team 007","staff":["G"]},
			{"teamId":"team-008","name":"Team 008","staff":["H"]},
			{"teamId":"team-009","name":"Team 009","staff":["I"]},
			{"teamId":"team-010","name":"Team 010","staff":["J"]}
		]
	}`))
	if err == nil {
		t.Fatalf("expected duplicate staff error")
	}
}

func TestBuildImportPlanCreatesTeamUpdates(t *testing.T) {
	assignments := testAssignments()
	players := []mongomodel.Player{
		{ID: "staff-a", Nickname: "A", Role: playerRoleStaff},
		{ID: "staff-b", Nickname: "B", Role: playerRoleStaff},
		{ID: "staff-c", Nickname: "C", Role: playerRoleStaff},
		{ID: "staff-d", Nickname: "D", Role: playerRoleStaff},
		{ID: "staff-e", Nickname: "E", Role: playerRoleStaff},
		{ID: "staff-f", Nickname: "F", Role: playerRoleStaff},
		{ID: "staff-g", Nickname: "G", Role: playerRoleStaff},
		{ID: "staff-h", Nickname: "H", Role: playerRoleStaff},
		{ID: "staff-i", Nickname: "I", Role: playerRoleStaff},
		{ID: "staff-j", Nickname: "J", Role: playerRoleStaff},
		{ID: "player-a", Nickname: "A"},
	}

	plan, err := buildImportPlan(assignments, players)
	if err != nil {
		t.Fatalf("build import plan: %v", err)
	}
	if len(plan.Teams) != 10 {
		t.Fatalf("expected 10 teams, got %d", len(plan.Teams))
	}
	if len(plan.Updates) != 10 {
		t.Fatalf("expected 10 updates, got %d", len(plan.Updates))
	}
	if plan.Updates[0].PlayerID != "staff-a" || plan.Updates[0].TeamID != "team-001" {
		t.Fatalf("unexpected first update: %#v", plan.Updates[0])
	}
}

func TestBuildImportPlanRejectsMissingStaffPlayer(t *testing.T) {
	_, err := buildImportPlan(testAssignments(), []mongomodel.Player{
		{ID: "staff-a", Nickname: "A", Role: playerRoleStaff},
	})
	if err == nil {
		t.Fatalf("expected missing staff player error")
	}
}

func TestBuildImportPlanRejectsDuplicateStaffPlayer(t *testing.T) {
	players := testStaffPlayers()
	players = append(players, mongomodel.Player{ID: "staff-a-2", Nickname: "A", Role: playerRoleStaff})

	_, err := buildImportPlan(testAssignments(), players)
	if err == nil {
		t.Fatalf("expected duplicate staff player error")
	}
}

func TestPrintSummaryDoesNotExposePlayerID(t *testing.T) {
	plan := importPlan{
		Teams: []mongomodel.Team{{ID: "team-001", Name: "Team 001"}},
		Updates: []staffTeamUpdate{{
			PlayerID: "secret-staff-token",
			Nickname: "A",
			TeamID:   "team-001",
		}},
	}

	output := captureStdout(t, func() {
		printSummary(plan, true)
	})
	if strings.Contains(output, "secret-staff-token") {
		t.Fatalf("expected summary to redact player id, got %q", output)
	}
	if !strings.Contains(output, "A -> team-001") {
		t.Fatalf("expected summary to include nickname and team id, got %q", output)
	}
}

func testAssignments() assignmentFile {
	return assignmentFile{Teams: []staffTeamAssignment{
		{TeamID: "team-001", Name: "Team 001", Staff: []string{"A"}},
		{TeamID: "team-002", Name: "Team 002", Staff: []string{"B"}},
		{TeamID: "team-003", Name: "Team 003", Staff: []string{"C"}},
		{TeamID: "team-004", Name: "Team 004", Staff: []string{"D"}},
		{TeamID: "team-005", Name: "Team 005", Staff: []string{"E"}},
		{TeamID: "team-006", Name: "Team 006", Staff: []string{"F"}},
		{TeamID: "team-007", Name: "Team 007", Staff: []string{"G"}},
		{TeamID: "team-008", Name: "Team 008", Staff: []string{"H"}},
		{TeamID: "team-009", Name: "Team 009", Staff: []string{"I"}},
		{TeamID: "team-010", Name: "Team 010", Staff: []string{"J"}},
	}}
}

func testStaffPlayers() []mongomodel.Player {
	players := make([]mongomodel.Player, 0, 10)
	for _, name := range []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J"} {
		players = append(players, mongomodel.Player{
			ID:       "staff-" + strings.ToLower(name),
			Nickname: name,
			Role:     playerRoleStaff,
		})
	}
	return players
}

func teamHasStaff(assignments assignmentFile, teamID string, nickname string) bool {
	for _, team := range assignments.Teams {
		if team.TeamID != teamID {
			continue
		}
		for _, staff := range team.Staff {
			if staff == nickname {
				return true
			}
		}
	}
	return false
}

func captureStdout(t *testing.T, fn func()) string {
	t.Helper()

	originalStdout := os.Stdout
	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatalf("create stdout pipe: %v", err)
	}
	os.Stdout = writer
	defer func() {
		os.Stdout = originalStdout
	}()

	fn()

	if err := writer.Close(); err != nil {
		t.Fatalf("close stdout writer: %v", err)
	}
	output, err := io.ReadAll(reader)
	if err != nil {
		t.Fatalf("read stdout: %v", err)
	}
	if err := reader.Close(); err != nil {
		t.Fatalf("close stdout reader: %v", err)
	}
	return string(output)
}
