package ws_exchange

type ExitGamePayload struct {
	PlayerId string `json:"player_id"`
	GameId   string `json:"game_id"`
}
