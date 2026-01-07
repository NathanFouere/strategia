package model

import (
	"errors"
	"fmt"
	"server/internal/ws_exchange"

	"github.com/google/uuid"
)

type Game struct {
	ID               uuid.UUID
	TilesToRender    map[string]string  // pos x - y => id player
	TilesDict        map[string]string  // pos x - y => id player
	NbTilesPerPlayer map[string]int     // id player => nb tiles controlled
	Players          map[string]*Player // id player => nb tiles controlled
	Finished         bool
}

func InitGame() *Game {
	return &Game{
		ID:               uuid.New(),
		TilesToRender:    map[string]string{},
		TilesDict:        map[string]string{},
		Players:          map[string]*Player{},
		NbTilesPerPlayer: map[string]int{},
		Finished:         false,
	}
}

func (g *Game) AddPlayer(player *Player) {
	g.Players[player.ID.String()] = player
}

func (g *Game) RemovePlayer(playerId uuid.UUID) error {
	if _, exists := g.Players[playerId.String()]; exists {
		delete(g.Players, playerId.String())
	}

	return errors.New("couldn't find player of id in game")
}

func (g *Game) ReceivePixelClick(pixelClick *ws_exchange.PixelClickPayload) {
	tile := fmt.Sprintf("%d-%d", pixelClick.X, pixelClick.Y)
	g.TilesToRender[tile] = pixelClick.IdPlayer
	g.TilesDict[tile] = pixelClick.IdPlayer
	g.NbTilesPerPlayer[pixelClick.IdPlayer]++
}

func (g *Game) FindPlayerOfIdInGame(playerId uuid.UUID) (*Player, error) {
	if _, exists := g.Players[playerId.String()]; exists {
		return g.Players[playerId.String()], nil
	}

	return nil, errors.New("player not found")
}

func (g *Game) ResetTilesToRender() {
	g.TilesToRender = map[string]string{}
}

func (g *Game) ResetState() {
	g.TilesDict = map[string]string{}
}

func (g *Game) UpdatePlayers() error {
	for _, v := range g.Players {
		if val, exists := g.NbTilesPerPlayer[v.ID.String()]; exists {
			v.NbTilesControlled = val
		} else {
			return errors.New("didn't find player in nbtilesperplayer map")
		}
		v.UpdatePopulation()
	}

	return nil
}

func (g *Game) CheckGameFinished() {
	// TODO => à compléter avec les conditions de victoire
	g.Finished = g.GetNbPlayers() == 0
}

func (g *Game) GetNbPlayers() int {
	return len(g.Players)
}
