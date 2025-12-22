package sender

import (
	"encoding/json"
	"server/internal/model"
	"server/internal/ws_exchange"
	"server/pkg/logger"

	"github.com/gorilla/websocket"
)

type ConnectionExchangeSender struct {
	logger *logger.LoggerService
}

func NewConnectionExchangeSender(logger *logger.LoggerService) *ConnectionExchangeSender {
	return &ConnectionExchangeSender{
		logger: logger,
	}
}

func (s *ConnectionExchangeSender) Send(player *model.Player) {

	connexionExchange := &ws_exchange.ConnectionPayload{
		PlayerId:     player.ID.String(),
		PlayerPseudo: player.Pseudo,
	}

	data, err := json.Marshal(connexionExchange.ToWsExchange())
	err = player.Client.Conn.WriteMessage(websocket.TextMessage, data)
	if err != nil {
		s.logger.Error("Error while send Connection Exchange", "err", err)
	}
}
