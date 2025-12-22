package service

import (
	"os"
	"server/internal/broadcaster"
	"server/internal/model"
	"server/pkg/logger"
	"strconv"
	"time"
)

type StartGameService struct {
	logger                *logger.LoggerService
	gameUpdateBroadcaster *broadcaster.GameUpdateBroadcaster
}

func NewStartGameService(logger *logger.LoggerService, gameUpdateBroadcaster *broadcaster.GameUpdateBroadcaster) *StartGameService {
	return &StartGameService{
		logger:                logger,
		gameUpdateBroadcaster: gameUpdateBroadcaster,
	}
}

func (s *StartGameService) Start(g *model.Game) {
	s.logger.Info("New game launched with ID", "id", g.ID)

	tickerUpdateGameMs, err := strconv.Atoi(os.Getenv("TICKER_UPDATE_GAME_MS"))
	if err != nil {
		panic("Couldn't load TICKER_MAIN_MENU_SECONDS in .env")
	}

	ticker := time.NewTicker(time.Duration(tickerUpdateGameMs) * time.Millisecond)
	go func() {
		for range ticker.C {
			err := s.gameUpdateBroadcaster.BroadcastGameState(g)
			if err != nil {
				s.logger.Error("Error while updating game", "game id", g.ID, "error", err)
				// TODO => see how to handle error
				return
			}
		}
	}()

}
