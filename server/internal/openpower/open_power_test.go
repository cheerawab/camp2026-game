package openpower

import (
	"context"
	"reflect"
	"testing"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func TestTotalPipeline(t *testing.T) {
	got := TotalPipeline("7H9K2Q")
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

func TestTotalFromCursor(t *testing.T) {
	cursor, err := mongo.NewCursorFromDocuments([]any{
		bson.D{{Key: "total", Value: 1280}},
	}, nil, nil)
	if err != nil {
		t.Fatalf("new cursor: %v", err)
	}

	total, err := TotalFromCursor(context.Background(), cursor)
	if err != nil {
		t.Fatalf("open power total: %v", err)
	}
	if total != 1280 {
		t.Fatalf("expected total 1280, got %d", total)
	}
}

func TestTotalFromCursorWithoutRecords(t *testing.T) {
	cursor, err := mongo.NewCursorFromDocuments(nil, nil, nil)
	if err != nil {
		t.Fatalf("new cursor: %v", err)
	}

	total, err := TotalFromCursor(context.Background(), cursor)
	if err != nil {
		t.Fatalf("open power total: %v", err)
	}
	if total != 0 {
		t.Fatalf("expected total 0, got %d", total)
	}
}
