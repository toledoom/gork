package query

import (
	"github.com/toledoom/gork/internal/domain/leaderboard"
)

type GetTopPlayers struct {
	NumPlayers int64
}

type GetTopPlayersResponse struct {
	MemberList []*leaderboard.Member
}

func GetTopPlayersHandler(ranking leaderboard.Ranking) func(*GetTopPlayers) (*GetTopPlayersResponse, error) {
	return func(q *GetTopPlayers) (*GetTopPlayersResponse, error) {
		limit := q.NumPlayers

		membersModel, err := ranking.GetTopPlayers(limit)
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
}
