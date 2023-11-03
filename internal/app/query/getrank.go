package query

import (
	"fmt"

	"github.com/toledoom/gork/internal/domain/leaderboard"
	"github.com/toledoom/gork/pkg/cqrs"
)

const GetRankQueryID = "GetRank"

type GetRank struct {
	PlayerID string
}

func (grq *GetRank) QueryID() string {
	return GetRankQueryID
}

type GetRankHandler struct {
	ranking leaderboard.Ranking
}

func NewGetRankHandler(ranking leaderboard.Ranking) *GetRankHandler {
	return &GetRankHandler{
		ranking: ranking,
	}
}

func (grh *GetRankHandler) QueryID() string {
	return GetRankQueryID
}

func (grh *GetRankHandler) Handle(q cqrs.Query) (cqrs.QueryResponse[any], error) {
	grq, ok := q.(*GetRank)
	if !ok {
		return nil, fmt.Errorf("wrong query: %v", q)
	}
	playerID := grq.PlayerID

	rank, err := grh.ranking.GetRank(playerID)
	if err != nil {
		return nil, err
	}

	return &GetRankResponse{
		Rank: rank,
	}, nil
}

type GetRankResponse struct {
	Rank int64
}

func (grr *GetRankResponse) Data() any {
	return grr
}
