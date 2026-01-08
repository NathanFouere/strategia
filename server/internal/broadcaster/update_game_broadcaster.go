package broadcaster

import (
	"encoding/json"
	"errors"
	"fmt"
	"server/internal/model"
	"server/internal/ws_exchange"
	"server/pkg/logger"
	"strings"

	"github.com/google/uuid"
)

type GameUpdateBroadcaster struct {
	logger *logger.LoggerService
}

func NewGameUpdateBroadcaster(logger *logger.LoggerService) *GameUpdateBroadcaster {
	return &GameUpdateBroadcaster{
		logger: logger,
	}
}

func (s *GameUpdateBroadcaster) buildTileUpdate(unparsedPos string, unparsedPlayerId string, game *model.Game) (*ws_exchange.ServerUpdateData, error) {
	parts := strings.Split(unparsedPos, "-")
	x := parts[0]
	y := parts[1]
	playerId, err := uuid.Parse(unparsedPlayerId)
	if err != nil {
		return nil, errors.New("error while parsing player id")
	}

	player, err := game.FindPlayerOfIdInGame(playerId)
	if err != nil {
		game.RemovePlayer(playerId)
		return nil, err
	}
	r, g, b, a := player.Color.RGBA()
	return &ws_exchange.ServerUpdateData{
		X:     x,
		Y:     y,
		Color: fmt.Sprintf("rgba(%d,%d,%d,%d)", r/257, g/257, b/257, a/257),
	}, nil
}

func (s *GameUpdateBroadcaster) BroadcastGameState(g *model.Game) error {
	updates := []ws_exchange.ServerUpdateData{}
	for key, val := range g.TilesDict {
		serverUpdateData, err := s.buildTileUpdate(key, val, g)
		if err != nil {
			return err
		}
		updates = append(updates, *serverUpdateData)
	}

	data := &ws_exchange.ServerUpdatePayload{
		ServerUpdateDatas: updates,
	}

	bytes, err := json.Marshal(data.ToWsExchange())
	if err != nil {
		return err
	}

	for _, player := range g.Players {
		player.Client.Send <- bytes
	}

	return nil
}
