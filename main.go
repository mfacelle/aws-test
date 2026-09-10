package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func main() {
	// 1. Load AWS configuration (automatically picks up EC2 Instance Profile roles)
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-2"))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	// 2. Create DynamoDB Client
	dbClient := dynamodb.NewFromConfig(cfg)

	// 3. Put an item into the "Users" table
	tableName := "Users"
	input := &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item: map[string]types.AttributeValue{
			"UserID":   &types.AttributeValueMemberS{Value: "user_12345"},
			"UserName": &types.AttributeValueMemberS{Value: "Alice Dev"},
		},
	}

	_, err = dbClient.PutItem(context.TODO(), input)
	if err != nil {
		log.Fatalf("Failed to put item into DynamoDB: %v", err)
	}

	fmt.Println("Successfully added 'Alice Dev' to the Users table!")
}
