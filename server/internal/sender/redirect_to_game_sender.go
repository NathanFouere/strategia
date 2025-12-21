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
	for client := range s.pr.ClientsInLobby {
		isClientWaitingForGame := false
		_, ok := s.pr.WaitingGameClients[client]
		if ok {
			isClientWaitingForGame = true
		}
		data := &ws_exchange.WaitingGamePayload{
			SecondsBeforeLaunch:    10 - s.gr.CounterBetweenGames, // TODO => enlever hardcode
			GameId:                 s.gr.PendingGame.ID.String(),
			NumberOfWaitingPlayers: len(s.pr.WaitingGameClients),
			IsPlayerWaitingForGame: isClientWaitingForGame,
			IsGameLaunching:        s.gr.CounterBetweenGames == 0,
		}

		bytes, err := json.Marshal(data.ToWsExchange())
		if err != nil {
			return err
		}

		s.pr.ClientsInLobby[client].Client.Send <- bytes
	}

	return nil
}
