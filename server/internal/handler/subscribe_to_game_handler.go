package handler

import (
	"server/internal/repository"
	"server/internal/ws_exchange"
	"server/pkg/logger"

	"github.com/google/uuid"
)

type SubscribeToGameHandler struct {
	pr     *repository.PlayerRepository
	gr     *repository.GameRepository
	logger *logger.LoggerService
}

func NewSubscribeToGameHandler(
	pr *repository.PlayerRepository,
	gr *repository.GameRepository,
	logger *logger.LoggerService,
) *SubscribeToGameHandler {
	return &SubscribeToGameHandler{
		pr:     pr,
		gr:     gr,
		logger: logger,
	}
}

func (h *SubscribeToGameHandler) Handle(gameSubscriptionPayload *ws_exchange.GameSubscriptionPayload) {
	parsedUUid, err := uuid.Parse(gameSubscriptionPayload.PlayerId)
	if err != nil {
		h.logger.Error("Couldn't parse uuid from game subscription payload", "uuid", gameSubscriptionPayload.PlayerId, "err", err)
	}
	_, ok := h.pr.WaitingGameClients[parsedUUid]
	if ok {
		h.logger.Error("Player already subscribed to game", "playerid", parsedUUid, "gameid", h.gr.PendingGame.ID)
		return
	}

	h.logger.Info("Subscribe player to game", "playerid", parsedUUid, "gameid", h.gr.PendingGame.ID)
	subscribingPlayer, err := h.pr.GetPlayerFromId(parsedUUid)
	h.gr.PendingGame.AddPlayer(subscribingPlayer)
	connCorrespondingToUuid := h.pr.ClientsInLobby[parsedUUid]
	h.pr.WaitingGameClients[parsedUUid] = connCorrespondingToUuid
}
