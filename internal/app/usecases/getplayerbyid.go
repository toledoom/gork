package usecases

import (
	"github.com/toledoom/gork/internal/app/query"
	"github.com/toledoom/gork/internal/domain/player"
	"github.com/toledoom/gork/pkg/gork/cqrs"
)

type GetPlayerByIDUseCase struct {
	qr *cqrs.QueryRegistry
}

type GetPlayerByIDInput struct {
	PlayerID string
}

type GetPlayerByIDOutput struct {
	Player *player.Player
}

func (gpbiuc *GetPlayerByIDUseCase) Execute(gpbii GetPlayerByIDInput) (GetPlayerByIDOutput, error) {
	q := query.GetPlayerByID{
		PlayerID: gpbii.PlayerID,
	}
	response, err := cqrs.HandleQuery[*query.GetPlayerByID, *query.GetPlayerByIDResponse](gpbiuc.qr, &q)
	if err != nil {
		return GetPlayerByIDOutput{}, err
	}

	return GetPlayerByIDOutput{
		Player: response.Player,
	}, nil
}
