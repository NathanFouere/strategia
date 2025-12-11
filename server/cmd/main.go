package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"server/internal/model"
	"server/internal/ws_exchange"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}
var clients = make(map[*websocket.Conn]bool)
var mutex = &sync.Mutex{}
var broadcast = make(chan []byte) // Broadcast channel

// from: https://freedium-mirror.cfd/https://medium.com/wisemonks/implementing-websockets-in-golang-d3e8e219733b
func wsHandler(w http.ResponseWriter, r *http.Request, g *model.Game) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error upgrading:", err)
		return
	}

	mutex.Lock()
	clients[conn] = true
	fmt.Println("New client_old:", conn.RemoteAddr())
	player := model.InitPlayer(g)
	g.AddPlayer(player)
	mutex.Unlock()

	// envoie du connexion exchange
	connexionExchange := &ws_exchange.ConnexionExchange{
		PlayerId: player.ID.String(),
	}

	data, err := json.Marshal(connexionExchange.ToWsExchange())
	err = conn.WriteMessage(websocket.TextMessage, data)
	if err != nil {
		fmt.Println("error", err)
		return
	}

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			mutex.Lock()
			delete(clients, conn)
			mutex.Unlock()
			break
		}
		broadcast <- message
	}
}

func handleMessages(g *model.Game) {
	for {
		message := <-broadcast
		mutex.Lock()
		for client := range clients {
			fmt.Println(message)
			var click ws_exchange.PixelClickExchange
			err := json.Unmarshal(message, &click)
			if err != nil {
				fmt.Println("JSON error:", err)
				return
			}
			g.ReceivePixelClick(&click)
			if err != nil {
				err := client.Close()
				if err != nil {
					return
				}
				delete(clients, client)
			}
		}
		mutex.Unlock()
	}
}

func makeWsHandler(g *model.Game) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		wsHandler(w, r, g)
	}
}

func update(g *model.Game) {
	fmt.Println("Update")
	for conn := range clients {
		err := g.Update(conn)
		if err != nil {
			return
		}
	}
	g.ResetState()
}

func main() {
	game := model.InitGame()
	http.HandleFunc("/ws", makeWsHandler(game))
	go handleMessages(game)
	ticker := time.NewTicker(1 * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				update(game)
			}
		}
	}()
	fmt.Println("Serveur WebSocket démarré sur :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
