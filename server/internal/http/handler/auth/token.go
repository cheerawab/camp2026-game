package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"

	mongomodel "github.com/sitcon-tw/camp2026-game/internal/mongodb/model"
)

const authTokenBytes = 32

var errAuthTokenStale = errors.New("auth token is no longer current")

func newAuthToken() (string, error) {
	randomBytes := make([]byte, authTokenBytes)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", fmt.Errorf("generate auth token: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(randomBytes), nil
}

func (h *Handler) rotateAuthToken(ctx context.Context, playerID string, currentToken string) (string, error) {
	nextToken, err := newAuthToken()
	if err != nil {
		return "", err
	}

	result, err := h.db.Collection(mongomodel.PlayersCollection).UpdateOne(
		ctx,
		authTokenRotationFilter(playerID, currentToken),
		authTokenRotationUpdate(nextToken),
	)
	if err != nil {
		return "", err
	}
	if result.MatchedCount == 0 {
		return "", errAuthTokenStale
	}
	return nextToken, nil
}

func authTokenRotationFilter(playerID string, currentToken string) bson.M {
	return bson.M{
		"_id":        playerID,
		"auth_token": currentToken,
	}
}

func authTokenRotationUpdate(nextToken string) bson.M {
	return bson.M{
		"$set": bson.M{
			"auth_token": nextToken,
		},
	}
}
