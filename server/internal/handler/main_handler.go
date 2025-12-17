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
	clientsInLobby      map[uuid.UUID]*websocket.Conn
	waitingGameClients  map[uuid.UUID]*websocket.Conn
	mutex               *sync.Mutex
	broadcast           chan []byte
	logger              *logger.LoggerService
	playerRepository    *repository.PlayerRepository
	ongoingGames        []*model.Game
	pendingGame         *model.Game
	counterBetweenGames int
	gameRepository      *repository.GameRepository
}

func NewMainHandler(logger *logger.LoggerService, playerRepository *repository.PlayerRepository, gameRepository *repository.GameRepository) *MainHandler {
	mainHandler := &MainHandler{
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		},
		clientsInLobby:      make(map[uuid.UUID]*websocket.Conn),
		waitingGameClients:  make(map[uuid.UUID]*websocket.Conn),
		mutex:               &sync.Mutex{},
		broadcast:           make(chan []byte),
		logger:              logger,
		playerRepository:    playerRepository,
		ongoingGames:        []*model.Game{},
		pendingGame:         model.InitGame(),
		counterBetweenGames: 0,
		gameRepository:      gameRepository,
	}

	mainHandler.gameRepository.AddGame(mainHandler.pendingGame)

	return mainHandler
}

func (mh *MainHandler) wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := mh.upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error upgrading:", err)
		return
	}

	mh.mutex.Lock()
	mh.logger.Info("New client as connected with addr", "addr", conn.RemoteAddr())
	player := model.InitPlayer("test", conn)
	mh.playerRepository.AddPlayer(player)
	mh.clientsInLobby[player.ID] = conn
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
			// delete(mh.clientsInLobby, conn)
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
		case "game_subscription":
			fmt.Println("received game subscription evt")
			payload, err := ws_exchange.ExtractTypedPayload[ws_exchange.GameSubscriptionPayload](&exchangeRaw)
			if err != nil {
				fmt.Println("Payload ERROR:", err)
				continue
			}
			mh.handleGameSubscription(payload)
		case "pixel_click_evt":
			fmt.Println("Received pixel_click_evt")
			payload, err := ws_exchange.ExtractTypedPayload[ws_exchange.PixelClickPayload](&exchangeRaw)
			if err != nil {
				fmt.Println("Payload ERROR:", err)
				continue
			}
			mh.handlePixelEvt(payload)
		case "game_unsubscribe":
			fmt.Println("Received game_unsubscribe")
			payload, err := ws_exchange.ExtractTypedPayload[ws_exchange.GameUnsubscribePayload](&exchangeRaw)
			if err != nil {
				fmt.Println("Payload ERROR:", err)
				continue
			}
			mh.handleGameUnsubscribeEvt(payload)
		case "set_in_waiting_lobby":
			fmt.Println("Received set_in_waiting_lobby")
			payload, err := ws_exchange.ExtractTypedPayload[ws_exchange.SetInWaitingLobbyPayload](&exchangeRaw)
			if err != nil {
				fmt.Println("Payload ERROR:", err)
				continue
			}
			mh.handleSetInWaitingLobbyEvt(payload)
		case "exit_game":
			fmt.Println("Received game unsubscribed")
			payload, err := ws_exchange.ExtractTypedPayload[ws_exchange.ExitGamePayload](&exchangeRaw)
			if err != nil {
				fmt.Println("Payload ERROR:", err)
				continue
			}
			mh.handleExitGame(payload)
		}

		mh.mutex.Unlock()
	}
}

func (mh *MainHandler) handleSetInWaitingLobbyEvt(setInWaitingLobbyPayload *ws_exchange.SetInWaitingLobbyPayload) {
	parsedPlayerId, err := uuid.Parse(setInWaitingLobbyPayload.PlayerId)
	if err != nil {
		mh.logger.Error("Couldn't parse uuid from game subscription payload", "uuid", setInWaitingLobbyPayload.PlayerId, "err", err)
		return
	}

	player, err := mh.playerRepository.GetPlayerFromId(parsedPlayerId)
	if err != nil {
		return
	}

	mh.clientsInLobby[player.ID] = player.WsCon
}

func (mh *MainHandler) handleExitGame(exitGamePayload *ws_exchange.ExitGamePayload) {
	parsedPlayerId, err := uuid.Parse(exitGamePayload.PlayerId)
	if err != nil {
		mh.logger.Error("Couldn't parse uuid from game subscription payload", "uuid", exitGamePayload.PlayerId, "err", err)
		return
	}

	parsedGameId, err := uuid.Parse(exitGamePayload.GameId)
	if err != nil {
		mh.logger.Error("Couldn't parse uuid from game subscription payload", "uuid", exitGamePayload.GameId, "err", err)
		return
	}

	game, err := mh.gameRepository.GetGameOfId(parsedGameId)

	if err != nil {
		mh.logger.Error("Couldn't pase gameId", "err", err)
		return
	}

	err = game.RemovePlayer(parsedPlayerId)
	if err != nil {
		mh.logger.Error("Couldn't find player id ingame")
		return
	}
}

