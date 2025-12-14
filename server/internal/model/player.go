package model

import (
	"image/color"
	"math/rand"

	"github.com/google/uuid"
)

type Player struct {
	ID     uuid.UUID
	Pseudo string
	Game   *Game
	Color  color.Color
}

func getRandomColor() color.Color {
	idx := rand.Intn(10)
	colors := [10]color.Color{
		color.RGBA{R: 0, G: 255, B: 255, A: 255},
		color.RGBA{R: 255, G: 0, B: 255, A: 255},
		color.RGBA{R: 255, G: 255, B: 0, A: 255},
		color.RGBA{R: 255, G: 255, B: 100, A: 255},
		color.RGBA{R: 255, G: 100, B: 255, A: 255},
		color.RGBA{R: 100, G: 255, B: 255, A: 255},
		color.RGBA{R: 8, G: 255, B: 255, A: 255},
		color.RGBA{R: 255, G: 8, B: 255, A: 255},
		color.RGBA{R: 255, G: 255, B: 8, A: 255},
		color.RGBA{R: 255, G: 4, B: 2, A: 255},
	}

	return colors[idx]
}

func InitPlayer(pseudo string) *Player {
	return &Player{
		ID:     uuid.New(),
		Pseudo: pseudo,
		Game:   nil,
		Color:  getRandomColor(),
	}
}

func (p *Player) AssignToGame(game *Game) {
	p.Game = game
}
