package tmp

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

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type MainHandler struct {
	upgrader             websocket.Upgrader
	mutex                *sync.Mutex
	broadcast            chan []byte
	logger               *logger.LoggerService
	playerRepository     *repository.PlayerRepository
	ongoingGames         []*model.Game
	counterBetweenGames  int
	gameRepository       *repository.GameRepository
	messageRouterService *service.MessageRouterService
}

func NewMainHandler(
	logger *logger.LoggerService,
	playerRepository *repository.PlayerRepository,
	gameRepository *repository.GameRepository,
	messageRouterService *service.MessageRouterService,
) *MainHandler {
	mainHandler := &MainHandler{
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		},
		mutex:                &sync.Mutex{},
		broadcast:            make(chan []byte),
		logger:               logger,
		playerRepository:     playerRepository,
		ongoingGames:         []*model.Game{},
		counterBetweenGames:  0,
		gameRepository:       gameRepository,
		messageRouterService: messageRouterService,
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
	mh.playerRepository.ClientsInLobby[player.ID] = conn
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

	mh.playerRepository.ClientsInLobby[player.ID] = player.WsCon
}

func (mh *MainHandler) handleUpdatePlayerPseudo(updatePlayerPseudoPayload *ws_exchange.UpdatePlayerPseudoPayload) {
	parsedPlayerId, err := uuid.Parse(updatePlayerPseudoPayload.PlayerId)
	if err != nil {
		mh.logger.Error("Couldn't parse uuid from game subscription payload", "uuid", updatePlayerPseudoPayload.PlayerId, "err", err)
	}

	player, err := mh.playerRepository.GetPlayerFromId(parsedPlayerId)
	if err != nil {
		return
	}

	player.UpdatePseudo(updatePlayerPseudoPayload.NewPseudo)
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

func (mh *MainHandler) update() error {
	mh.logger.Info("UPDATE: ", "pending game id", mh.gameRepository.PendingGame.ID, "counter", mh.counterBetweenGames)
	mh.counterBetweenGames++
	if mh.counterBetweenGames == 10 && len(mh.gameRepository.PendingGame.Players) > 0 { // TODO => enelever hardcode
		mh.sendRedirectToGame()
		mh.gameRepository.PendingGame.Start()
		mh.ongoingGames = append(mh.ongoingGames, mh.gameRepository.PendingGame)
		mh.logger.Info("New game launched with ID", "id", mh.gameRepository.PendingGame.ID)
		mh.gameRepository.PendingGame = model.InitGame()
		mh.gameRepository.AddGame(mh.gameRepository.PendingGame)
		mh.playerRepository.WaitingGameClients = make(map[uuid.UUID]*websocket.Conn)
		mh.counterBetweenGames = 0
		return nil
	} else if mh.counterBetweenGames == 10 && len(mh.gameRepository.PendingGame.Players) == 0 {
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
		GameId: mh.gameRepository.PendingGame.ID.String(),
	}

	bytes, err := json.Marshal(data.ToWsExchange())
	if err != nil {
		return err
	}

	for client := range mh.playerRepository.WaitingGameClients {
		err = mh.playerRepository.WaitingGameClients[client].WriteMessage(websocket.TextMessage, bytes)
		delete(mh.playerRepository.WaitingGameClients, client)
		delete(mh.playerRepository.ClientsInLobby, client)
		if err != nil {
			return err
		}
	}

	return nil
}

func (mh *MainHandler) sendPendingGameUpdate() error {
	for client := range mh.playerRepository.ClientsInLobby {
		isClientWaitingForGame := false
		_, ok := mh.playerRepository.WaitingGameClients[client]
		if ok {
			isClientWaitingForGame = true
		}
		data := &ws_exchange.WaitingGamePayload{
			SecondsBeforeLaunch:    10 - mh.counterBetweenGames, // TODO => enlever hardcode
			GameId:                 mh.gameRepository.PendingGame.ID.String(),
			NumberOfWaitingPlayers: len(mh.playerRepository.WaitingGameClients),
			IsPlayerWaitingForGame: isClientWaitingForGame,
			IsGameLaunching:        mh.counterBetweenGames == 0,
		}

		bytes, err := json.Marshal(data.ToWsExchange())
		if err != nil {
			return err
		}

		err = mh.playerRepository.ClientsInLobby[client].WriteMessage(websocket.TextMessage, bytes)
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
