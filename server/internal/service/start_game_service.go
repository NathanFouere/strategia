package service

import (
	"server/internal/model"
	"server/internal/repository"
	"server/pkg/logger"
)

type StartGameService struct {
	logger          *logger.LoggerService
	gameLoopService *GameLoopService
	gameRepository  *repository.GameRepository
}

func NewStartGameService(logger *logger.LoggerService, gameLoopService *GameLoopService, gameRepository *repository.GameRepository) *StartGameService {
	return &StartGameService{
		logger:          logger,
		gameLoopService: gameLoopService,
		gameRepository:  gameRepository,
	}
}

func (s *StartGameService) Start(g *model.Game) {
	s.logger.Info("New game launched with ID", "id", g.ID)

	s.gameLoopService.Run(g)
}
