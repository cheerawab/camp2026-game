package model

const PlayerSitonesCollection = "player_sitones"

type PlayerSitone struct {
	ID       string `bson:"_id"`
	PlayerID string `bson:"player_id"`
	SitoneID string `bson:"sitone_id"`
	Quantity int    `bson:"quantity"`
}
