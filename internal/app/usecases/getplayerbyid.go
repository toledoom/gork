package usecases

import (
	"github.com/toledoom/gork/internal/app/query"
	"github.com/toledoom/gork/internal/domain/player"
	"github.com/toledoom/gork/pkg/gork"
)

type GetPlayerByIDInput struct {
	PlayerID string
}

type GetPlayerByIDOutput struct {
	Player *player.Player
}

func GetPlayerByID(qr *gork.QueryRegistry) func(gpbid GetPlayerByIDInput) (GetPlayerByIDOutput, error) {
	return func(gpbid GetPlayerByIDInput) (GetPlayerByIDOutput, error) {
		q := query.GetPlayerByID{
			PlayerID: gpbid.PlayerID,
		}
		response, err := gork.HandleQuery[*query.GetPlayerByID, *query.GetPlayerByIDResponse](qr, &q)
		if err != nil {
			return GetPlayerByIDOutput{}, err
		}

		return GetPlayerByIDOutput{
			Player: response.Player,
		}, nil
	}
}
