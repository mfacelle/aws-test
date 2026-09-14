package main

import (
	common "aws-test/common"
	"context"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/hajimehoshi/ebiten/v2"
)

// launch "server-side" version of the game, which (for now) randomly moves all NPCs around the grid.
// will really just be launched locally but, ideally, this would run on the AWS instance
func main() {

	ctx := context.Background()
	region := getenv("AWS_REGION", "us-east-2")
	playersTableName := getenv("PLAYERS_TABLE", "Players")
	npcsTableName := getenv("NPC_TABLE", "NPCs")
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		log.Fatalf("unable to load AWS config: %v", err)
	}

	playerRepo := common.NewPlayerRepository(dynamodb.NewFromConfig(cfg), playersTableName)
	allPlayers, err := playerRepo.GetAllPlayers(ctx)
	if err != nil {
		log.Fatalf("load players: %v", err)
	}

	npcRepo := common.NewPlayerRepository(dynamodb.NewFromConfig(cfg), npcsTableName)
	allNpcs, err := npcRepo.GetAllPlayers(ctx)
	if err != nil {
		log.Fatalf("load NPCs: %v", err)
	}

	game := NewServerGame(ctx, playerRepo, allPlayers, npcRepo, allNpcs)
	ebiten.SetWindowSize(common.WindowWidth, common.WindowHeight)
	ebiten.SetWindowTitle("Server Viewer")
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
