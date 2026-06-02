package me

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"github.com/sitcon-tw/camp2026-game/internal/content"
	"github.com/sitcon-tw/camp2026-game/internal/http/authctx"
	"github.com/sitcon-tw/camp2026-game/internal/http/httpx"
	mongomodel "github.com/sitcon-tw/camp2026-game/internal/mongodb/model"
)

func TestQRCodeResponse(t *testing.T) {
	handler := New(Dependencies{})
	req := authenticatedRequest(mongomodel.Player{
		ID:          "7H9K2Q",
		AuthToken:   "auth_token_123456",
		QRCodeToken: "qr_token_123456",
	})
	res := httptest.NewRecorder()

	handler.QRCode(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d: %s", http.StatusOK, res.Code, res.Body.String())
	}

	var body map[string]string
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body["qrcodeToken"] != "qr_token_123456" {
		t.Fatalf("expected qrcode token, got %#v", body)
	}
	if _, ok := body["authToken"]; ok {
		t.Fatalf("expected auth token to be omitted, got %#v", body)
	}
}

func TestQRCodeRequiresPlayerContext(t *testing.T) {
	handler := New(Dependencies{})
	req := httptest.NewRequest(http.MethodGet, "/api/me/qrcode", nil)
	res := httptest.NewRecorder()

	handler.QRCode(res, req)

	assertProblem(t, res, http.StatusUnauthorized)
}

func TestQRCodeRequiresToken(t *testing.T) {
	handler := New(Dependencies{})
	req := authenticatedRequest(mongomodel.Player{ID: "7H9K2Q"})
	res := httptest.NewRecorder()

	handler.QRCode(res, req)

	assertProblem(t, res, http.StatusInternalServerError)
}

func TestStatusResponse(t *testing.T) {
	response := statusResponse(
		mongomodel.Player{
			ID:        "7H9K2Q",
			Nickname:  "Alice",
			TeamID:    "8M4RXP",
			AvatarURL: "https://example.test/avatar/alice.png",
		},
		mongomodel.Team{
			ID:   "8M4RXP",
			Name: "Blue Team",
		},
		1280,
	)

	if response.PlayerID != "7H9K2Q" {
		t.Fatalf("expected player id, got %q", response.PlayerID)
	}
	if response.Team.TeamID != "8M4RXP" {
		t.Fatalf("expected team id, got %q", response.Team.TeamID)
	}
	if response.OpenPower != 1280 {
		t.Fatalf("expected open power 1280, got %d", response.OpenPower)
	}
	if response.AvatarURL == "" {
		t.Fatalf("expected avatar url")
	}
}

func TestOpenPowerTotalPipeline(t *testing.T) {
	got := openPowerTotalPipeline("7H9K2Q")
	want := mongo.Pipeline{
		{{Key: "$match", Value: bson.D{{Key: "player_id", Value: "7H9K2Q"}}}},
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: nil},
			{Key: "total", Value: bson.D{{Key: "$sum", Value: "$amount"}}},
		}}},
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected pipeline: %#v", got)
	}
}

func TestOpenPowerTotalFromCursor(t *testing.T) {
	cursor, err := mongo.NewCursorFromDocuments([]any{
		bson.D{{Key: "total", Value: 1280}},
	}, nil, nil)
	if err != nil {
		t.Fatalf("new cursor: %v", err)
	}

	total, err := openPowerTotalFromCursor(context.Background(), cursor)
	if err != nil {
		t.Fatalf("open power total: %v", err)
	}
	if total != 1280 {
		t.Fatalf("expected total 1280, got %d", total)
	}
}

func TestMapPlayerSitones(t *testing.T) {
	sitones, err := mapPlayerSitones(loadTestContent(t), []mongomodel.PlayerSitone{
		{
			ID:       "owned-sitone-001",
			PlayerID: "7H9K2Q",
			SitoneID: "sitone-engineering",
			Quantity: 1,
		},
	})
	if err != nil {
		t.Fatalf("map sitones: %v", err)
	}
	if len(sitones) != 1 {
		t.Fatalf("expected 1 sitone, got %d", len(sitones))
	}
	if sitones[0].Sitone.Name != "工程型小石" {
		t.Fatalf("expected catalog sitone name, got %#v", sitones[0])
	}
}

func TestMapPlayerSitonesRequiresCatalogDefinition(t *testing.T) {
	_, err := mapPlayerSitones(loadTestContent(t), []mongomodel.PlayerSitone{
		{ID: "owned-sitone-001", SitoneID: "sitone-missing", Quantity: 1},
	})
	if err == nil {
		t.Fatal("expected missing sitone error")
	}
}

func TestMapPlayerItems(t *testing.T) {
	items, err := mapPlayerItems(loadTestContent(t), []mongomodel.PlayerItem{
		{
			ID:       "owned-item-001",
			PlayerID: "7H9K2Q",
			ItemID:   "item-crafting-fragment",
			Quantity: 3,
		},
	})
	if err != nil {
		t.Fatalf("map items: %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(items))
	}
	if items[0].Item.Name != "合成碎片" {
		t.Fatalf("expected catalog item name, got %#v", items[0])
	}
}

func TestMapPlayerItemsReturnsEmptySlice(t *testing.T) {
	items, err := mapPlayerItems(loadTestContent(t), nil)
	if err != nil {
		t.Fatalf("map items: %v", err)
	}
	if items == nil {
		t.Fatal("expected empty slice, got nil")
	}
	if len(items) != 0 {
		t.Fatalf("expected 0 items, got %d", len(items))
	}
}

func authenticatedRequest(player mongomodel.Player) *http.Request {
	req := httptest.NewRequest(http.MethodGet, "/api/me/qrcode", strings.NewReader(""))
	return req.WithContext(authctx.WithPlayer(req.Context(), player))
}

func assertProblem(t *testing.T, res *httptest.ResponseRecorder, status int) httpx.ProblemDetails {
	t.Helper()

	if res.Code != status {
		t.Fatalf("expected status %d, got %d: %s", status, res.Code, res.Body.String())
	}
	if contentType := res.Header().Get("Content-Type"); contentType != "application/problem+json" {
		t.Fatalf("expected problem content type, got %q", contentType)
	}

	var problem httpx.ProblemDetails
	if err := json.NewDecoder(res.Body).Decode(&problem); err != nil {
		t.Fatalf("decode problem: %v", err)
	}
	if problem.Status != status {
		t.Fatalf("expected problem status %d, got %d", status, problem.Status)
	}
	return problem
}

func loadTestContent(t *testing.T) *content.Store {
	t.Helper()

	store, err := content.Load("../../../../content")
	if err != nil {
		t.Fatalf("load test content: %v", err)
	}
	return store
}
