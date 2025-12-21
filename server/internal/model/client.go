package model

import "github.com/gorilla/websocket"

// inspired by https://github.com/gorilla/websocket/blob/v1.4.0/examples/chat/client.go
type Client struct {
	Conn *websocket.Conn
	Send chan []byte
}

func InitClient(conn *websocket.Conn) *Client {
	return &Client{
		Conn: conn,
		Send: make(chan []byte, 256),
	}
}
