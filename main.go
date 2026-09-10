package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func main() {
	if len(os.Args) < 2 || os.Args[1] != "add-user" {
		log.Fatalf("usage: %s add-user --id USER_ID --name USER_NAME", os.Args[0])
	}

	flags := flag.NewFlagSet("add-user", flag.ExitOnError)
	userID := flags.String("id", "", "unique user ID")
	userName := flags.String("name", "", "user name")
	_ = flags.Parse(os.Args[2:])
	if *userID == "" || *userName == "" {
		flags.Usage()
		os.Exit(2)
	}

	ctx := context.Background()
	region := getenv("AWS_REGION", "us-east-2")
	tableName := getenv("USERS_TABLE", "Users")
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		log.Fatalf("unable to load AWS config: %v", err)
	}

	repository := NewUserRepository(dynamodb.NewFromConfig(cfg), tableName)
	if err := repository.AddUser(ctx, User{ID: *userID, Name: *userName}); err != nil {
		log.Fatalf("add user: %v", err)
	}

	log.Printf("added user %q to DynamoDB table %q", *userID, tableName)
}

func getenv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
