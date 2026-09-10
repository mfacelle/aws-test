package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type User struct {
	ID   string
	Name string
}

type dynamoDB interface {
	PutItem(context.Context, *dynamodb.PutItemInput, ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
}

// UserRepository persists users in DynamoDB.
type UserRepository struct {
	db        dynamoDB
	tableName string
}

func NewUserRepository(db dynamoDB, tableName string) *UserRepository {
	return &UserRepository{db: db, tableName: tableName}
}

func (repository *UserRepository) AddUser(ctx context.Context, user User) error {
	user.ID = strings.TrimSpace(user.ID)
	user.Name = strings.TrimSpace(user.Name)
	if user.ID == "" {
		return fmt.Errorf("user ID is required")
	}
	if user.Name == "" {
		return fmt.Errorf("user name is required")
	}

	_, err := repository.db.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(repository.tableName),
		Item: map[string]types.AttributeValue{
			"UserID":   &types.AttributeValueMemberS{Value: user.ID},
			"UserName": &types.AttributeValueMemberS{Value: user.Name},
		},
		ConditionExpression: aws.String("attribute_not_exists(UserID)"),
	})
	if err != nil {
		return fmt.Errorf("put user %q: %w", user.ID, err)
	}
	return nil
}
