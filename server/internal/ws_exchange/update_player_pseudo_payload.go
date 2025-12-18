package ws_exchange

type UpdatePlayerPseudoPayload struct {
	PlayerId  string `json:"player_id"`
	NewPseudo string `json:"new_pseudo"`
}