func (mh *MainHandler) handleGameUnsubscribeEvt(gameUnsubscribePayload *ws_exchange.GameUnsubscribePayload) {
	parsedUUid, err := uuid.Parse(gameUnsubscribePayload.PlayerId)
	if err != nil {
		mh.logger.Error("Couldn't parse uuid from game subscription payload", "uuid", gameUnsubscribePayload.PlayerId, "err", err)
		return
	}
	_, ok := mh.waitingGameClients[parsedUUid]
	if ok {
		mh.logger.Info("Unsubscribe player from game", "playerid", parsedUUid, "gameid", mh.pendingGame.ID)
		fmt.Println("NB PLAYER IN PENDING GAME BEFORE", len(mh.pendingGame.Players))
		err := mh.pendingGame.RemovePlayer(parsedUUid)
		if err != nil {
			return
		}
		fmt.Println("NB PLAYER IN PENDING GAME AFTER", len(mh.pendingGame.Players))
		delete(mh.waitingGameClients, parsedUUid)
		return
	}

	mh.logger.Error("Tried to unsubscribe unsubscribed player", "playerid", parsedUUid, "gameid", mh.pendingGame.ID)
}

func (mh *MainHandler) handleGameSubscription(gameSubscriptionPayload *ws_exchange.GameSubscriptionPayload) {
	parsedUUid, err := uuid.Parse(gameSubscriptionPayload.PlayerId)
	if err != nil {
		mh.logger.Error("Couldn't parse uuid from game subscription payload", "uuid", gameSubscriptionPayload.PlayerId, "err", err)
	}
	_, ok := mh.waitingGameClients[parsedUUid]
	if ok {
		mh.logger.Error("Player already subscribed to game", "playerid", parsedUUid, "gameid", mh.pendingGame.ID)
		return
	}

	mh.logger.Info("Subscribe player to game", "playerid", parsedUUid, "gameid", mh.pendingGame.ID)
	subscribingPlayer, err := mh.playerRepository.GetPlayerFromId(parsedUUid)
	mh.pendingGame.AddPlayer(subscribingPlayer)
	connCorrespondingToUuid := mh.clientsInLobby[parsedUUid]
	mh.waitingGameClients[parsedUUid] = connCorrespondingToUuid
}

func (mh *MainHandler) handlePixelEvt(pixelClickPayload *ws_exchange.PixelClickPayload) {
	gameId, err := uuid.Parse(pixelClickPayload.GameId)
	if err != nil {
		mh.logger.Error("Couldn't parse uuid from pixel click payload", "uuid", pixelClickPayload.GameId, "err", err)
		return
	}
	game, err := mh.gameRepository.GetGameOfId(gameId)
	if err != nil {
		mh.logger.Error("Couldn't find id from game", "uuid", pixelClickPayload.GameId, "err", err)
		return
	}

	game.ReceivePixelClick(pixelClickPayload)
}

func (mh *MainHandler) update() error {
	mh.logger.Info("UPDATE: ", "pending game id", mh.pendingGame.ID, "counter", mh.counterBetweenGames)
	mh.counterBetweenGames++
	if mh.counterBetweenGames == 10 && len(mh.pendingGame.Players) > 0 { // TODO => enelever hardcode
		mh.sendRedirectToGame()
		mh.pendingGame.Start()
		mh.ongoingGames = append(mh.ongoingGames, mh.pendingGame)
		mh.logger.Info("New game launched with ID", "id", mh.pendingGame.ID)
		mh.pendingGame = model.InitGame()
		mh.gameRepository.AddGame(mh.pendingGame)
		mh.waitingGameClients = make(map[uuid.UUID]*websocket.Conn)
		mh.counterBetweenGames = 0
		return nil
	} else if mh.counterBetweenGames == 10 && len(mh.pendingGame.Players) == 0 {
		mh.counterBetweenGames = 0
		return nil
	}

	err := mh.sendPendingGameUpdate()
	if err != nil {
		return err
	}

	return nil
}

func (mh *MainHandler) sendRedirectToGame() error {
	data := &ws_exchange.RedirectToGamePayload{
		GameId: mh.pendingGame.ID.String(),
	}

	bytes, err := json.Marshal(data.ToWsExchange())
	if err != nil {
		return err
	}

	for client := range mh.waitingGameClients {
		err = mh.waitingGameClients[client].WriteMessage(websocket.TextMessage, bytes)
		delete(mh.waitingGameClients, client)
		delete(mh.clientsInLobby, client)
		if err != nil {
			return err
		}
	}

	return nil
}

func (mh *MainHandler) sendPendingGameUpdate() error {
	for client := range mh.clientsInLobby {
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
			IsGameLaunching:        mh.counterBetweenGames == 0,
		}

		bytes, err := json.Marshal(data.ToWsExchange())
		if err != nil {
			return err
		}

		err = mh.clientsInLobby[client].WriteMessage(websocket.TextMessage, bytes)
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
				mh.update()
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
