package openpower

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	mongomodel "github.com/sitcon-tw/camp2026-game/internal/mongodb/model"
)

func TotalForPlayer(ctx context.Context, db *mongo.Database, playerID string) (int, error) {
	cursor, err := db.Collection(mongomodel.OpenPowerRecordsCollection).Aggregate(ctx, TotalPipeline(playerID))
	if err != nil {
		return 0, err
	}
	defer func() {
		_ = cursor.Close(ctx)
	}()

	return TotalFromCursor(ctx, cursor)
}

func TotalFromCursor(ctx context.Context, cursor *mongo.Cursor) (int, error) {
	var totals []struct {
		Total int `bson:"total"`
	}
	if err := cursor.All(ctx, &totals); err != nil {
		return 0, err
	}
	if len(totals) == 0 {
		return 0, nil
	}
	return totals[0].Total, nil
}

func TotalPipeline(playerID string) mongo.Pipeline {
	return mongo.Pipeline{
		{{Key: "$match", Value: bson.D{{Key: "player_id", Value: playerID}}}},
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: nil},
			{Key: "total", Value: bson.D{{Key: "$sum", Value: "$amount"}}},
		}}},
	}
}
