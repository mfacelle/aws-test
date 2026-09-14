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

// launch client-side application, where the player can control their specific character
func main() {
	if len(os.Args) != 2 || os.Args[1] == "" {
		log.Fatalf("usage: %s PLAYER_ID", os.Args[0])
	}

	ctx := context.Background()
	region := getenv("AWS_REGION", "us-east-2")
	tableName := getenv("PLAYERS_TABLE", "Players")
	npcsTableName := getenv("NPC_TABLE", "NPCs")
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		log.Fatalf("unable to load AWS config: %v", err)
	}

	playerRepo := common.NewPlayerRepository(dynamodb.NewFromConfig(cfg), tableName)
	player, err := playerRepo.GetPlayer(ctx, os.Args[1])
	if err != nil {
		log.Fatalf("load player: %v", err)
	}
	allPlayers, err := playerRepo.GetAllPlayers(ctx)
	if err != nil {
		log.Fatalf("load all players: %v", err)
	}

	npcRepo := common.NewPlayerRepository(dynamodb.NewFromConfig(cfg), npcsTableName)
	allNpcs, err := npcRepo.GetAllPlayers(ctx)
	if err != nil {
		log.Fatalf("load NPCs: %v", err)
	}

	game := NewPlayerGame(ctx, playerRepo, player, allPlayers, npcRepo, allNpcs)
	ebiten.SetWindowSize(common.WindowWidth, common.WindowHeight)
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
