package query

import (
	"github.com/toledoom/gork/internal/domain/leaderboard"
)

const GetRankQueryID = "GetRank"

type GetRank struct {
	PlayerID string
}

type GetRankResponse struct {
	Rank uint64
}

func GetRankHandler(ranking leaderboard.Ranking) func(*GetRank) (*GetRankResponse, error) {
	return func(q *GetRank) (*GetRankResponse, error) {
		playerID := q.PlayerID

		rank, err := ranking.GetRank(playerID)
		if err != nil {
			return nil, err
		}

		return &GetRankResponse{
			Rank: rank,
		}, nil
	}
}
