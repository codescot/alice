package igdb

import (
	"fmt"

	"github.com/codescot/go-common/httputil"
)

var igdbUserKey string

// IGDB internet game database
type IGDB struct {
	ClientID      string
	Authorization string
}

// GameInfo result from IGDB
type GameInfo struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Summary   string `json:"summary"`
	Platforms []int  `json:"platforms"`
	URL       string `json:"url"`
}

type gamelist []GameInfo

const baseURL = "https://api.igdb.com/v4"

// GetGameInfo get game info for specified game name
func (ig IGDB) GetGameInfo(game string) (GameInfo, error) {
	req := httputil.HTTP{
		TargetURL: fmt.Sprintf("%s/games", baseURL),
		Method:    "POST",
		Headers: map[string]string{
			"Client-ID":     ig.ClientID,
			"Authorization": ig.Authorization,
		},
		Body: fmt.Sprintf("fields id,name,summary,url; search \"%s\"; where release_dates.platform = (48,49,6,8,9); limit 1;", game),
	}

	var games gamelist

	req.JSON(&games)

	if len(games) == 0 {
		return GameInfo{}, fmt.Errorf("no game info for %s", game)
	}

	return games[0], nil
}

// GetGameInfoByID get game info for specified game id
func (ig IGDB) GetGameInfoByID(gameID int) (GameInfo, error) {
	req := httputil.HTTP{
		TargetURL: fmt.Sprintf("%s/games", baseURL),
		Method:    "POST",
		Headers: map[string]string{
			"Client-ID":     ig.ClientID,
			"Authorization": ig.Authorization,
		},
		Body: fmt.Sprintf("fields id,name,summary,url; where id = %d & release_dates.platform = (48,49,6,8,9); limit 1;", gameID),
	}

	var games gamelist

	req.JSON(&games)

	if len(games) == 0 {
		return GameInfo{}, fmt.Errorf("no game info for %d", gameID)
	}

	return games[0], nil
}
