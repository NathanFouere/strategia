package handler

import (
	"server/internal/repository"
	"server/internal/ws_exchange"
	"server/pkg/logger"

	"github.com/google/uuid"
)

type PixelClickHandler struct {
	pr     *repository.PlayerRepository
	gr     *repository.GameRepository
	logger *logger.LoggerService
}

func NewPixelClickHandler(
	pr *repository.PlayerRepository,
	gr *repository.GameRepository,
	logger *logger.LoggerService,
) *PixelClickHandler {
	return &PixelClickHandler{
		pr:     pr,
		gr:     gr,
		logger: logger,
	}
}

func (h *PixelClickHandler) Handle(pixelClickPayload *ws_exchange.PixelClickPayload) {
	gameId, err := uuid.Parse(pixelClickPayload.GameId)
	if err != nil {
		h.logger.Error("Couldn't parse uuid from pixel click payload", "uuid", pixelClickPayload.GameId, "err", err)
		return
	}
	game, err := h.gr.GetGameOfId(gameId)
	if err != nil {
		h.logger.Error("Couldn't find id from game", "uuid", pixelClickPayload.GameId, "err", err)
		return
	}

	game.ReceivePixelClick(pixelClickPayload)
}
