package ws_exchange

type RedirectToGamePayload struct {
}

func (redirectToGamePayload *RedirectToGamePayload) ToWsExchange() *WsExchangeTemplate[*RedirectToGamePayload] {
	return &WsExchangeTemplate[*RedirectToGamePayload]{
		Type:    "redirect_to_game",
		Payload: redirectToGamePayload,
	}
}
