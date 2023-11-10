package command

import (
	"github.com/toledoom/gork/internal/domain/player"
)

type CreatePlayer struct {
	PlayerID, Name string
}

func CreatePlayerHandler(pr player.Repository) func(c *CreatePlayer) error {
	return func(c *CreatePlayer) error {
		id := c.PlayerID
		name := c.Name

		p := player.New(id, name)
		err := pr.Add(p)

		if err != nil {
			return err
		}

		return nil
	}
}
