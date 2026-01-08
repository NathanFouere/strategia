package model

import (
	"errors"
	"fmt"
	"os"
	"server/internal/ws_exchange"
	"strconv"

	"github.com/google/uuid"
)

type Game struct {
	ID               uuid.UUID
	TilesToRender    map[string]string  // pos x - y => id player
	TilesDict        map[string]string  // pos x - y => id player
	NbTilesPerPlayer map[string]int     // id player => nb tiles controlled
	Players          map[string]*Player // id player => nb tiles controlled
	Started          bool
	Finished         bool
	TimerBeforeStart int
}

func InitGame() *Game {
	tickerUpdateGameMs, err := strconv.Atoi(os.Getenv("TICKER_UPDATE_GAME_MS"))
	if err != nil {
		panic("Couldn't read TICKER_UPDATE_GAME_MS env var")
	}
	startDelaySec, err := strconv.Atoi(os.Getenv("GAME_START_DELAY_SEC"))
	if err != nil {
		panic("Couldn't read GAME_START_DELAY_SEC env var")
	}

	timerBeforeStartTicks := startDelaySec * 1000 / tickerUpdateGameMs
	return &Game{
		ID:               uuid.New(),
		TilesToRender:    map[string]string{},
		TilesDict:        map[string]string{},
		Players:          map[string]*Player{},
		NbTilesPerPlayer: map[string]int{},
		Started:          false,
		Finished:         false,
		TimerBeforeStart: timerBeforeStartTicks,
	}
}

func (g *Game) AddPlayer(player *Player) {
	g.Players[player.ID.String()] = player
}

func (g *Game) RemovePlayer(playerId uuid.UUID) error {
	if _, exists := g.Players[playerId.String()]; !exists {
		return errors.New("couldn't find player of id in game")
	}

	delete(g.Players, playerId.String())
	return nil
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

func (g *Game) CheckGameStarted() {
	if g.TimerBeforeStart <= 0 {
		g.Started = true
	}
}

func (g *Game) GetNbPlayers() int {
	return len(g.Players)
}
