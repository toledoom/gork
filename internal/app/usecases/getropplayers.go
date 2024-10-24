package usecases

import (
	"github.com/toledoom/gork/internal/app/query"
	"github.com/toledoom/gork/internal/domain/leaderboard"
	"github.com/toledoom/gork/pkg/gork/cqrs"
)

type GetTopPlayersUseCase struct {
	qr *cqrs.QueryRegistry
}

type GetTopPlayersInput struct {
	NumPlayers int64
}

type GetTopPlayersOutput struct {
	MemberList []*leaderboard.Member
}

func (gtpuc *GetTopPlayersUseCase) Execute(gtpi GetTopPlayersInput) (GetTopPlayersOutput, error) {
	q := &query.GetTopPlayers{
		NumPlayers: gtpi.NumPlayers,
	}

	response, err := cqrs.HandleQuery[*query.GetTopPlayers, *query.GetTopPlayersResponse](gtpuc.qr, q)
	if err != nil {
		return GetTopPlayersOutput{}, nil
	}

	return GetTopPlayersOutput{
		MemberList: response.MemberList,
	}, nil
}
