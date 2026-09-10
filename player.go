package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// this class has player object representations as well as functions to handle
// database operations for players

// represents player information
type Player struct {
	ID     string
	Health int
	X      int
	Y      int
}

type dynamoDB interface {
	PutItem(context.Context, *dynamodb.PutItemInput, ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
	UpdateItem(context.Context, *dynamodb.UpdateItemInput, ...func(*dynamodb.Options)) (*dynamodb.UpdateItemOutput, error)
	GetItem(context.Context, *dynamodb.GetItemInput, ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error)
	Scan(context.Context, *dynamodb.ScanInput, ...func(*dynamodb.Options)) (*dynamodb.ScanOutput, error)
}

type PlayerRepository struct {
	db        dynamoDB
	tableName string
}

func NewPlayerRepository(db dynamoDB, tableName string) *PlayerRepository {
	return &PlayerRepository{db: db, tableName: tableName}
}

func (repository *PlayerRepository) CreatePlayer(ctx context.Context, player Player) error {
	player.ID = strings.TrimSpace(player.ID)
	if player.ID == "" {
		return fmt.Errorf("player ID is required")
	}
	if player.Health < 0 {
		return fmt.Errorf("health cannot be negative")
	}

	_, err := repository.db.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(repository.tableName),
		Item: map[string]types.AttributeValue{
			"PlayerID":  &types.AttributeValueMemberS{Value: player.ID},
			"Health":    &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", player.Health)},
			"PositionX": &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", player.X)},
			"PositionY": &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", player.Y)},
		},
		ConditionExpression: aws.String("attribute_not_exists(PlayerID)"),
	})
	if err != nil {
		return fmt.Errorf("create player %q: %w", player.ID, err)
	}
	return nil
}

func (repository *PlayerRepository) MovePlayer(ctx context.Context, id string, x, y int) error {
	return repository.update(ctx, id, "SET PositionX = :x, PositionY = :y", map[string]types.AttributeValue{
		":x": &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", x)},
		":y": &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", y)},
	})
}

func (repository *PlayerRepository) DamagePlayer(ctx context.Context, id string, damage int) error {
	if damage <= 0 {
		return fmt.Errorf("damage amount must be greater than zero")
	}
	return repository.update(ctx, id, "SET Health = Health - :damage", map[string]types.AttributeValue{
		":damage": &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", damage)},
	}, "Health >= :damage")
}

func (repository *PlayerRepository) update(ctx context.Context, id, expression string, values map[string]types.AttributeValue, conditions ...string) error {
	id = strings.TrimSpace(id)
	if id == "" {
		return fmt.Errorf("player ID is required")
	}
	input := &dynamodb.UpdateItemInput{
		TableName:                 aws.String(repository.tableName),
		Key:                       map[string]types.AttributeValue{"PlayerID": &types.AttributeValueMemberS{Value: id}},
		UpdateExpression:          aws.String(expression),
		ExpressionAttributeValues: values,
		ConditionExpression:       aws.String("attribute_exists(PlayerID)"),
	}
	if len(conditions) > 0 {
		input.ConditionExpression = aws.String("attribute_exists(PlayerID) AND " + conditions[0])
	}
	if _, err := repository.db.UpdateItem(ctx, input); err != nil {
		return fmt.Errorf("update player %q: %w", id, err)
	}
	return nil
}

func (repository *PlayerRepository) GetPlayer(ctx context.Context, id string) (Player, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return Player{}, fmt.Errorf("player ID is required")
	}
	result, err := repository.db.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(repository.tableName),
		Key:       map[string]types.AttributeValue{"PlayerID": &types.AttributeValueMemberS{Value: id}},
	})
	if err != nil {
		return Player{}, fmt.Errorf("read player %q: %w", id, err)
	}
	if len(result.Item) == 0 {
		return Player{}, fmt.Errorf("player %q was not found", id)
	}
	player, err := playerFromItem(result.Item)
	if err != nil {
		return Player{}, fmt.Errorf("decode player %q: %w", id, err)
	}
	return player, nil
}

func (repository *PlayerRepository) GetOtherPlayers(ctx context.Context, activeID string) ([]Player, error) {
	activeID = strings.TrimSpace(activeID)
	if activeID == "" {
		return nil, fmt.Errorf("player ID is required")
	}

	var players []Player
	var startKey map[string]types.AttributeValue
	for {
		result, err := repository.db.Scan(ctx, &dynamodb.ScanInput{
			TableName:         aws.String(repository.tableName),
			ExclusiveStartKey: startKey,
		})
		if err != nil {
			return nil, fmt.Errorf("scan players: %w", err)
		}

		for _, item := range result.Items {
			player, err := playerFromItem(item)
			if err != nil {
				return nil, fmt.Errorf("decode scanned player: %w", err)
			}
			// could probably just always get all players, and save off a pointer to the "current" player?
			// fine for now, but if I ever expand on this, might be worth doing
			if player.ID != activeID {
				players = append(players, player)
			}
		}

		if len(result.LastEvaluatedKey) == 0 {
			return players, nil
		}
		startKey = result.LastEvaluatedKey
	}
}

func playerFromItem(item map[string]types.AttributeValue) (Player, error) {
	idAttribute, ok := item["PlayerID"].(*types.AttributeValueMemberS)
	if !ok || strings.TrimSpace(idAttribute.Value) == "" {
		return Player{}, fmt.Errorf("missing PlayerID")
	}
	return Player{
		ID:     idAttribute.Value,
		Health: number(item["Health"]),
		X:      number(item["PositionX"]),
		Y:      number(item["PositionY"]),
	}, nil
}

func number(value types.AttributeValue) int {
	var parsed int
	if numeric, ok := value.(*types.AttributeValueMemberN); ok {
		_, _ = fmt.Sscanf(numeric.Value, "%d", &parsed)
	}
	return parsed
}
