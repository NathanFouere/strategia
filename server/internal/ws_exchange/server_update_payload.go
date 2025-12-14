package ws_exchange

type ServerUpdateData struct {
	X     string `json:"x"`
	Y     string `json:"y"`
	Color string `json:"color"`
}

type ServerUpdatePayload struct {
	ServerUpdateDatas []ServerUpdateData `json:"update_datas"`
}

func (serverUpdatePayload *ServerUpdatePayload) ToWsExchange() *WsExchangeTemplate {
	return &WsExchangeTemplate{
		Type:    "server-update-datas",
		Payload: serverUpdatePayload.ServerUpdateDatas,
	}
}
