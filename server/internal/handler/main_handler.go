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
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type MainHandler struct {
	upgrader            websocket.Upgrader
	clients             map[uuid.UUID]*websocket.Conn
	waitingGameClients  map[uuid.UUID]*websocket.Conn
	mutex               *sync.Mutex
	broadcast           chan []byte
	logger              *logger.LoggerService
	playerRepository    *repository.PlayerRepository
	ongoingGames        []*model.Game
	pendingGame         *model.Game
	counterBetweenGames int
}

func NewMainHandler(logger *logger.LoggerService, playerRepository *repository.PlayerRepository) *MainHandler {
	return &MainHandler{
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		},
		clients:             make(map[uuid.UUID]*websocket.Conn),
		waitingGameClients:  make(map[uuid.UUID]*websocket.Conn),
		mutex:               &sync.Mutex{},
		broadcast:           make(chan []byte),
		logger:              logger,
		playerRepository:    playerRepository,
		ongoingGames:        []*model.Game{},
		pendingGame:         model.InitGame(),
		counterBetweenGames: 0,
	}
}

func (mh *MainHandler) wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := mh.upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error upgrading:", err)
		return
	}

	mh.mutex.Lock()
	mh.logger.Info("New client as connected with addr", "addr", conn.RemoteAddr())
	player := model.InitPlayer("test")
	mh.clients[player.ID] = conn
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
			// delete(mh.clients, conn)
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
		switch exchangeRaw.Type {
		case "game-subscription":
			fmt.Println("received game subscription evt")
			payload, err := ws_exchange.ExtractTypedPayload[ws_exchange.GameSubscriptionPayload](&exchangeRaw)
			if err != nil {
				fmt.Println("Payload ERROR:", err)
				continue
			}
			mh.handleGameSubscription(payload)
		}

		mh.mutex.Unlock()
	}
}

func (mh *MainHandler) handleGameSubscription(gameSubscriptionPayload *ws_exchange.GameSubscriptionPayload) {
	parsedUUid, err := uuid.Parse(gameSubscriptionPayload.PlayerId)
	if err != nil {
		mh.logger.Error("Couldn't parse uuid from game subscription payload", "uuid", gameSubscriptionPayload.PlayerId, "err", err)
	}
	_, ok := mh.waitingGameClients[parsedUUid]
	if ok {
		mh.logger.Info("Unsubscribe player from game", "playerid", parsedUUid, "gameid", mh.pendingGame.ID)
		delete(mh.waitingGameClients, parsedUUid)
		return
	}

	mh.logger.Info("Subscribe player to game", "playerid", parsedUUid, "gameid", mh.pendingGame.ID)
	subscribingPlayer, err := mh.playerRepository.GetPlayerFromId(parsedUUid)
	mh.pendingGame.AddPlayer(subscribingPlayer)
	connCorrespondingToUuid := mh.clients[parsedUUid]
	mh.waitingGameClients[parsedUUid] = connCorrespondingToUuid
}

func (mh *MainHandler) update() error {
	mh.logger.Info("UPDATE: ", "pending game id", mh.pendingGame.ID, "counter", mh.counterBetweenGames)
	mh.counterBetweenGames++
	if mh.counterBetweenGames == 10 { // TODO => enelever hardcode
		mh.pendingGame.Start()
		mh.ongoingGames = append(mh.ongoingGames, mh.pendingGame)
		mh.logger.Info("New game launched with ID", "id", mh.pendingGame.ID)
		mh.pendingGame = model.InitGame()
		mh.waitingGameClients = make(map[uuid.UUID]*websocket.Conn)
		mh.counterBetweenGames = 0
	}

	err := mh.sendPendingGameUpdate()
	if err != nil {
		return err
	}

	return nil
}

func (mh *MainHandler) sendPendingGameUpdate() error {
	for client := range mh.clients {
		isClientWaitingForGame := false
		_, ok := mh.waitingGameClients[client]
		if ok {
			isClientWaitingForGame = true
		}
		data := &ws_exchange.WaitingGamePayload{
			SecondsBeforeLaunch:    10 - mh.counterBetweenGames, // TODO => enlever hardcode
			GameId:                 mh.pendingGame.ID.String(),
			NumberOfWaitingPlayers: len(mh.waitingGameClients),
			IsPlayerWaitingForGame: isClientWaitingForGame,
		}

		bytes, err := json.Marshal(data.ToWsExchange())
		if err != nil {
			return err
		}

		err = mh.clients[client].WriteMessage(websocket.TextMessage, bytes)
		if err != nil {
			return err
		}
	}

	return nil

}

func (mh *MainHandler) Launch() {
	ticker := time.NewTicker(1 * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				err := mh.update()
				if err != nil {
					return
				}
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
