package service

import (
	"errors"
	"server/internal/model"
	"server/internal/repository"
	"server/pkg/logger"

	"github.com/google/uuid"
)

type GameGeneratorService struct {
	CooldownSinceLastGame int
	AttendingPlayers      []*model.Player
	logger                *logger.LoggerService
	playerRepository      *repository.PlayerRepository
}

func NewGameGeneratorService(loggerService *logger.LoggerService, playerRepository *repository.PlayerRepository) *GameGeneratorService {
	return &GameGeneratorService{
		CooldownSinceLastGame: 0,
		AttendingPlayers:      []*model.Player{},
		logger:                loggerService,
		playerRepository:      playerRepository,
	}
}

func (gcs *GameGeneratorService) AddAttendingPlayer(player *model.Player) {
	gcs.AttendingPlayers = append(gcs.AttendingPlayers, player)
}

func (gcs *GameGeneratorService) RemoveAttendingPlayer(uuid uuid.UUID) error {
	for i := 0; i < len(gcs.AttendingPlayers); i++ {
		if gcs.AttendingPlayers[i].ID == uuid {
			// cf . https://stackoverflow.com/questions/37334119/how-to-delete-an-element-from-a-slice-in-golang
			gcs.AttendingPlayers[i] = gcs.AttendingPlayers[len(gcs.AttendingPlayers)-1]
			gcs.AttendingPlayers = gcs.AttendingPlayers[:len(gcs.AttendingPlayers)-1]

			return nil
		}
	}

	gcs.logger.Error("Didn't find the player of uuid", "uuid", uuid)
	return errors.New("Didnt find player")
}

func (gcs *GameGeneratorService) StartGame() {
	game := model.InitGame()
	for _, attendingPlayer := range gcs.AttendingPlayers {
		game.AddPlayer(attendingPlayer)
	}
	gcs.AttendingPlayers = []*model.Player{}
}

func (gcs *GameGeneratorService) Update() {
	gcs.CooldownSinceLastGame++
	if gcs.CooldownSinceLastGame > 10 {
		gcs.StartGame()
		gcs.CooldownSinceLastGame = 0
	}
}
