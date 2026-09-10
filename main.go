package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("usage: %s <create-player|move-player|damage-player|get-player> [flags]", os.Args[0])
	}

	ctx := context.Background()
	region := getenv("AWS_REGION", "us-east-2")
	tableName := getenv("PLAYERS_TABLE", "Players")
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		log.Fatalf("unable to load AWS config: %v", err)
	}

	repository := NewPlayerRepository(dynamodb.NewFromConfig(cfg), tableName)
	if err := runAction(ctx, repository, os.Args[1], os.Args[2:]); err != nil {
		log.Fatal(err)
	}
}

func runAction(ctx context.Context, repository *PlayerRepository, action string, args []string) error {
	// TODO eventually just make this into a game loop that can call each of these as separate functions
	// maybe use command-line arg to set what player is being controlled?

	switch action {
	case "create-player":
		flags := flag.NewFlagSet(action, flag.ContinueOnError)
		id := flags.String("id", "", "unique player ID")
		health := flags.Int("health", 100, "starting health")
		x := flags.Int("x", 0, "starting X position")
		y := flags.Int("y", 0, "starting Y position")
		if err := flags.Parse(args); err != nil {
			return err
		}
		if err := repository.CreatePlayer(ctx, Player{ID: *id, Health: *health, X: *x, Y: *y}); err != nil {
			return fmt.Errorf("create player: %w", err)
		}
		log.Printf("created player %q", *id)
		return nil
	case "move-player":
		flags := flag.NewFlagSet(action, flag.ContinueOnError)
		id := flags.String("id", "", "player ID")
		x := flags.Int("x", 0, "new X position")
		y := flags.Int("y", 0, "new Y position")
		if err := flags.Parse(args); err != nil {
			return err
		}
		if err := repository.MovePlayer(ctx, *id, *x, *y); err != nil {
			return fmt.Errorf("move player: %w", err)
		}
		log.Printf("moved player %q to (%d, %d)", *id, *x, *y)
		return nil
	case "damage-player":
		flags := flag.NewFlagSet(action, flag.ContinueOnError)
		id := flags.String("id", "", "player ID")
		damage := flags.Int("amount", 0, "damage amount")
		if err := flags.Parse(args); err != nil {
			return err
		}
		if err := repository.DamagePlayer(ctx, *id, *damage); err != nil {
			return fmt.Errorf("damage player: %w", err)
		}
		log.Printf("damaged player %q by %d", *id, *damage)
		return nil
	case "get-player":
		flags := flag.NewFlagSet(action, flag.ContinueOnError)
		id := flags.String("id", "", "player ID")
		if err := flags.Parse(args); err != nil {
			return err
		}
		player, err := repository.GetPlayer(ctx, *id)
		if err != nil {
			return fmt.Errorf("get player: %w", err)
		}
		fmt.Printf("player=%s health=%d x=%d y=%d\n", player.ID, player.Health, player.X, player.Y)
		return nil
	default:
		return fmt.Errorf("unknown action %q", action)
	}
}

func getenv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
