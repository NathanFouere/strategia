package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"server/internal/model"
	"server/internal/repository"
	"server/internal/sender"
	"server/internal/service"
	"server/internal/ws_exchange"
	"server/pkg/logger"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

type MainHandler struct {
	upgrader                 websocket.Upgrader
	broadcast                chan []byte
	logger                   *logger.LoggerService
	playerRepository         *repository.PlayerRepository
	gameRepository           *repository.GameRepository
	messageRouterService     *service.MessageRouterService
	updateService            *service.UpdateService
	connectionExchangeSender *sender.ConnectionExchangeSender
}

func NewMainHandler(
	logger *logger.LoggerService,
	playerRepository *repository.PlayerRepository,
	gameRepository *repository.GameRepository,
	messageRouterService *service.MessageRouterService,
	updateService *service.UpdateService,
	connectionExchangeSender *sender.ConnectionExchangeSender,
) *MainHandler {
	mainHandler := &MainHandler{
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		},
		broadcast:                make(chan []byte),
		logger:                   logger,
		playerRepository:         playerRepository,
		gameRepository:           gameRepository,
		messageRouterService:     messageRouterService,
		updateService:            updateService,
		connectionExchangeSender: connectionExchangeSender,
	}

	pending := model.InitGame()
	gameRepository.PendingGame = pending
	gameRepository.AddGame(pending)

	return mainHandler
}

func (mh *MainHandler) wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := mh.upgrader.Upgrade(w, r, nil)
	if err != nil {
		mh.logger.Error("Error while upgrading ws connection", err)
		return
	}

	mh.logger.Info("New client has connected with addr", "addr", conn.RemoteAddr())
	player := model.InitPlayer(conn)
	mh.playerRepository.AddPlayer(player)
	mh.playerRepository.ClientsInLobby[player.ID] = player

	go func() {
		for msg := range player.Client.Send {
			err := player.Client.Conn.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				mh.logger.Error("Error when sending messsage to player", "player id", player.ID, "error", err)
				return
			}
		}
	}()

	mh.connectionExchangeSender.Send(player)

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			mh.playerRepository.RemovePlayer(player.ID)
			break
		}
		mh.broadcast <- message
	}
}

func (mh *MainHandler) handleMessages() {
	for {
		message := <-mh.broadcast

		var exchangeRaw ws_exchange.WsExchangeTemplateRaw
		err := json.Unmarshal(message, &exchangeRaw)
		if err != nil {
			mh.logger.Error("Error while unmarshalling message", "message", message, "error", err)
			continue
		}
		mh.messageRouterService.HandleMessage(exchangeRaw)

	}
}

func (mh *MainHandler) Launch() {
	tickerMainMenuSecond, err := strconv.Atoi(os.Getenv("TICKER_MAIN_MENU_SECONDS"))
	if err != nil {
		panic("Couldn't load TICKER_MAIN_MENU_SECONDS in .env")
	}

	ticker := time.NewTicker(time.Duration(tickerMainMenuSecond) * time.Second)
	go func() {
		for range ticker.C {
			mh.updateService.Update()
		}
	}()
	http.HandleFunc("/ws", mh.wsHandler)
	go mh.handleMessages()
	mh.logger.Info("Starting server on port", "port", os.Getenv("HTTP_PORT"))
	err = http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("HTTP_PORT")), nil)
	if err != nil {
		mh.logger.Error("Error starting server", "error", err)
		panic("Error starting server")
	}
}
