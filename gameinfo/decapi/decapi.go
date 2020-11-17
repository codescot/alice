package decapi

import (
	"fmt"

	"github.com/codescot/go-common/httputil"
)

// GameCategoryForChannel game for the current twitch channel
func GameCategoryForChannel(channel string) string {
	req := httputil.HTTP{
		TargetURL: fmt.Sprintf("https://decapi.me/twitch/game/%s", channel),
	}

	game, _ := req.String()

	return game
}
