package handler

import (
	"server/internal/repository"
	"server/internal/ws_exchange"
	"server/pkg/logger"

	"github.com/google/uuid"
)

type UpdatePlayerPseudoHandler struct {
	pr     *repository.PlayerRepository
	logger *logger.LoggerService
}

func NewUpdatePlayerPseudoHandler(
	pr *repository.PlayerRepository,
	logger *logger.LoggerService,
) *UpdatePlayerPseudoHandler {
	return &UpdatePlayerPseudoHandler{
		pr:     pr,
		logger: logger,
	}
}

func (h *UpdatePlayerPseudoHandler) Handle(updatePlayerPseudoPayload *ws_exchange.UpdatePlayerPseudoPayload) {
	parsedPlayerId, err := uuid.Parse(updatePlayerPseudoPayload.PlayerId)
	if err != nil {
		h.logger.Error("Couldn't parse uuid from game subscription payload", "uuid", updatePlayerPseudoPayload.PlayerId, "err", err)
	}

	player, err := h.pr.GetPlayerFromId(parsedPlayerId)
	if err != nil {
		return
	}

	player.UpdatePseudo(updatePlayerPseudoPayload.NewPseudo)
}
