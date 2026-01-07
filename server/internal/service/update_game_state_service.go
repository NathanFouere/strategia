package service

import (
	"server/internal/broadcaster"
	"server/internal/model"
	"server/pkg/logger"
)

type UpdateGameStateService struct {
	logger                *logger.LoggerService
	gameUpdateBroadcaster *broadcaster.GameUpdateBroadcaster
}

func NewUpdateGameService(
	logger *logger.LoggerService,
	gameUpdateBroadcaster *broadcaster.GameUpdateBroadcaster,
) *UpdateGameStateService {
	return &UpdateGameStateService{
		logger:                logger,
		gameUpdateBroadcaster: gameUpdateBroadcaster,
	}
}

func (s *UpdateGameStateService) UpdateGameState(g *model.Game) error {
	g.CheckGameFinished()
	if g.Finished {
		s.logger.Info("Game finished", "game id", g.ID)
		return nil
	}

	err := s.gameUpdateBroadcaster.BroadcastGameState(g)
	if err != nil {
		s.logger.Error("Error while broadcasting game state", "game id", g.ID, "error", err)
		return err
	}
	g.UpdatePlayers()
	g.ResetTilesToRender()

	return nil
}
