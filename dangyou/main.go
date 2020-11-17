package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
)

type event struct{}

var thingsToDo = []string{
	"alert length",
	"mosquito",
	"salmonberry",
	"autosave",
	"Twitch",
	"washing machine",
	"sound clip",
	"children",
	"Logitech",
	"focaccia",
	"pizza crust",
	"Alicia",
}

func nextInt(n int) int {
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	return r.Intn(n)
}

func nextString() string {
	var l = len(thingsToDo)
	var i = nextInt(l)
	return thingsToDo[i]
}

func handleRequest(ctx context.Context, name event) (string, error) {
	return fmt.Sprintf("Dang you, %s! hekiFist", nextString()), nil
}

func main() {
	lambda.Start(handleRequest)
}
