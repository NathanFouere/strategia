package repository

import (
	"errors"
	"server/internal/model"
	"server/pkg/logger"

	"github.com/google/uuid"
)

type PlayerRepository struct {
	players            []*model.Player
	logger             *logger.LoggerService
	ClientsInLobby     map[uuid.UUID]*model.Player
	WaitingGameClients map[uuid.UUID]*model.Player
}

func NewPlayerRepository(loggerService *logger.LoggerService) *PlayerRepository {
	return &PlayerRepository{
		players:            []*model.Player{},
		logger:             loggerService,
		ClientsInLobby:     make(map[uuid.UUID]*model.Player),
		WaitingGameClients: make(map[uuid.UUID]*model.Player),
	}
}

func (pr *PlayerRepository) AddPlayer(player *model.Player) {
	pr.players = append(pr.players, player)
}

func (pr *PlayerRepository) GetPlayerFromId(uuid uuid.UUID) (*model.Player, error) {
	for i := 0; i < len(pr.players); i++ {
		if pr.players[i].ID == uuid {
			return pr.players[i], nil
		}
	}

	return nil, errors.New("Couldn't find player of id " + uuid.String())
}

func (pr *PlayerRepository) AddPlayerToClientLobby(p *model.Player) {
	pr.ClientsInLobby[p.ID] = p
}

func (pr *PlayerRepository) RemovePlayer(uuid uuid.UUID) error {
	for i := 0; i < len(pr.players); i++ {
		if pr.players[i].ID == uuid {
			// cf . https://stackoverflow.com/questions/37334119/how-to-delete-an-element-from-a-slice-in-golang
			pr.players[i] = pr.players[len(pr.players)-1]
			return nil
		}
	}

	pr.logger.Error("Didn't find the player of uuid", "uuid", uuid)
	return errors.New("Didn't find player")
}
