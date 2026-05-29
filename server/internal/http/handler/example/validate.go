package example

import "strings"

func normalizeValidationExampleRequest(body *ValidationExampleRequest) {
	for i := range body.Players {
		body.Players[i].DisplayName = strings.TrimSpace(body.Players[i].DisplayName)
		body.Players[i].FavoritePebbleType = strings.TrimSpace(body.Players[i].FavoritePebbleType)
	}
}
