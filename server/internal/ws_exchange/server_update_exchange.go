package ws_exchange

type ServerUpdateData struct {
	X     string `json:"x"`
	Y     string `json:"y"`
	Color string `json:"color"`
}

type ServerUpdateModel struct {
	ServerUpdateDatas []ServerUpdateData `json:"update-datas"`
}

func (serverUpdateModel *ServerUpdateModel) ToWsExchange() *WsExchangeTemplate {
	return &WsExchangeTemplate{
		Type:    "server-update-datas",
		Payload: serverUpdateModel.ServerUpdateDatas,
	}
}
