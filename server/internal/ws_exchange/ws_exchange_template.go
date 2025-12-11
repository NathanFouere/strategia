package ws_exchange

type WsExchangeTemplate struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}
