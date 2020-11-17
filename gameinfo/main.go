package main

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/codescot/gameinfo/decapi"
	"github.com/codescot/gameinfo/igdb"
	"github.com/codescot/gameinfo/twitch"
)

func shortCircuit(game string) (int, bool) {
	switch game {
	case "Games + Demos":
		return 0, true
	case "Just Chatting":
		return 0, true
	case "Among Us":
		return 111469, true
	case "Unravel":
		return 11170, true
	case "Cities: Skylines":
		return 9066, true
	default:
		return 0, false
	}
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
	twitchClientID := os.Getenv("TWITCH_CLIENT_ID")
	twitchClientSecret := os.Getenv("TWITCH_CLIENT_SECRET")

	tw := twitch.Twitch{
		ClientID:     twitchClientID,
		ClientSecret: twitchClientSecret,
	}

	if err := tw.GetCredentials(); err != nil {
		panic(err.Error)
	}

	ig := igdb.IGDB{
		ClientID:      twitchClientID,
		Authorization: fmt.Sprintf("Bearer %s", tw.AccessToken),
	}

	channel := request.QueryStringParameters["channel"]
	currentGame := decapi.GameCategoryForChannel(channel)

	var gameInfo igdb.GameInfo
	var err error
	if id, ok := shortCircuit(currentGame); ok {
		if id == 0 {
			return okString(currentGame)
		}

		gameInfo, err = ig.GetGameInfoByID(id)
		if err != nil {
			return okString(err.Error())
		}
	} else {
		gameInfo, err = ig.GetGameInfo(currentGame)
		if err != nil {
			return okString(err.Error())
		}
	}

	summary := gameInfo.Summary
	summary = strings.ReplaceAll(summary, "\r", "")
	summary = strings.ReplaceAll(summary, "\n", "")
	summary = strings.ReplaceAll(summary, "  ", " ")

	url := gameInfo.URL
	if len(summary) > 400 {
		m := 396 - len(url)
		summary = truncateString(summary, m)
	}

	return okString(fmt.Sprintf("%s %s", summary, url))
}

func main() {
	lambda.Start(handleRequest)
}
