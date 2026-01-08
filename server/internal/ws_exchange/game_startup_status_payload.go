package ws_exchange

type GameStartupStatusPayload struct {
	ProgressionPercentage int  `json:"progression_percentage"`
	GameStarted           bool `json:"game_started"`
}

func (gameStartupStatusPayload *GameStartupStatusPayload) ToWsExchange() *WsExchangeTemplate[*GameStartupStatusPayload] {
	return &WsExchangeTemplate[*GameStartupStatusPayload]{
		Type:    "game_startup_status",
		Payload: gameStartupStatusPayload,
	}
}
