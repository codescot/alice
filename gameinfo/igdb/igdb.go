package igdb

import (
	"errors"
	"fmt"

	"github.com/codescot/go-common/httputil"
)

var igdbUserKey string

// IGDB internet game database
type IGDB struct {
	Key string
}

// GameInfo result from IGDB
type GameInfo struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Summary string `json:"summary"`
}

type gamelist []GameInfo

// SteamURI result from IGDB
type SteamURI struct {
	ID  int    `json:"id"`
	URL string `json:"url"`
}

type steamlist []SteamURI

// GetGameInfo get game info for specified game name
func (ig IGDB) GetGameInfo(game string) (GameInfo, error) {
	req := httputil.HTTP{
		TargetURL: "https://api-v3.igdb.com/games",
		Headers: map[string]string{
			"user-key": ig.Key,
		},
		Body: fmt.Sprintf("fields id,name,summary; search \"%s\"; where name != \"Fall Guy\"; limit 3;", game),
	}

	var games gamelist

	req.JSON(&games)

	if len(games) == 0 {
		return GameInfo{}, fmt.Errorf("no game info for %s", game)
	}

	return games[0], nil
}

// GetSteamURI get the steam link for specified game id
func (ig IGDB) GetSteamURI(gameid int) (SteamURI, error) {
	req := httputil.HTTP{
		TargetURL: "https://api-v3.igdb.com/websites",
		Headers: map[string]string{
			"user-key": ig.Key,
		},
		Body: fmt.Sprintf("fields url; where category=13 & game=%d;", gameid),
	}

	var uris steamlist

	req.JSON(&uris)

	if len(uris) == 0 {
		return SteamURI{}, errors.New("no steam uri")
	}

	return uris[0], nil
}
