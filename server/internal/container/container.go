//nolint:sloglint // usage du logger global accepté ici car c'est l'initialisation du container
package container

import (
	"log/slog"
	"server/internal/handler"
	"server/internal/repository"
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

	err = GetContainer().Provide(func(logger *logger.LoggerService, playerRepository *repository.PlayerRepository) *service.GameGeneratorService {
		return service.NewGameGeneratorService(logger, playerRepository)
	})
	if err != nil {
		slog.Error("Error occured while providing game generator service", "err", err)
	}

	err = GetContainer().Provide(func(logger *logger.LoggerService, playerRepository *repository.PlayerRepository) *handler.MainHandler {
		return handler.NewMainHandler(logger, playerRepository)
	})
	if err != nil {
		slog.Error("Error occured while providing main handler", "err", err)
	}

	slog.Info("Container successfully initiated !")
	return err
}
