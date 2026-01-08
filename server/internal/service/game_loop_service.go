package service

import (
	"os"
	"server/internal/model"
	"server/internal/repository"
	"server/pkg/logger"
	"strconv"
	"time"
)

type GameLoopService struct {
	logger                 *logger.LoggerService
	updateGameStateService *UpdateGameStateService
	gameRepository         *repository.GameRepository
}

func NewGameLoopService(
	logger *logger.LoggerService,
	updateGameStateService *UpdateGameStateService,
	gameRepository *repository.GameRepository,
) *GameLoopService {
	return &GameLoopService{
		logger:                 logger,
		updateGameStateService: updateGameStateService,
		gameRepository:         gameRepository,
	}
}

func (s *GameLoopService) Run(g *model.Game) {
	tickerUpdateGameMs, err := strconv.Atoi(os.Getenv("TICKER_UPDATE_GAME_MS"))
	if err != nil {
		panic("Couldn't load TICKER_MAIN_MENU_SECONDS in .env")
	}

	ticker := time.NewTicker(time.Duration(tickerUpdateGameMs) * time.Millisecond)
	go func() {
		for range ticker.C {
			err := s.updateGameStateService.UpdateGameState(g)
			if err != nil {
				s.logger.Error("Error while updating game", "game id", g.ID, "error", err)
				return
			}
			if g.Finished {
				ticker.Stop()
				err = s.gameRepository.RemoveGame(g.ID)
				if err != nil {
					s.logger.Error("Error while removing finished game from repository", "game id", g.ID, "error", err)
					return
				}
				s.logger.Info("Stopping game loop for finished game", "game id", g.ID)
				return
			}
		}
	}()
}
