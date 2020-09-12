package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/codescot/go-common/httputil"
)

var igdbUserKey string

type gameinfo struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Summary string `json:"summary"`
}

type gamelist []gameinfo

type steamuri struct {
	ID  int    `json:"id"`
	URL string `json:"url"`
}

type steamlist []steamuri

func getCurrentGame(channel string) string {
	req := httputil.HTTP{
		TargetURL: fmt.Sprintf("https://decapi.me/twitch/game/%s", channel),
	}

	game, _ := req.String()

	return game
}

func ignoreGame(game string) bool {
	switch game {
	case "Games + Demos":
		return true
	default:
		return false
	}
}

func getGameInfo(game string) (gameinfo, error) {
	req := httputil.HTTP{
		TargetURL: "https://api-v3.igdb.com/games",
		Headers: map[string]string{
			"user-key": igdbUserKey,
		},
		Body: fmt.Sprintf("fields id,name,summary; search \"%s\"; where name != \"Fall Guy\"; limit 3;", game),
	}

	var games gamelist

	req.JSON(&games)

	if len(games) == 0 {
		return gameinfo{}, fmt.Errorf("no game info for %s", game)
	}

	return games[0], nil
}

func getSteamURI(gameid int) (steamuri, error) {
	req := httputil.HTTP{
		TargetURL: "https://api-v3.igdb.com/websites",
		Headers: map[string]string{
			"user-key": igdbUserKey,
		},
		Body: fmt.Sprintf("fields url; where category=13 & game=%d;", gameid),
	}

	var uris steamlist

	req.JSON(&uris)

	if len(uris) == 0 {
		return steamuri{}, errors.New("no steam uri")
	}

	return uris[0], nil
}

func okString(response string) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       response,
	}, nil
}

func truncateString(s string, m int) string {
	var b bytes.Buffer
	b.WriteString(s[0:m])
	b.WriteString("...")
	return b.String()
}

func handleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	igdbUserKey = os.Getenv("IGDB_KEY")

	channel := request.QueryStringParameters["channel"]
	currentGame := getCurrentGame(channel)
	if ignoreGame(currentGame) {
		return okString(currentGame)
	}

	gameInfo, err := getGameInfo(currentGame)
	if err != nil {
		return okString(err.Error())
	}

	steamURI, err := getSteamURI(gameInfo.ID)
	if err != nil {
		summary := gameInfo.Summary
		if len(summary) > 500 {
			summary = truncateString(summary, 497)
		}

		return okString(summary)
	}

	summary := gameInfo.Summary
	steam := steamURI.URL
	if len(summary) > 500 {
		m := 496 - len(steam)
		summary = truncateString(summary, m)
	}

	return okString(fmt.Sprintf("%s %s", summary, steam))
}

func main() {
	lambda.Start(handleRequest)
}
