package main

import (
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
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Summary string `json:"summary"`
}

type gamelist []gameinfo

type steamuri struct {
	Id  int    `json:"id"`
	Url string `json:"url"`
}

type steamlist []steamuri

func getCurrentGame(channel string) string {
	req := httputil.HTTP{
		TargetURL: fmt.Sprintf("https://decapi.me/twitch/game/%s", channel),
		Headers: map[string]string{
			"user-key": igdbUserKey,
		},
	}

	game, _ := req.String()

	return game
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
		return gameinfo{}, errors.New(fmt.Sprintf("no game info for %s", game))
	}

	return games[0], nil
}

func getSteamUri(gameid int) (steamuri, error) {
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

func handleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	igdbUserKey = os.Getenv("IGDB_KEY")

	channel := request.QueryStringParameters["channel"]
	currentGame := getCurrentGame(channel)

	gameInfo, err := getGameInfo(currentGame)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Body:       err.Error(),
		}, nil
	}

	steamUri, err := getSteamUri(gameInfo.Id)
	if err != nil {
		summary := gameInfo.Summary
		if len(summary) > 500 {
			summary = summary[0:497] + "..."
		}

		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Body:       summary,
		}, nil
	}

	summary := gameInfo.Summary
	steam := steamUri.Url
	if len(summary) > 500 {
		maxLen := 496 - len(steam)
		summary = summary[0:maxLen] + "..."
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       fmt.Sprintf("%s %s", summary, steam),
	}, nil
}

func main() {
	lambda.Start(handleRequest)
}
