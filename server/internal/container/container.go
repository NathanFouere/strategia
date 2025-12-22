//nolint:sloglint // usage du logger global accepté ici car c'est l'initialisation du container
package container

import (
	"log/slog"
	"server/internal"
	"server/internal/broadcaster"
	"server/internal/handler"
	"server/internal/repository"
	"server/internal/sender"
	"server/internal/service"
	"server/pkg/logger"
	"sync"

	"go.uber.org/dig"
)

var singleton *dig.Container //nolint:gochecknoglobals // on accepte la variable globale singleton pour le container
var once sync.Once           //nolint:gochecknoglobals // on accepte la variable globale singleton pour le container

func GetContainer() *dig.Container {
	once.Do(func() {
		singleton = dig.New()
	})
	return singleton
}

func SetupContainer() error {
	err := GetContainer().Provide(slog.Default)
	if err != nil {
		slog.Error("Error occurred while providing root slog.Logger", "err", err)
		return err
	}

	err = GetContainer().Provide(logger.NewLoggerService)
	if err != nil {
		slog.Error("Error occured while providing logger service", "err", err)
	}

	err = GetContainer().Provide(func(logger *logger.LoggerService) *repository.PlayerRepository {
		return repository.NewPlayerRepository(logger)
	})
	if err != nil {
		slog.Error("Error occured while providing player repository", "err", err)
	}

	err = GetContainer().Provide(func(logger *logger.LoggerService) *repository.GameRepository {
		return repository.NewGameRepository(logger)
	})
	if err != nil {
		slog.Error("Error occured while providing game repository", "err", err)
	}

	err = GetContainer().Provide(func(logger *logger.LoggerService, playerRepository *repository.PlayerRepository) *service.GameGeneratorService {
		return service.NewGameGeneratorService(logger, playerRepository)
	})
	if err != nil {
		slog.Error("Error occured while providing game generator service", "err", err)
	}

	err = GetContainer().Provide(func(logger *logger.LoggerService, playerRepository *repository.PlayerRepository, gameRepository *repository.GameRepository) *handler.ExitGameHandler {
		return handler.NewExitGameHandler(playerRepository, gameRepository, logger)
	})
	if err != nil {
		slog.Error("Error occured while providing exit game handler", "err", err)
	}

	err = GetContainer().Provide(func(logger *logger.LoggerService, playerRepository *repository.PlayerRepository, gameRepository *repository.GameRepository) *handler.PixelClickHandler {
		return handler.NewPixelClickHandler(playerRepository, gameRepository, logger)
	})
	if err != nil {
		slog.Error("Error occured while providing pixel click handler", "err", err)
	}

	err = GetContainer().Provide(func(logger *logger.LoggerService, playerRepository *repository.PlayerRepository) *handler.SetInWaitingLobbyHandler {
		return handler.NewSetInWaitingLobbyHandler(logger, playerRepository)
	})
	if err != nil {
		slog.Error("Error occured while providing set in waiting lobby handler", "err", err)
	}

	err = GetContainer().Provide(func(logger *logger.LoggerService, playerRepository *repository.PlayerRepository, gameRepository *repository.GameRepository) *handler.SubscribeToGameHandler {
		return handler.NewSubscribeToGameHandler(playerRepository, gameRepository, logger)
	})
	if err != nil {
		slog.Error("Error occured while providing subscribe to game handler", "err", err)
	}

	err = GetContainer().Provide(func(logger *logger.LoggerService, playerRepository *repository.PlayerRepository, gameRepository *repository.GameRepository) *handler.UnsubscribeFromGameHandler {
		return handler.NewUnsubscribeFromGameHandler(logger, playerRepository, gameRepository)
	})
	if err != nil {
		slog.Error("Error occured while providing unsubscribe to game handler", "err", err)
	}

	err = GetContainer().Provide(func(logger *logger.LoggerService, playerRepository *repository.PlayerRepository) *handler.UpdatePlayerPseudoHandler {
		return handler.NewUpdatePlayerPseudoHandler(playerRepository, logger)
	})
	if err != nil {
		slog.Error("Error occured while providing unsubscribe to game handler", "err", err)
	}

	err = GetContainer().Provide(func(
		logger *logger.LoggerService,
		exitGameHandler *handler.ExitGameHandler,
		clickHandler *handler.PixelClickHandler,
		setInWaitingLobbyHandler *handler.SetInWaitingLobbyHandler,
		toGameHandler *handler.SubscribeToGameHandler,
		fromGameHandler *handler.UnsubscribeFromGameHandler,
		pseudoPayload *handler.UpdatePlayerPseudoHandler,
	) *service.MessageRouterService {
		return service.NewMessageRouterService(
			logger,
			exitGameHandler,
			clickHandler,
			setInWaitingLobbyHandler,
			toGameHandler,
			fromGameHandler,
			pseudoPayload,
		)
	})
	if err != nil {
		slog.Error("Error occured while providing message router service", "err", err)
	}

	err = GetContainer().Provide(func(
		logger *logger.LoggerService,
		pr *repository.PlayerRepository,
		gr *repository.GameRepository,
	) *sender.PendingGameUpdateSender {
		return sender.NewPendingGameUpdateSender(
			logger,
			pr,
			gr,
		)
	})
	if err != nil {
		slog.Error("Error occured while providing pending game update sender", "err", err)
	}

	err = GetContainer().Provide(func(
		logger *logger.LoggerService,
		pr *repository.PlayerRepository,
		gr *repository.GameRepository,
	) *sender.RedirectToGameSender {
		return sender.NewRedirectToGameSender(
			pr,
			gr,
			logger,
		)
	})
	if err != nil {
		slog.Error("Error occured while providing pending redirect to game sender", "err", err)
	}

	err = GetContainer().Provide(func(
		logger *logger.LoggerService,
		updateGameService *broadcaster.GameUpdateBroadcaster,
	) *service.StartGameService {
		return service.NewStartGameService(
			logger,
			updateGameService,
		)
	})
	if err != nil {
		slog.Error("Error occured while providing pending start game service", "err", err)
	}

	err = GetContainer().Provide(func(
		logger *logger.LoggerService,
		pr *repository.PlayerRepository,
		gr *repository.GameRepository,
		redirectToGameSender *sender.RedirectToGameSender,
		pendingGameUpdateSender *sender.PendingGameUpdateSender,
		startGameService *service.StartGameService,
	) *service.UpdateService {
		return service.NewUpdateService(
			pr,
			gr,
			logger,
			pendingGameUpdateSender,
			redirectToGameSender,
			startGameService,
		)
	})
	if err != nil {
		slog.Error("Error occured while providing update service", "err", err)
	}

	err = GetContainer().Provide(func(
		logger *logger.LoggerService,
	) *sender.ConnectionExchangeSender {
		return sender.NewConnectionExchangeSender(logger)
	})
	if err != nil {
		slog.Error("Error occured while providing connection exchange sender", "err", err)
	}

	err = GetContainer().Provide(func(
		logger *logger.LoggerService,
	) *broadcaster.GameUpdateBroadcaster {
		return broadcaster.NewGameUpdateBroadcaster(logger)
	})

	err = GetContainer().Provide(func(
		logger *logger.LoggerService,
		playerRepository *repository.PlayerRepository,
		gameRepository *repository.GameRepository,
		messageRouterService *service.MessageRouterService,
		updateService *service.UpdateService,
		connectionExchangeSender *sender.ConnectionExchangeSender,
	) *internal.MainHandler {
		return internal.NewMainHandler(logger, playerRepository, gameRepository, messageRouterService, updateService, connectionExchangeSender)
	})
	if err != nil {
		slog.Error("Error occured while providing main handler", "err", err)
	}

	slog.Info("Container successfully initiated !")
	return err
}
