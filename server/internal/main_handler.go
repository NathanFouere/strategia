package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
	"server/internal/model"
	"server/internal/repository"
	"server/internal/service"
	"server/internal/ws_exchange"
	"server/pkg/logger"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type MainHandler struct {
	upgrader             websocket.Upgrader
	mutex                *sync.Mutex
	broadcast            chan []byte
	logger               *logger.LoggerService
	playerRepository     *repository.PlayerRepository
	gameRepository       *repository.GameRepository
	messageRouterService *service.MessageRouterService
	updateService        *service.UpdateService
}

func NewMainHandler(
	logger *logger.LoggerService,
	playerRepository *repository.PlayerRepository,
	gameRepository *repository.GameRepository,
	messageRouterService *service.MessageRouterService,
	updateService *service.UpdateService,
) *MainHandler {
	mainHandler := &MainHandler{
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		},
		mutex:                &sync.Mutex{},
		broadcast:            make(chan []byte),
		logger:               logger,
		playerRepository:     playerRepository,
		gameRepository:       gameRepository,
		messageRouterService: messageRouterService,
		updateService:        updateService,
	}

	pending := model.InitGame()
	gameRepository.PendingGame = pending
	gameRepository.AddGame(pending)

	return mainHandler
}

func (mh *MainHandler) wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := mh.upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error upgrading:", err)
		return
	}

	mh.mutex.Lock()
	mh.logger.Info("New client has connected with addr", "addr", conn.RemoteAddr())
	player := model.InitPlayer(conn)
	mh.playerRepository.AddPlayer(player)
	mh.playerRepository.ClientsInLobby[player.ID] = player
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
			// TODO => résoudre ça
			// delete(mh.playerRepository.ClientsInLobby, conn)
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

		var exchangeRaw ws_exchange.WsExchangeTemplateRaw
		err := json.Unmarshal(message, &exchangeRaw)
		if err != nil {
			fmt.Println("JSON ERROR:", err)
		}
		mh.messageRouterService.HandleMessage(exchangeRaw)

		mh.mutex.Unlock()
	}
}

func (mh *MainHandler) Launch() {
	ticker := time.NewTicker(1 * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				mh.updateService.Update()
			}
		}
	}()
	http.HandleFunc("/ws", mh.wsHandler)
	go mh.handleMessages()
	fmt.Println("Serveur WebSocket démarré sur :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
