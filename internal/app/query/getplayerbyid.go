package query

import (
	"fmt"

	"github.com/toledoom/gork/internal/domain/player"
	"github.com/toledoom/gork/pkg/cqrs"
)

const GetPlayerByIDQueryID = "GetPlayerByID"

type GetPlayerByID struct {
	PlayerID string
}

func (gpq *GetPlayerByID) QueryID() string {
	return GetPlayerByIDQueryID
}

type GetPlayerByIDHandler struct {
	pr player.Repository
}

func NewGetPlayerByIDHandler(pr player.Repository) *GetPlayerByIDHandler {
	return &GetPlayerByIDHandler{
		pr: pr,
	}
}

func (gph *GetPlayerByIDHandler) QueryID() string {
	return GetPlayerByIDQueryID
}

func (gph *GetPlayerByIDHandler) Handle(q cqrs.Query) (cqrs.QueryResponse[any], error) {
	gpq, ok := q.(*GetPlayerByID)
	if !ok {
		return nil, fmt.Errorf("wrong query: %v", q)
	}
	id := gpq.PlayerID
	p, err := gph.pr.GetByID(id)
	if err != nil {
		return nil, err
	}

	return &GetPlayerByIDResponse{
		Player: p,
	}, nil
}

type GetPlayerByIDResponse struct {
	Player *player.Player
}

func (gpr *GetPlayerByIDResponse) Data() any {
	return gpr
}
