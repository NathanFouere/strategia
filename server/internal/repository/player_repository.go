package repository

import (
	"errors"
	"server/internal/model"
	"server/pkg/logger"

	"github.com/google/uuid"
)

type PlayerRepository struct {
	LoggedPlayers []*model.Player // Solution temporaire ou on stocke tt les joueurs en mémoire
	logger        *logger.LoggerService
}

func NewPlayerRepository(loggerService *logger.LoggerService) *PlayerRepository {
	return &PlayerRepository{
		LoggedPlayers: []*model.Player{},
		logger:        loggerService,
	}
}

func (pr *PlayerRepository) AddPlayer(pseudo string) {
	newPlayer := model.InitPlayer(pseudo)
	pr.LoggedPlayers = append(pr.LoggedPlayers, newPlayer)
}

func (pr *PlayerRepository) RemovePlayer(uuid uuid.UUID) error {
	// moche et gourmand mais marche pour le moment
	for i := 0; i < len(pr.LoggedPlayers); i++ {
		if pr.LoggedPlayers[i].ID == uuid {
			// cf . https://stackoverflow.com/questions/37334119/how-to-delete-an-element-from-a-slice-in-golang
			pr.LoggedPlayers[i] = pr.LoggedPlayers[len(pr.LoggedPlayers)-1]
			return nil
		}
	}

	pr.logger.Error("Didn't find the player of uuid", "uuid", uuid)
	return errors.New("Didnt find player")
}
