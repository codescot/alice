package main

import (
	"bytes"
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/codescot/gameinfo/decapi"
	"github.com/codescot/gameinfo/igdb"
)

func ignoreGame(game string) bool {
	switch game {
	case "Games + Demos":
		return true
	default:
		return false
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
	ig := igdb.IGDB{
		Key: os.Getenv("IGDB_KEY"),
	}

	channel := request.QueryStringParameters["channel"]
	currentGame := decapi.GameCategoryForChannel(channel)
	if ignoreGame(currentGame) {
		return okString(currentGame)
	}

	gameInfo, err := ig.GetGameInfo(currentGame)
	if err != nil {
		return okString(err.Error())
	}

	steamURI, err := ig.GetSteamURI(gameInfo.ID)
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
