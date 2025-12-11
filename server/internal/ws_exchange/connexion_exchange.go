package ws_exchange

type ConnexionExchange struct {
	PlayerId string `json:"player-id"`
}

func (connexionExchange *ConnexionExchange) ToWsExchange() *WsExchangeTemplate {
	return &WsExchangeTemplate{
		Type:    "connexion-exchange",
		Payload: connexionExchange,
	}
}
