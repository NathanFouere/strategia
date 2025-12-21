package model

import (
	"image/color"
	"math/rand"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Player struct {
	ID     uuid.UUID
	Pseudo string
	Game   *Game
	Color  color.Color
	Client *Client
}

func getRandomColor() color.Color {
	idx := rand.Intn(20)
	colors := [20]color.Color{
		color.RGBA{R: 255, G: 0, B: 0, A: 255},     // Rouge vif
		color.RGBA{R: 0, G: 255, B: 0, A: 255},     // Vert vif
		color.RGBA{R: 0, G: 0, B: 255, A: 255},     // Bleu vif
		color.RGBA{R: 255, G: 255, B: 0, A: 255},   // Jaune
		color.RGBA{R: 255, G: 0, B: 255, A: 255},   // Magenta
		color.RGBA{R: 0, G: 255, B: 255, A: 255},   // Cyan
		color.RGBA{R: 255, G: 128, B: 0, A: 255},   // Orange vif
		color.RGBA{R: 128, G: 0, B: 255, A: 255},   // Violet
		color.RGBA{R: 0, G: 128, B: 255, A: 255},   // Bleu azur
		color.RGBA{R: 0, G: 255, B: 128, A: 255},   // Vert turquoise
		color.RGBA{R: 255, G: 64, B: 0, A: 255},    // Orange rouge
		color.RGBA{R: 255, G: 0, B: 128, A: 255},   // Rose fuchsia
		color.RGBA{R: 128, G: 255, B: 0, A: 255},   // Vert citron
		color.RGBA{R: 0, G: 255, B: 64, A: 255},    // Vert menthe
		color.RGBA{R: 0, G: 64, B: 255, A: 255},    // Bleu électrique
		color.RGBA{R: 128, G: 128, B: 255, A: 255}, // Lavande vive
		color.RGBA{R: 255, G: 128, B: 255, A: 255}, // Rose clair saturé
		color.RGBA{R: 255, G: 255, B: 128, A: 255}, // Jaune pastel vif
		color.RGBA{R: 128, G: 255, B: 255, A: 255}, // Cyan clair vif
		color.RGBA{R: 255, G: 64, B: 64, A: 255},   // Rouge corai
	}

	return colors[idx]
}

func InitPlayer(conn *websocket.Conn) *Player {
	playerId := uuid.New()
	return &Player{
		ID:     playerId,
		Pseudo: playerId.String(),
		Game:   nil,
		Color:  getRandomColor(),
		Client: InitClient(conn),
	}
}

func (p *Player) AssignToGame(game *Game) {
	p.Game = game
}

func (p *Player) UpdatePseudo(newPseudo string) {
	p.Pseudo = newPseudo
}
