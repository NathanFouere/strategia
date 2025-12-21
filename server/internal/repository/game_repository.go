package repository

import (
	"errors"
	"server/internal/model"
	"server/pkg/logger"

	"github.com/google/uuid"
)

type GameRepository struct {
	games               []*model.Game
	logger              *logger.LoggerService
	PendingGame         *model.Game
	OngoingGames        []*model.Game
	CounterBetweenGames int // TODO => move
}

func NewGameRepository(loggerService *logger.LoggerService) *GameRepository {
	return &GameRepository{
		games:               []*model.Game{},
		logger:              loggerService,
		PendingGame:         nil,
		CounterBetweenGames: 0,
		OngoingGames:        []*model.Game{},
	}
}

func (gr *GameRepository) AddGame(game *model.Game) {
	gr.games = append(gr.games, game)
}

func (gr *GameRepository) GetGameOfId(gameId uuid.UUID) (*model.Game, error) {
	for i := 0; i < len(gr.games); i++ {
		if gr.games[i].ID == gameId {
			return gr.games[i], nil
		}
	}

	gr.logger.Error("Didn't find the game of id", "id", gameId)
	return nil, errors.New("Didn't find game")
}

func (gr *GameRepository) RemoveGame(gameId uuid.UUID) error {
	for i := 0; i < len(gr.games); i++ {
		if gr.games[i].ID == gameId {
			// cf . https://stackoverflow.com/questions/37334119/how-to-delete-an-element-from-a-slice-in-golang
			gr.games[i] = gr.games[len(gr.games)-1]
			gr.games = gr.games[:len(gr.games)-1]
			return nil
		}
	}

	gr.logger.Error("Didn't find the game of id", "id", gameId)
	return errors.New("Didn't find game")
}
