package handler

import (
	"server/internal/repository"
	"server/internal/ws_exchange"
	"server/pkg/logger"

	"github.com/google/uuid"
)

type UnsubscribeFromGameHandler struct {
	logger *logger.LoggerService
	pr     *repository.PlayerRepository
	gr     *repository.GameRepository
}

func NewUnsubscribeFromGameHandler(
	logger *logger.LoggerService,
	pr *repository.PlayerRepository,
	gr *repository.GameRepository,
) *UnsubscribeFromGameHandler {
	return &UnsubscribeFromGameHandler{
		logger: logger,
		pr:     pr,
		gr:     gr,
	}
}

func (h *UnsubscribeFromGameHandler) Handle(gameUnsubscribePayload *ws_exchange.GameUnsubscribePayload) {
	parsedUUid, err := uuid.Parse(gameUnsubscribePayload.PlayerId)
	if err != nil {
		h.logger.Error("Couldn't parse uuid from game subscription payload", "uuid", gameUnsubscribePayload.PlayerId, "err", err)
		return
	}

	// TODO => refactoriser cette fonction
	_, ok := h.pr.WaitingGameClients[parsedUUid]
	if ok {
		h.logger.Info("Unsubscribe player from game", "playerid", parsedUUid, "gameid", h.gr.PendingGame.ID)
		err := h.gr.PendingGame.RemovePlayer(parsedUUid)
		if err != nil {
			return
		}
		delete(h.pr.WaitingGameClients, parsedUUid)
		return
	}

	h.logger.Error("Tried to unsubscribe unsubscribed player", "playerid", parsedUUid, "gameid", h.gr.PendingGame.ID)
}
