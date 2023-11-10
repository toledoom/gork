package command

import (
	"time"

	"github.com/toledoom/gork/internal/domain/battle"
	"github.com/toledoom/gork/internal/domain/player"
)

type StartBattle struct {
	BattleID             string
	Player1ID, Player2ID string
}

func StartBattleHandler(br battle.Repository, pr player.Repository) func(c *StartBattle) error {
	return func(c *StartBattle) error {
		battleID := c.BattleID
		player1ID := c.Player1ID
		player2ID := c.Player2ID

		player1, err := pr.GetByID(player1ID)
		if err != nil {
			return err
		}
		player2, err := pr.GetByID(player2ID)
		if err != nil {
			return err
		}
		originalPlayer1Score := player1.Score
		originalPlayer2Score := player2.Score

		b := battle.New(battleID, player1ID, player2ID, originalPlayer1Score, originalPlayer2Score, time.Now().UTC())
		err = br.Add(b)
		if err != nil {
			return err
		}

		return nil
	}

}
