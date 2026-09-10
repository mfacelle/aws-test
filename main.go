package main

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	if len(os.Args) != 2 || os.Args[1] == "" {
		log.Fatalf("usage: %s PLAYER_ID", os.Args[0])
	}

	ctx := context.Background()
	region := getenv("AWS_REGION", "us-east-2")
	tableName := getenv("PLAYERS_TABLE", "Players")
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		log.Fatalf("unable to load AWS config: %v", err)
	}

	repository := NewPlayerRepository(dynamodb.NewFromConfig(cfg), tableName)
	player, err := repository.GetPlayer(ctx, os.Args[1])
	if err != nil {
		log.Fatalf("load player: %v", err)
	}
	otherPlayers, err := repository.GetOtherPlayers(ctx, player.ID)
	if err != nil {
		log.Fatalf("load other players: %v", err)
	}

	game := NewPlayerGame(ctx, repository, player, otherPlayers)
	ebiten.SetWindowSize(windowWidth, windowHeight)
	ebiten.SetWindowTitle("Player Viewer - " + player.ID)
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

func getenv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
