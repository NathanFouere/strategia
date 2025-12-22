package sender

import (
	"encoding/json"
	"os"
	"server/internal/repository"
	"server/internal/ws_exchange"
	"server/pkg/logger"
	"strconv"
)

type PendingGameUpdateSender struct {
	pr     *repository.PlayerRepository
	gr     *repository.GameRepository
	logger *logger.LoggerService
}

func NewPendingGameUpdateSender(
	logger *logger.LoggerService,
	pr *repository.PlayerRepository,
	gr *repository.GameRepository,
) *PendingGameUpdateSender {
	return &PendingGameUpdateSender{
		pr:     pr,
		gr:     gr,
		logger: logger,
	}
}

func (s *PendingGameUpdateSender) SendPendingGameUpdate() error {
	for client := range s.pr.ClientsInLobby {
		isClientWaitingForGame := false
		_, ok := s.pr.WaitingGameClients[client]
		if ok {
			isClientWaitingForGame = true
		}

		timeBetweenGame, err := strconv.Atoi(os.Getenv("TIME_BETWEEN_GAME"))
		if err != nil {
			return err
		}

		data := &ws_exchange.WaitingGamePayload{
			SecondsBeforeLaunch:    timeBetweenGame - s.gr.CounterBetweenGames,
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
