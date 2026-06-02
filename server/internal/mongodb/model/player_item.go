package model

const PlayerItemsCollection = "player_items"

type PlayerItem struct {
	ID       string `bson:"_id"`
	PlayerID string `bson:"player_id"`
	ItemID   string `bson:"item_id"`
	Quantity int    `bson:"quantity"`
}
