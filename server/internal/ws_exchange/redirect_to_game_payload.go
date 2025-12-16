package ws_exchange

type RedirectToGamePayload struct {
	GameId string `json:"game_id"`
}

func (redirectToGamePayload *RedirectToGamePayload) ToWsExchange() *WsExchangeTemplate[*RedirectToGamePayload] {
	return &WsExchangeTemplate[*RedirectToGamePayload]{
		Type:    "redirect_to_game",
		Payload: redirectToGamePayload,
	}
}
