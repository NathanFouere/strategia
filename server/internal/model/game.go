package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"server/internal/ws_exchange"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Game struct {
	ID            uuid.UUID
	TilesToRender map[string]string // pos x - y => id player
	TilesDict     map[string]string // pos x - y => id player
	Players       []*Player
}

func InitGame() *Game {
	return &Game{
		ID:            uuid.New(),
		TilesToRender: map[string]string{},
		TilesDict:     map[string]string{},
	}
}

func (g *Game) AddPlayer(player *Player) {
	g.Players = append(g.Players, player)
}
func (g *Game) Start() {
	fmt.Println("Game started !")
	ticker := time.NewTicker(1 * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				err := g.Update()
				if err != nil {
					fmt.Println("ERROR in game update")
					return
				}
			}
		}
	}()
}

func (g *Game) ReceivePixelClick(pixelClick *ws_exchange.PixelClickPayload) {
	fmt.Println("GAME RECEIVE PIXELCLICK", pixelClick)
	tile := fmt.Sprintf("%d-%d", pixelClick.X, pixelClick.Y)
	g.TilesToRender[tile] = pixelClick.IdPlayer
	g.TilesDict[tile] = pixelClick.IdPlayer
}

func (g *Game) findPlayerOfIdInGame(playerId uuid.UUID) (*Player, error) {
	for _, val := range g.Players {
		if val.ID == playerId {
			return val, nil
		}
	}

	return nil, errors.New("player not found")
}

func (g *Game) generateServerUpdate() *ws_exchange.ServerUpdatePayload {
	updates := []ws_exchange.ServerUpdateData{}
	fmt.Println(g.TilesDict)
	for key, val := range g.TilesDict {
		parts := strings.Split(key, "-")
		x := parts[0]
		y := parts[1]
		playerId, err := uuid.Parse(val)
		if err != nil {
			fmt.Println("error", err)
			continue
			// todo => should throw error
		}

		player, err := g.findPlayerOfIdInGame(playerId)
		if err != nil {
			fmt.Println("error", err)
			continue
			// todo => should throw error
		}
		r, g, b, a := player.Color.RGBA()
		updates = append(updates, ws_exchange.ServerUpdateData{
			X:     x,
			Y:     y,
			Color: fmt.Sprintf("rgba(%d,%d,%d,%d)", r/257, g/257, b/257, a/257),
		})
	}

	return &ws_exchange.ServerUpdatePayload{
		ServerUpdateDatas: updates,
	}
}

func (g *Game) Update() error {
	fmt.Println("UPDATE GAME !")
	data := g.generateServerUpdate()
	fmt.Println(data)

	bytes, err := json.Marshal(data.ToWsExchange())
	if err != nil {
		return err
	}

	for _, player := range g.Players {
		err = player.WsCon.WriteMessage(websocket.TextMessage, bytes)
		if err != nil {
			return err
		}
	}
	g.TilesToRender = map[string]string{}
	return nil
}

func (g *Game) ResetState() {
	g.TilesDict = map[string]string{}
}
