package model

import "time"

const MatchItemDropsCollection = "match_item_drops"

type MatchItemDrop struct {
	ID        string    `bson:"_id"`
	MatchID   string    `bson:"match_id"`
	PlayerID  string    `bson:"player_id"`
	ItemID    string    `bson:"item_id,omitempty"`
	Quantity  int       `bson:"quantity,omitempty"`
	DropRate  int       `bson:"drop_rate"`
	Dropped   bool      `bson:"dropped"`
	Granted   bool      `bson:"granted,omitempty"`
	Source    string    `bson:"source"`
	CreatedAt time.Time `bson:"created_at"`
}
