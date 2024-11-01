package usecases

import (
	"github.com/toledoom/gork/internal/app/query"
	"github.com/toledoom/gork/pkg/gork/cqrs"
)

type GetRankInput struct {
	PlayerID string
}

type GetRankOutput struct {
	Rank uint64
}

func GetRank(qr *cqrs.QueryRegistry) func(gri GetRankInput) (GetRankOutput, error) {
	return func(gri GetRankInput) (GetRankOutput, error) {
		getRankQuery := query.GetRank{
			PlayerID: gri.PlayerID,
		}

		queryResponse, err := cqrs.HandleQuery[*query.GetRank, *query.GetRankResponse](qr, &getRankQuery)
		if err != nil {
			return GetRankOutput{}, err
		}

		return GetRankOutput{
			Rank: queryResponse.Rank,
		}, nil

	}
}
