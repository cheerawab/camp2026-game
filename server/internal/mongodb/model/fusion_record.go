package model

import "time"

const FusionRecordsCollection = "fusion_records"

type FusionRecord struct {
	ID        string            `bson:"_id"`
	PlayerID  string            `bson:"player_id"`
	RecipeID  string            `bson:"recipe_id"`
	Inputs    []FusionComponent `bson:"inputs"`
	Outputs   []FusionComponent `bson:"outputs"`
	CreatedAt time.Time         `bson:"created_at"`
}

type FusionComponent struct {
	Kind     string `bson:"kind"`
	RefID    string `bson:"ref_id"`
	Quantity int    `bson:"quantity"`
}
