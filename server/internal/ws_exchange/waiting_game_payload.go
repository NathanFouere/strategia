package ws_exchange

type WaitingGamePayload struct {
	SecondsBeforeLaunch    int    `json:"seconds_before_launch"`
	GameId                 string `json:"game_id"`
	NumberOfWaitingPlayers int    `json:"number_of_waiting_players"`
	IsPlayerWaitingForGame bool   `json:"is_player_waiting_for_game"`
	IsGameLaunching        bool   `json:"is_game_launching"`
}

func (waitingGamePayload *WaitingGamePayload) ToWsExchange() *WsExchangeTemplate[*WaitingGamePayload] {
	return &WsExchangeTemplate[*WaitingGamePayload]{
		Type:    "waiting_game_exchange",
		Payload: waitingGamePayload,
	}
}
