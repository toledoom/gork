package query

import (
	"github.com/toledoom/gork/internal/domain/battle"
)

type GetBattleResult struct {
	BattleID string
}

type GetBattleResultResponse struct {
	Player1Score int64
	Player2Score int64
}

func GetBattleResultHandler(br battle.Repository) func(q *GetBattleResult) (*GetBattleResultResponse, error) {
	return func(q *GetBattleResult) (*GetBattleResultResponse, error) {
		id := q.BattleID
		b, err := br.GetByID(id)
		if err != nil {
			return nil, err
		}

		return &GetBattleResultResponse{
			Player1Score: b.FinalPlayer1Score,
			Player2Score: b.FinalPlayer2Score,
		}, nil
	}
}
