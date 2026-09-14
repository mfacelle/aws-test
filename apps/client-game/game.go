package main

import (
	common "aws-test/common"
	"context"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font/basicfont"
)

type PlayerGame struct {
	common.AwsTestGame
	player common.Player
}

func NewPlayerGame(ctx context.Context,
	playerRepo *common.PlayerRepository,
	player common.Player,
	allPlayers []common.Player,
	npcRepo *common.PlayerRepository,
	allNpcs []common.Player) *PlayerGame {
	return &PlayerGame{
		AwsTestGame: *common.NewAwsTestGame(ctx, playerRepo, allPlayers, npcRepo, allNpcs),
		player:      player,
	}
}

// move or damage active player
func (game *PlayerGame) Update() error {
	// call base class logic (updates all players from database)
	game.AwsTestGame.Update()

	// note: using arrow keys appears to crash the application, even with nothing in Update.
	// not sure exactly why yet but, since this is just for learning... it doesn't matter. just use WASD
	newX, newY := game.player.X, game.player.Y
	if inpututil.IsKeyJustPressed(ebiten.KeyA) {
		newX--
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyD) {
		newX++
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyW) {
		newY--
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyS) {
		newY++
	}
	if newX != game.player.X || newY != game.player.Y {
		if newX >= -common.GridSize/2 &&
			newX <= common.GridSize/2 &&
			newY >= -common.GridSize/2 &&
			newY <= common.GridSize/2 {
			if err := game.AwsTestGame.PlayerRepo.MovePlayer(game.AwsTestGame.Ctx, game.player.ID, newX, newY); err != nil {
				log.Printf("move player: %v", err)
				return nil
			}
			// TODO don't actually want to update player object position here... really just want to update the db!
			// leaving for now, because it's interesting to see this way
			game.player.X, game.player.Y = newX, newY
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyZ) && game.player.Health > 0 {
		if err := game.AwsTestGame.PlayerRepo.DamagePlayer(game.AwsTestGame.Ctx, game.player.ID, 10); err != nil {
			log.Printf("damage player: %v", err)
			return nil
		}
		game.player.Health -= 10
		if game.player.Health < 0 {
			game.player.Health = 0
		}
	}
	return nil
}

// draws a basic grid and player position, with some debug text
func (game *PlayerGame) Draw(screen *ebiten.Image) {

	// base class handles all main rendering
	game.AwsTestGame.Draw(screen)

	// draw player (essentially just drawing another rect on top of base game rendering)
	playerX := common.GridOriginX + (game.player.X+common.GridSize/2)*common.CellSize + (common.CellSize-common.PlayerSize)/2
	playerY := common.GridOriginY + (game.player.Y+common.GridSize/2)*common.CellSize + (common.CellSize-common.PlayerSize)/2
	vector.FillRect(screen, float32(playerX), float32(playerY), float32(common.PlayerSize), float32(common.PlayerSize), color.RGBA{R: 150, G: 250, B: 200, A: 255}, false)

	// draw instructions at top of screen
	text.Draw(screen, "WASD: move    Z: damage 10", basicfont.Face7x13, 0, 10, color.RGBA{R: 150, G: 150, B: 200, A: 255})
}
