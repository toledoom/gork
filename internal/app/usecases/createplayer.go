package usecases

import (
	"github.com/toledoom/gork/internal/app/command"
	"github.com/toledoom/gork/internal/app/query"
	"github.com/toledoom/gork/internal/domain/player"
	"github.com/toledoom/gork/pkg/gork"
)

type CreatePlayerInput struct {
	PlayerID, Name string
}

type CreatePlayerOutput struct {
	Player player.Player
}

func CreatePlayer(cr *gork.CommandRegistry, qr *gork.QueryRegistry) gork.UseCase[CreatePlayerInput, CreatePlayerOutput] {
	return func(cpi CreatePlayerInput) (CreatePlayerOutput, error) {
		createPlayerCommand := command.CreatePlayer{
			PlayerID: cpi.PlayerID,
			Name:     cpi.Name,
		}

		err := gork.HandleCommand(cr, &createPlayerCommand)
		if err != nil {
			return CreatePlayerOutput{}, err
		}

		getPlayerQuery := query.GetPlayerByID{
			PlayerID: cpi.PlayerID,
		}

		queryResponse, err := gork.HandleQuery[*query.GetPlayerByID, *query.GetPlayerByIDResponse](qr, &getPlayerQuery)
		if err != nil {
			return CreatePlayerOutput{}, err
		}

		return CreatePlayerOutput{Player: *queryResponse.Player}, nil
	}
}
