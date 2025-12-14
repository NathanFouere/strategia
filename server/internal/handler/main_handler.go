package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"server/internal/model"
	"server/internal/repository"
	"server/internal/ws_exchange"
	"server/pkg/logger"
	"sync"

	"github.com/gorilla/websocket"
)

type MainHandler struct {
	upgrader         websocket.Upgrader
	clients          map[*websocket.Conn]bool
	mutex            *sync.Mutex
	broadcast        chan []byte
	logger           *logger.LoggerService
	playerRepository *repository.PlayerRepository
}

func NewMainHandler(logger *logger.LoggerService, playerRepository *repository.PlayerRepository) *MainHandler {
	return &MainHandler{
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		},
		clients:          make(map[*websocket.Conn]bool),
		mutex:            &sync.Mutex{},
		broadcast:        make(chan []byte),
		logger:           logger,
		playerRepository: playerRepository,
	}
}

func (mh *MainHandler) wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := mh.upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error upgrading:", err)
		return
	}

	mh.mutex.Lock()
	mh.clients[conn] = true
	mh.logger.Info("New client as connected with addr", "addr", conn.RemoteAddr())
	player := model.InitPlayer("test")
	mh.mutex.Unlock()

	// todo => bouger dans un service
	connexionExchange := &ws_exchange.ConnectionPayload{
		PlayerId:     player.ID.String(),
		PlayerPseudo: player.Pseudo,
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
			mh.mutex.Lock()
			delete(mh.clients, conn)
			mh.mutex.Unlock()
			break
		}
		mh.broadcast <- message
	}
}

func (mh *MainHandler) handleMessages() {
	for {
		message := <-mh.broadcast
		mh.mutex.Lock()
		for client := range mh.clients {
			fmt.Println(message)
			if err != nil {
				err := client.Close()
				if err != nil {
					return
				}
				delete(mh.clients, client)
			}
		}
		mh.mutex.Unlock()
	}
}

func (mh *MainHandler) Launch() {
	http.HandleFunc("/ws", mh.wsHandler)
	go mh.handleMessages()
	fmt.Println("Serveur WebSocket démarré sur :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
