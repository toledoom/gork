package query

import (
	"fmt"

	"github.com/toledoom/gork/internal/domain/leaderboard"
	"github.com/toledoom/gork/pkg/cqrs"
)

const GetTopPlayersQueryID = "GetTopPlayers"

type GetTopPlayers struct {
	NumPlayers int64
}

func (gtp *GetTopPlayers) QueryID() string {
	return GetTopPlayersQueryID
}

type GetTopPlayersHandler struct {
	ranking leaderboard.Ranking
}

func NewGetTopPlayersHandler(ranking leaderboard.Ranking) *GetTopPlayersHandler {
	return &GetTopPlayersHandler{
		ranking: ranking,
	}
}

func (gtph *GetTopPlayersHandler) Handle(q cqrs.Query) (cqrs.QueryResponse[any], error) {
	gtpq, ok := q.(*GetTopPlayers)
	if !ok {
		return nil, fmt.Errorf("wrong query: %v", q)
	}

	limit := gtpq.NumPlayers

	membersModel, err := gtph.ranking.GetTopPlayers(limit)
	if err != nil {
		return nil, err
	}
	var members []*leaderboard.Member
	for _, mm := range membersModel {
		m := &leaderboard.Member{
			PlayerID: mm.PlayerID,
			Score:    mm.Score,
		}
		members = append(members, m)
	}

	return &GetTopPlayersResponse{
		MemberList: members,
	}, err
}

func (gtph *GetTopPlayersHandler) QueryID() string {
	return GetTopPlayersQueryID
}

type GetTopPlayersResponse struct {
	MemberList []*leaderboard.Member
}

func (gtpr *GetTopPlayersResponse) Data() any {
	return gtpr
}
