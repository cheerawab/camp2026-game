package auth

import (
	"github.com/sitcon-tw/camp2026-game/internal/http/apimodel"
	"github.com/sitcon-tw/camp2026-game/internal/http/httpx"
)

var (
	_ = apimodel.AuthLoginRequest{}
	_ = apimodel.AuthLoginResponse{}
	_ = httpx.ProblemDetails{}
)
