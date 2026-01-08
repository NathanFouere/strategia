package broadcaster

import (
	"encoding/json"
	"os"
	"server/internal/model"
	"server/internal/ws_exchange"
	"server/pkg/logger"
	"strconv"
)

type GameStartupUpdateBroadcaster struct {
	logger *logger.LoggerService
}

func NewGameStartupUpdateBroadcaster(logger *logger.LoggerService) *GameStartupUpdateBroadcaster {
	return &GameStartupUpdateBroadcaster{
		logger: logger,
	}
}

func (s *GameStartupUpdateBroadcaster) BroadcastGameStartupUpdate(g *model.Game) error {
	tickerUpdateGameMs, err := strconv.Atoi(os.Getenv("TICKER_UPDATE_GAME_MS"))
	if err != nil {
		panic("Couldn't read TICKER_UPDATE_GAME_MS env var")
	}
	startDelaySec, err := strconv.Atoi(os.Getenv("GAME_START_DELAY_SEC"))
	if err != nil {
		panic("Couldn't read GAME_START_DELAY_SEC env var")
	}
	totalTicks := startDelaySec * 1000 / tickerUpdateGameMs
	percentage := 0
	if totalTicks > 0 {
		ticksElapsed := totalTicks - g.TimerBeforeStart
		percentage = (ticksElapsed * 100) / totalTicks
	}

	data := &ws_exchange.GameStartupStatusPayload{
		ProgressionPercentage: percentage,
		GameStarted:           g.Started,
	}

	bytes, err := json.Marshal(data.ToWsExchange())
	if err != nil {
		return err
	}

	for _, player := range g.Players {
		player.Client.Send <- bytes
	}

	return nil
}
