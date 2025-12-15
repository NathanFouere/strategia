package ws_exchange

type ConnectionPayload struct {
	PlayerId     string `json:"player_id"`
	PlayerPseudo string `json:"player_pseudo"`
}

func (connectionPayload *ConnectionPayload) ToWsExchange() *WsExchangeTemplate[*ConnectionPayload] {
	return &WsExchangeTemplate[*ConnectionPayload]{
		Type:    "connexion-exchange",
		Payload: connectionPayload,
	}
}
