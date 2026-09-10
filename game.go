package main

import (
	"context"
	"fmt"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font/basicfont"
)

const (
	windowWidth  = 800
	windowHeight = 600
	gridSize     = 21
	cellSize     = 24
	gridOriginX  = (windowWidth - gridSize*cellSize) / 2
	gridOriginY  = 45
	playerSize   = 18
)

type PlayerGame struct {
	ctx        context.Context
	repository *PlayerRepository
	player     Player
}

func NewPlayerGame(ctx context.Context, repository *PlayerRepository, player Player) *PlayerGame {
	return &PlayerGame{ctx: ctx, repository: repository, player: player}
}

// move or damage active player
func (game *PlayerGame) Update() error {
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
		if newX >= -gridSize/2 && newX <= gridSize/2 && newY >= -gridSize/2 && newY <= gridSize/2 {
			if err := game.repository.MovePlayer(game.ctx, game.player.ID, newX, newY); err != nil {
				log.Printf("move player: %v", err)
				return nil
			}
			game.player.X, game.player.Y = newX, newY
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyZ) && game.player.Health > 0 {
		if err := game.repository.DamagePlayer(game.ctx, game.player.ID, 10); err != nil {
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
	screen.Fill(color.RGBA{R: 22, G: 27, B: 38, A: 255})

	gridColor := color.RGBA{R: 65, G: 76, B: 96, A: 255}
	for column := 0; column <= gridSize; column++ {
		x := gridOriginX + column*cellSize
		vector.FillRect(screen, float32(x), float32(gridOriginY), 1, float32(gridSize*cellSize), gridColor, false)
	}
	for row := 0; row <= gridSize; row++ {
		y := gridOriginY + row*cellSize
		vector.FillRect(screen, float32(gridOriginX), float32(y), float32(gridSize*cellSize), 1, gridColor, false)
	}

	playerX := gridOriginX + (game.player.X+gridSize/2)*cellSize + (cellSize-playerSize)/2
	playerY := gridOriginY + (game.player.Y+gridSize/2)*cellSize + (cellSize-playerSize)/2
	// log.Printf("drawing player at (%d, %d)", playerX, playerY)
	vector.FillRect(screen, float32(playerX), float32(playerY), float32(playerSize), float32(playerSize), color.RGBA{R: 80, G: 210, B: 130, A: 255}, false)

	debug := fmt.Sprintf("%s\nhealth=%d\npos=(%d, %d)", game.player.ID, game.player.Health, game.player.X, game.player.Y)
	text.Draw(screen, debug, basicfont.Face7x13, playerX+playerSize, playerY+playerSize, color.White)
	text.Draw(screen, "WASD: move    Z: damage 10", basicfont.Face7x13, 0, 10, color.RGBA{R: 180, G: 190, B: 205, A: 255})
}

func (game *PlayerGame) Layout(_, _ int) (int, int) {
	return windowWidth, windowHeight
}
