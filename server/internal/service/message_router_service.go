package service

import (
	"fmt"
	"server/internal/handler"
	"server/internal/ws_exchange"
	"server/pkg/logger"
)

type MessageRouterService struct {
	logger                     *logger.LoggerService
	exitGameHandler            *handler.ExitGameHandler
	pixelClickHandler          *handler.PixelClickHandler
	setInWaitingLobbyHandler   *handler.SetInWaitingLobbyHandler
	subscribeToGameHandler     *handler.SubscribeToGameHandler
	unsubscribeFromGameHandler *handler.UnsubscribeFromGameHandler
	updatePlayerPseudoHandler  *handler.UpdatePlayerPseudoHandler
}

func NewMessageRouterService(
	logger *logger.LoggerService,
	exitGameHandler *handler.ExitGameHandler,
	pixelClickHandler *handler.PixelClickHandler,
	setInWaitingLobbyHandler *handler.SetInWaitingLobbyHandler,
	subscribeToGameHandler *handler.SubscribeToGameHandler,
	unsubscribeFromGameHandler *handler.UnsubscribeFromGameHandler,
	updatePlayerPseudoHandler *handler.UpdatePlayerPseudoHandler,
) *MessageRouterService {
	return &MessageRouterService{
		logger:                     logger,
		exitGameHandler:            exitGameHandler,
		pixelClickHandler:          pixelClickHandler,
		setInWaitingLobbyHandler:   setInWaitingLobbyHandler,
		subscribeToGameHandler:     subscribeToGameHandler,
		unsubscribeFromGameHandler: unsubscribeFromGameHandler,
		updatePlayerPseudoHandler:  updatePlayerPseudoHandler,
	}
}

func (s *MessageRouterService) HandleMessage(exchangeRaw ws_exchange.WsExchangeTemplateRaw) {
	switch exchangeRaw.Type {
	case "game_subscription":
		fmt.Println("received game subscription evt")
		payload, err := ws_exchange.ExtractTypedPayload[ws_exchange.GameSubscriptionPayload](&exchangeRaw)
		if err != nil {
			fmt.Println("Payload ERROR:", err)
			return
		}
		s.subscribeToGameHandler.Handle(payload)
	case "pixel_click_evt":
		fmt.Println("Received pixel_click_evt")
		payload, err := ws_exchange.ExtractTypedPayload[ws_exchange.PixelClickPayload](&exchangeRaw)
		if err != nil {
			fmt.Println("Payload ERROR:", err)
			return
		}
		s.pixelClickHandler.Handle(payload)
	case "game_unsubscribe":
		fmt.Println("Received game_unsubscribe")
		payload, err := ws_exchange.ExtractTypedPayload[ws_exchange.GameUnsubscribePayload](&exchangeRaw)
		if err != nil {
			fmt.Println("Payload ERROR:", err)
			return
		}
		s.unsubscribeFromGameHandler.Handle(payload)
	case "set_in_waiting_lobby":
		fmt.Println("Received set_in_waiting_lobby")
		payload, err := ws_exchange.ExtractTypedPayload[ws_exchange.SetInWaitingLobbyPayload](&exchangeRaw)
		if err != nil {
			fmt.Println("Payload ERROR:", err)
			return
		}
		s.setInWaitingLobbyHandler.Handle(payload)
	case "exit_game":
		fmt.Println("Received game unsubscribed")
		payload, err := ws_exchange.ExtractTypedPayload[ws_exchange.ExitGamePayload](&exchangeRaw)
		if err != nil {
			fmt.Println("Payload ERROR:", err)
			return
		}
		s.exitGameHandler.Handle(payload)
	case "update_player_pseudo":
		fmt.Println("Received update player pseudo")
		payload, err := ws_exchange.ExtractTypedPayload[ws_exchange.UpdatePlayerPseudoPayload](&exchangeRaw)
		if err != nil {
			fmt.Println("Payload ERROR:", err)
			return
		}
		s.updatePlayerPseudoHandler.Handle(payload)
	}
}
