package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"server/internal/ws_exchange"
	"strings"

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

func (Game *Game) AddPlayer(player *Player) {
	Game.Players = append(Game.Players, player)
}
func (Game *Game) Start() {
	fmt.Println("Game started !")
}

func (game *Game) ReceivePixelClick(pixelClick *ws_exchange.PixelClickPayload) {
	tile := fmt.Sprintf("%d-%d", pixelClick.X, pixelClick.Y)
	game.TilesToRender[tile] = pixelClick.IdPlayer
	game.TilesDict[tile] = pixelClick.IdPlayer
}

func (game *Game) findPlayerOfIdInGame(playerId uuid.UUID) (*Player, error) {
	for _, val := range game.Players {
		if val.ID == playerId {
			return val, nil
		}
	}

	return nil, errors.New("player not found")
}

func (game *Game) generateServerUpdate() *ws_exchange.ServerUpdatePayload {
	updates := []ws_exchange.ServerUpdateData{}
	fmt.Println(game.TilesDict)
	for key, val := range game.TilesDict {
		parts := strings.Split(key, "-")
		x := parts[0]
		y := parts[1]
		playerId, err := uuid.Parse(val)
		if err != nil {
			fmt.Println("error", err)
			continue
			// todo => should throw error
		}

		player, err := game.findPlayerOfIdInGame(playerId)
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

func (game *Game) Update(conn *websocket.Conn) error {
	data := game.generateServerUpdate()
	fmt.Println(data)

	bytes, err := json.Marshal(data.ToWsExchange())
	if err != nil {
		return err
	}

	err = conn.WriteMessage(websocket.TextMessage, bytes)
	if err != nil {
		return err
	}

	return nil
}

func (game *Game) ResetState() {
	game.TilesDict = map[string]string{}
}
