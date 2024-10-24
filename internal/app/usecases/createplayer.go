package usecases

import (
	"github.com/toledoom/gork/internal/app/command"
	"github.com/toledoom/gork/internal/app/query"
	"github.com/toledoom/gork/internal/domain/player"
	"github.com/toledoom/gork/pkg/gork/cqrs"
)

type CreatePlayerUseCase struct {
	cr *cqrs.CommandRegistry
	qr *cqrs.QueryRegistry
}

type CreatePlayerInput struct {
	PlayerID, Name string
}

type CreatePlayerOutput struct {
	Player player.Player
}

func (cpuc *CreatePlayerUseCase) Execute(cpi CreatePlayerInput) (CreatePlayerOutput, error) {
	createPlayerCommand := command.CreatePlayer{
		PlayerID: cpi.PlayerID,
		Name:     cpi.Name,
	}

	err := cqrs.HandleCommand(cpuc.cr, &createPlayerCommand)
	if err != nil {
		return CreatePlayerOutput{}, err
	}

	getPlayerQuery := query.GetPlayerByID{
		PlayerID: cpi.PlayerID,
	}

	queryResponse, err := cqrs.HandleQuery[*query.GetPlayerByID, *query.GetPlayerByIDResponse](cpuc.qr, &getPlayerQuery)
	if err != nil {
		return CreatePlayerOutput{}, err
	}

	return CreatePlayerOutput{Player: *queryResponse.Player}, nil
}
