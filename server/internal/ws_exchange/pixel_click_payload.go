package ws_exchange

type PixelClickPayload struct {
	X        int    `json:"x"`
	Y        int    `json:"y"`
	IdPlayer string `json:"id_player"`
	GameId   string `json:"game_id"`
}

func (pixelCLickPayload *PixelClickPayload) ToWsExchange() *WsExchangeTemplate[*PixelClickPayload] {
	return &WsExchangeTemplate[*PixelClickPayload]{
		Type:    "pixel_click_evt",
		Payload: pixelCLickPayload,
	}
}
