package service

import (
	"server/internal/model"
	"server/internal/repository"
	"server/internal/sender"
	"server/pkg/logger"

	"github.com/google/uuid"
)

type UpdateService struct {
	pr                      *repository.PlayerRepository
	gr                      *repository.GameRepository
	logger                  *logger.LoggerService
	pendingGameUpdateSender *sender.PendingGameUpdateSender
	redirectToGameSender    *sender.RedirectToGameSender
	startGameService        *StartGameService
}

func NewUpdateService(
	pr *repository.PlayerRepository,
	gr *repository.GameRepository,
	logger *logger.LoggerService,
	pendingGameUpdateSender *sender.PendingGameUpdateSender,
	redirectToGameSender *sender.RedirectToGameSender,
	startGameService *StartGameService,
) *UpdateService {
	return &UpdateService{
		pr:                      pr,
		gr:                      gr,
		logger:                  logger,
		pendingGameUpdateSender: pendingGameUpdateSender,
		redirectToGameSender:    redirectToGameSender,
		startGameService:        startGameService,
	}
}

func (s *UpdateService) Update() error {
	s.logger.Info("UPDATE: ", "pending game id", s.gr.PendingGame.ID, "counter", s.gr.CounterBetweenGames)
	s.gr.CounterBetweenGames++
	if s.gr.CounterBetweenGames == 10 && len(s.gr.PendingGame.Players) > 0 { // TODO => enelever hardcode
		s.redirectToGameSender.SendRedirectToGame()
		s.startGameService.Start(s.gr.PendingGame)
		s.gr.OngoingGames = append(s.gr.OngoingGames, s.gr.PendingGame)
		s.gr.PendingGame = model.InitGame()
		s.gr.AddGame(s.gr.PendingGame)
		s.pr.WaitingGameClients = make(map[uuid.UUID]*model.Player)
		s.gr.CounterBetweenGames = 0
		return nil
	} else if s.gr.CounterBetweenGames == 10 && len(s.gr.PendingGame.Players) == 0 {
		s.gr.CounterBetweenGames = 0
		return nil
	}

	err := s.pendingGameUpdateSender.SendPendingGameUpdate()
	if err != nil {
		return err
	}

	return nil
}
