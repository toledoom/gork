package usecases

import (
	"github.com/toledoom/gork/internal/app/query"
	"github.com/toledoom/gork/internal/domain/leaderboard"
	"github.com/toledoom/gork/pkg/gork"
)

type GetTopPlayersInput struct {
	NumPlayers int64
}

type GetTopPlayersOutput struct {
	MemberList []*leaderboard.Member
}

func GetTopPlayers(cr *gork.CommandRegistry, qr *gork.QueryRegistry) gork.UseCase[GetTopPlayersInput, GetTopPlayersOutput] {
	return func(gtpi GetTopPlayersInput) (GetTopPlayersOutput, error) {
		q := &query.GetTopPlayers{
			NumPlayers: gtpi.NumPlayers,
		}

		response, err := gork.HandleQuery[*query.GetTopPlayers, *query.GetTopPlayersResponse](qr, q)
		if err != nil {
			return GetTopPlayersOutput{}, nil
		}

		return GetTopPlayersOutput{
			MemberList: response.MemberList,
		}, nil
	}
}
