package model

import (
	"errors"
	"fmt"
	"server/internal/ws_exchange"

	"github.com/google/uuid"
)

type Game struct {
	ID            uuid.UUID
	TilesToRender map[string]string // pos x - y => id player
	TilesDict     map[string]string // pos x - y => id player
	Players       []*Player
	Finished      bool
}

func InitGame() *Game {
	return &Game{
		ID:            uuid.New(),
		TilesToRender: map[string]string{},
		TilesDict:     map[string]string{},
		Finished:      false,
	}
}

func (g *Game) AddPlayer(player *Player) {
	g.Players = append(g.Players, player)
}

func (g *Game) RemovePlayer(playerId uuid.UUID) error {
	for i := 0; i < len(g.Players); i++ {
		if g.Players[i].ID == playerId {
			g.Players[i] = g.Players[len(g.Players)-1]
			g.Players = g.Players[:len(g.Players)-1]

			return nil
		}
	}

	return errors.New("couldn't find player of id in game")
}

func (g *Game) ReceivePixelClick(pixelClick *ws_exchange.PixelClickPayload) {
	tile := fmt.Sprintf("%d-%d", pixelClick.X, pixelClick.Y)
	g.TilesToRender[tile] = pixelClick.IdPlayer
	g.TilesDict[tile] = pixelClick.IdPlayer
}

func (g *Game) FindPlayerOfIdInGame(playerId uuid.UUID) (*Player, error) {
	for _, val := range g.Players {
		if val.ID == playerId {
			return val, nil
		}
	}

	return nil, errors.New("player not found")
}

func (g *Game) ResetState() {
	g.TilesDict = map[string]string{}
}

func (g *Game) CheckGameFinished() {
	// TODO => à compléter avec les conditions de victoire
	g.Finished = g.GetNbPlayers() == 0
}

func (g *Game) GetNbPlayers() int {
	return len(g.Players)
}
