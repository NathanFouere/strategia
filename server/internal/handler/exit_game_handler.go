package handler

import (
	"server/internal/repository"
	"server/internal/ws_exchange"
	"server/pkg/logger"

	"github.com/google/uuid"
)

type ExitGameHandler struct {
	pr     *repository.PlayerRepository
	gr     *repository.GameRepository
	logger *logger.LoggerService
}

func NewExitGameHandler(
	pr *repository.PlayerRepository,
	gr *repository.GameRepository,
	logger *logger.LoggerService,
) *ExitGameHandler {
	return &ExitGameHandler{
		pr:     pr,
		gr:     gr,
		logger: logger,
	}
}

func (h *ExitGameHandler) Handle(exitGamePayload *ws_exchange.ExitGamePayload) {
	parsedPlayerId, err := uuid.Parse(exitGamePayload.PlayerId)
	if err != nil {
		h.logger.Error("Couldn't parse uuid from game subscription payload", "uuid", exitGamePayload.PlayerId, "err", err)
		return
	}

	parsedGameId, err := uuid.Parse(exitGamePayload.GameId)
	if err != nil {
		h.logger.Error("Couldn't parse uuid from game subscription payload", "uuid", exitGamePayload.GameId, "err", err)
		return
	}

	game, err := h.gr.GetGameOfId(parsedGameId)

	if err != nil {
		h.logger.Error("Couldn't pase gameId", "err", err)
		return
	}

	err = game.RemovePlayer(parsedPlayerId)
	if err != nil {
		h.logger.Error("Couldn't find player id ingame")
		return
	}
}
