package ws_exchange

import (
	"encoding/json"
	"fmt"
)

type WsExchangeTemplate[T any] struct {
	Type    string `json:"type"`
	Payload T      `json:"payload"`
}

type WsExchangeTemplateRaw struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

func ExtractTypedPayload[T any](raw *WsExchangeTemplateRaw) (*T, error) {
	var payload T
	if err := json.Unmarshal(raw.Payload, &payload); err != nil {
		return nil, fmt.Errorf("failed to unmarshal payload: %w", err)
	}
	return &payload, nil
}
