package main

import (
	common "aws-test/common"
	"context"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	NpcMoveInterval = time.Second * 2
)

type MovePair struct {
	X int
	Y int
}

type ServerGame struct {
	common.AwsTestGame
	lastNpcUpdate   time.Time   // how frequently npcs will move
	randomMoveTable [5]MovePair // random moves an npc can make
}

func NewServerGame(ctx context.Context,
	playerRepo *common.PlayerRepository,
	allPlayers []common.Player,
	npcRepo *common.PlayerRepository,
	allNpcs []common.Player) *ServerGame {
	return &ServerGame{
		AwsTestGame: *common.NewAwsTestGame(ctx, playerRepo, allPlayers, npcRepo, allNpcs),
		randomMoveTable: [5]MovePair{
			{X: 0, Y: 0},
			{X: 0, Y: -1},
			{X: 0, Y: 1},
			{X: -1, Y: 0},
			{X: 1, Y: 0},
		},
	}
}

// move or damage active player
// TODO make database pull a function in common base game
func (game *ServerGame) Update() error {
	// call base class logic (updates all players from database)
	game.AwsTestGame.Update()

	// randomly pick a direction and move NPCs
	if time.Since(game.lastNpcUpdate) >= NpcMoveInterval {
		game.lastNpcUpdate = time.Now()

		for _, npc := range game.AwsTestGame.Npcs {
			// select a random move
			randomMoveIndex := rand.Intn(len(game.randomMoveTable))
			movement := game.randomMoveTable[randomMoveIndex]

			if movement.X != 0 || movement.Y != 0 {
				// apply movement and check if within grid bounds
				// TODO need to also check against player positions.  Allows collisions, for now.
				// implement a lookup based on position, for anything occupying that space?
				newX := npc.X + movement.X
				newY := npc.Y + movement.Y
				if newX >= -common.GridSize/2 &&
					newX <= common.GridSize/2 &&
					newY >= -common.GridSize/2 &&
					newY <= common.GridSize/2 {
					if err := game.AwsTestGame.NpcRepo.MovePlayer(game.AwsTestGame.Ctx, npc.ID, newX, newY); err != nil {
						log.Printf("move NPC: %v", err)
						return nil
					}
				}
			}
		}
	}

	return nil
}

// draws a basic grid and player position, with some debug text
// TODO make rendering a function in common base game
func (game *ServerGame) Draw(screen *ebiten.Image) {

	// base class handles all main rendering
	game.AwsTestGame.Draw(screen)
}
