package handler

import (
	"server/internal/repository"
	"server/internal/ws_exchange"
	"server/pkg/logger"

	"github.com/google/uuid"
)

type SetInWaitingLobbyHandler struct {
	logger           *logger.LoggerService
	playerRepository *repository.PlayerRepository
}

func NewSetInWaitingLobbyHandler(
	logger *logger.LoggerService,
	pr *repository.PlayerRepository,
) *SetInWaitingLobbyHandler {
	return &SetInWaitingLobbyHandler{
		logger:           logger,
		playerRepository: pr,
	}
}

func (h *SetInWaitingLobbyHandler) Handle(setInWaitingLobbyPayload *ws_exchange.SetInWaitingLobbyPayload) {
	parsedPlayerId, err := uuid.Parse(setInWaitingLobbyPayload.PlayerId)
	if err != nil {
		h.logger.Error("Couldn't parse uuid from game subscription payload", "uuid", setInWaitingLobbyPayload.PlayerId, "err", err)
		return
	}

	player, err := h.playerRepository.GetPlayerFromId(parsedPlayerId)
	if err != nil {
		return
	}

	h.playerRepository.AddPlayerToClientLobby(player)
}
