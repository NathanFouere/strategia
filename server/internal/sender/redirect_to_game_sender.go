package sender

import (
	"encoding/json"
	"server/internal/repository"
	"server/internal/ws_exchange"
	"server/pkg/logger"
)

type RedirectToGameSender struct {
	pr     *repository.PlayerRepository
	gr     *repository.GameRepository
	logger *logger.LoggerService
}

func NewRedirectToGameSender(
	pr *repository.PlayerRepository,
	gr *repository.GameRepository,
	logger *logger.LoggerService,
) *RedirectToGameSender {
	return &RedirectToGameSender{
		pr:     pr,
		gr:     gr,
		logger: logger,
	}
}

func (s *RedirectToGameSender) SendRedirectToGame() error {
	data := &ws_exchange.RedirectToGamePayload{
		GameId: s.gr.PendingGame.ID.String(),
	}

	bytes, err := json.Marshal(data.ToWsExchange())
	if err != nil {
		return err
	}

	for client := range s.pr.WaitingGameClients {
		s.pr.WaitingGameClients[client].Client.Send <- bytes
		delete(s.pr.WaitingGameClients, client)
		delete(s.pr.ClientsInLobby, client)
	}

	return nil
}
