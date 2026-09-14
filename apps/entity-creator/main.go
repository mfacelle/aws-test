package main

import (
	common "aws-test/common"
	"context"
	"flag"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

// launch quick app that creates a new player or npc in the database
func main() {
	flags := flag.NewFlagSet("entity-creator", flag.ExitOnError)
	npcID := flags.String("id", "", "unique NPC ID")
	health := flags.Int("health", 100, "starting NPC health")
	x := flags.Int("x", 0, "starting X position")
	y := flags.Int("y", 0, "starting Y position")
	flags.Parse(os.Args[1:])

	if *npcID == "" {
		flags.Usage()
		os.Exit(2)
	}

	// TODO probably want to add an arg to allows this app to create a player OR an NPC.
	// but this is good enough for now.  Can always just update the db manually

	ctx := context.Background()
	region := getenv("AWS_REGION", "us-east-2")
	tableName := getenv("NPC_TABLE", "NPCs")
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		log.Fatalf("unable to load AWS config: %v", err)
	}

	repository := common.NewPlayerRepository(dynamodb.NewFromConfig(cfg), tableName)
	if err := repository.CreatePlayer(ctx, common.Player{
		ID:     *npcID,
		Health: *health,
		X:      *x,
		Y:      *y,
	}); err != nil {
		log.Fatalf("create NPC: %v", err)
	}

	log.Printf("created NPC %q at (%d, %d) with health %d in table %q", *npcID, *x, *y, *health, tableName)
}

func getenv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
