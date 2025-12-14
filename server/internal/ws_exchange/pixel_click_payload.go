package ws_exchange

type PixelClickPayload struct {
	X        int    `json:"x"`
	Y        int    `json:"y"`
	IdPlayer string `json:"id_player"`
}
