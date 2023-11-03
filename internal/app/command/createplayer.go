package command

import (
	"fmt"

	"github.com/toledoom/gork/internal/domain/player"
	"github.com/toledoom/gork/pkg/cqrs"
)

const CreatePlayerHandlerCmdID = "CreatePlayer"

type CreatePlayerHandler struct {
	pr player.Repository
}

func (cph *CreatePlayerHandler) CmdID() string {
	return CreatePlayerHandlerCmdID
}

func NewCreatePlayerHandler(pr player.Repository) *CreatePlayerHandler {
	return &CreatePlayerHandler{
		pr: pr,
	}
}

type CreatePlayer struct {
	PlayerID, Name string
}

func (cpc *CreatePlayer) CmdID() string {
	return CreatePlayerHandlerCmdID
}

func (cp *CreatePlayerHandler) Handle(c cqrs.Command) error {
	cpc, ok := c.(*CreatePlayer)
	if !ok {
		return fmt.Errorf("wrong command: %v", c)
	}

	id := cpc.PlayerID
	name := cpc.Name

	p := player.New(id, name)
	err := cp.pr.Add(p)

	if err != nil {
		return err
	}

	return nil
}
