package query

import (
	"github.com/toledoom/gork/internal/domain/player"
)

type GetPlayerByID struct {
	PlayerID string
}

type GetPlayerByIDResponse struct {
	Player *player.Player
}

func GetPlayerByIDHandler(pr player.Repository) func(q *GetPlayerByID) (*GetPlayerByIDResponse, error) {
	return func(q *GetPlayerByID) (*GetPlayerByIDResponse, error) {
		id := q.PlayerID
		p, err := pr.GetByID(id)
		if err != nil {
			return nil, err
		}

		return &GetPlayerByIDResponse{
			Player: p,
		}, nil
	}
}
