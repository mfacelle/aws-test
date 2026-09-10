package main

import (
	"context"
	"fmt"
	"image/color"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font/basicfont"
)

const (
	windowWidth           = 800
	windowHeight          = 600
	gridSize              = 21 // number of grid cells: [-10, 10] for x and y
	cellSize              = 24 // pixel size of cells in grid
	gridOriginX           = (windowWidth - gridSize*cellSize) / 2
	gridOriginY           = 45
	playerSize            = 18
	playerRefreshInterval = time.Second
)

type PlayerGame struct {
	ctx                     context.Context
	repository              *PlayerRepository
	player                  Player
	otherPlayers            []Player
	lastOtherPlayersRefresh time.Time
}

func NewPlayerGame(ctx context.Context, repository *PlayerRepository, player Player, otherPlayers []Player) *PlayerGame {
	return &PlayerGame{ctx: ctx, repository: repository, player: player, otherPlayers: otherPlayers}
}

// move or damage active player
func (game *PlayerGame) Update() error {
	if time.Since(game.lastOtherPlayersRefresh) >= playerRefreshInterval {
		game.lastOtherPlayersRefresh = time.Now()
		otherPlayers, err := game.repository.GetOtherPlayers(game.ctx, game.player.ID)
		if err != nil {
			log.Printf("refresh other players: %v", err)
		} else {
			game.otherPlayers = otherPlayers
		}
	}

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

	// draw gridlines
	gridColor := color.RGBA{R: 65, G: 76, B: 96, A: 255}
	// vertical
	for column := 0; column <= gridSize; column++ {
		x := gridOriginX + column*cellSize
		vector.FillRect(screen, float32(x), float32(gridOriginY), 1, float32(gridSize*cellSize), gridColor, false)
	}
	// horizontal
	for row := 0; row <= gridSize; row++ {
		y := gridOriginY + row*cellSize
		vector.FillRect(screen, float32(gridOriginX), float32(y), float32(gridSize*cellSize), 1, gridColor, false)
	}

	// draw player
	playerX := gridOriginX + (game.player.X+gridSize/2)*cellSize + (cellSize-playerSize)/2
	playerY := gridOriginY + (game.player.Y+gridSize/2)*cellSize + (cellSize-playerSize)/2
	vector.FillRect(screen, float32(playerX), float32(playerY), float32(playerSize), float32(playerSize), color.RGBA{R: 80, G: 210, B: 130, A: 255}, false)

	// draw debug text under player
	debug := fmt.Sprintf("%s\nhealth=%d\npos=(%d, %d)", game.player.ID, game.player.Health, game.player.X, game.player.Y)
	text.Draw(screen, debug, basicfont.Face7x13, playerX+playerSize, playerY+playerSize, color.White)

	// draw other players
	otherPlayerColor := color.RGBA{R: 235, G: 165, B: 75, A: 255}
	for _, player := range game.otherPlayers {
		if player.X < -gridSize/2 || player.X > gridSize/2 || player.Y < -gridSize/2 || player.Y > gridSize/2 {
			continue
		}
		otherX := gridOriginX + (player.X+gridSize/2)*cellSize + (cellSize-playerSize)/2
		otherY := gridOriginY + (player.Y+gridSize/2)*cellSize + (cellSize-playerSize)/2
		vector.FillRect(screen, float32(otherX), float32(otherY), float32(playerSize), float32(playerSize), otherPlayerColor, false)

		// draw debug text under player
		debug := fmt.Sprintf("%s\nhealth=%d\npos=(%d, %d)", player.ID, player.Health, player.X, player.Y)
		text.Draw(screen, debug, basicfont.Face7x13, otherX+playerSize, otherY+playerSize, otherPlayerColor)

	}

	// draw instructions at top of screen
	text.Draw(screen, "WASD: move    Z: damage 10", basicfont.Face7x13, 0, 10, color.RGBA{R: 180, G: 190, B: 205, A: 255})
}

func (game *PlayerGame) Layout(_, _ int) (int, int) {
	return windowWidth, windowHeight
}
