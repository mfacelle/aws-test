package awstestcommon

import (
	"context"
	"fmt"
	"image/color"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font/basicfont"
)

const (
	WindowWidth           = 800
	WindowHeight          = 600
	GridSize              = 21 // number of grid cells: [-10, 10] for x and y
	CellSize              = 24 // pixel size of cells in grid
	GridOriginX           = (WindowWidth - GridSize*CellSize) / 2
	GridOriginY           = 45
	PlayerSize            = 18
	NpcSize               = 15
	PlayerRefreshInterval = time.Second
)

// need to export Ctx and repos so derived classes can use them.
// should ideally use getters or something, but going the lazy "barely change anything" route for now
type AwsTestGame struct {
	Ctx                context.Context
	PlayerRepo         *PlayerRepository
	players            []Player
	NpcRepo            *PlayerRepository
	Npcs               []Player
	lastplayersRefresh time.Time
}

func NewAwsTestGame(ctx context.Context, playerRepo *PlayerRepository, players []Player, npcRepo *PlayerRepository, npcs []Player) *AwsTestGame {
	return &AwsTestGame{Ctx: ctx, PlayerRepo: playerRepo, players: players, NpcRepo: npcRepo, Npcs: npcs}
}

// update all players from database
func (game *AwsTestGame) Update() error {
	// TODO do this every Update?
	if time.Since(game.lastplayersRefresh) >= PlayerRefreshInterval {
		game.lastplayersRefresh = time.Now()

		players, err := game.PlayerRepo.GetAllPlayers(game.Ctx)
		if err != nil {
			log.Printf("refresh all players: %v", err)
		} else {
			game.players = players
		}

		npcs, err := game.NpcRepo.GetAllPlayers(game.Ctx)
		if err != nil {
			log.Printf("refresh all NPCs: %v", err)
		} else {
			game.Npcs = npcs
		}
	}

	return nil
}

// draws a basic grid and player positions, with some debug text
func (game *AwsTestGame) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{R: 20, G: 20, B: 40, A: 255})

	// draw gridlines
	gridColor := color.RGBA{R: 40, G: 40, B: 80, A: 255}
	// vertical
	for column := 0; column <= GridSize; column++ {
		x := GridOriginX + column*CellSize
		vector.FillRect(screen, float32(x), float32(GridOriginY), 1, float32(GridSize*CellSize), gridColor, false)
	}
	// horizontal
	for row := 0; row <= GridSize; row++ {
		y := GridOriginY + row*CellSize
		vector.FillRect(screen, float32(GridOriginX), float32(y), float32(GridSize*CellSize), 1, gridColor, false)
	}

	// draw NPCs
	npcColor := color.RGBA{R: 100, G: 100, B: 200, A: 255}
	for _, npc := range game.Npcs {
		if npc.X < -GridSize/2 || npc.X > GridSize/2 || npc.Y < -GridSize/2 || npc.Y > GridSize/2 {
			continue
		}
		otherX := GridOriginX + (npc.X+GridSize/2)*CellSize + (CellSize-NpcSize)/2
		otherY := GridOriginY + (npc.Y+GridSize/2)*CellSize + (CellSize-NpcSize)/2
		vector.FillRect(screen, float32(otherX), float32(otherY), float32(NpcSize), float32(NpcSize), npcColor, false)

		// draw debug text under npc
		debug := fmt.Sprintf("%s\nhealth=%d\npos=(%d, %d)", npc.ID, npc.Health, npc.X, npc.Y)
		text.Draw(screen, debug, basicfont.Face7x13, otherX+NpcSize, otherY+NpcSize, npcColor)
	}

	// draw players (after NPCs so they show up on top)
	playerColor := color.RGBA{R: 200, G: 150, B: 50, A: 255}
	for _, player := range game.players {
		if player.X < -GridSize/2 || player.X > GridSize/2 || player.Y < -GridSize/2 || player.Y > GridSize/2 {
			continue
		}
		otherX := GridOriginX + (player.X+GridSize/2)*CellSize + (CellSize-PlayerSize)/2
		otherY := GridOriginY + (player.Y+GridSize/2)*CellSize + (CellSize-PlayerSize)/2
		vector.FillRect(screen, float32(otherX), float32(otherY), float32(PlayerSize), float32(PlayerSize), playerColor, false)

		// draw debug text under player
		debug := fmt.Sprintf("%s\nhealth=%d\npos=(%d, %d)", player.ID, player.Health, player.X, player.Y)
		text.Draw(screen, debug, basicfont.Face7x13, otherX+PlayerSize, otherY+PlayerSize, playerColor)
	}
}

func (game *AwsTestGame) Layout(_, _ int) (int, int) {
	return WindowWidth, WindowHeight
}
