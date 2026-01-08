package service

import (
	"server/internal/broadcaster"
	"server/internal/model"
	"server/pkg/logger"
)

type UpdateGameStateService struct {
	logger                       *logger.LoggerService
	gameUpdateBroadcaster        *broadcaster.GameUpdateBroadcaster
	gameStartupUpdateBroadcaster *broadcaster.GameStartupUpdateBroadcaster
}

func NewUpdateGameService(
	logger *logger.LoggerService,
	gameUpdateBroadcaster *broadcaster.GameUpdateBroadcaster,
	gameStartupUpdateBroadcaster *broadcaster.GameStartupUpdateBroadcaster,
) *UpdateGameStateService {
	return &UpdateGameStateService{
		logger:                       logger,
		gameUpdateBroadcaster:        gameUpdateBroadcaster,
		gameStartupUpdateBroadcaster: gameStartupUpdateBroadcaster,
	}
}

func (s *UpdateGameStateService) UpdateGameState(g *model.Game) error {
	err := s.gameUpdateBroadcaster.BroadcastGameState(g)
	if err != nil {
		s.logger.Error("Error while broadcasting game state", "game id", g.ID, "error", err)
		return err
	}

	if !g.Started {
		g.TimerBeforeStart--
		g.CheckGameStarted()
		err = s.gameStartupUpdateBroadcaster.BroadcastGameStartupUpdate(g)
		if err != nil {
			s.logger.Error("Error while broadcasting game startup update", "game id", g.ID, "error", err)
			return err
		}
		return nil
	}

	g.CheckGameFinished()
	if g.Finished {
		s.logger.Info("Game finished", "game id", g.ID)
		return nil
	}
	err = g.UpdatePlayers()
	if err != nil {
		return err
	}

	g.ResetTilesToRender()

	return nil
}
