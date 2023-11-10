package command

import (
	"time"

	"github.com/toledoom/gork/internal/domain/battle"
	"github.com/toledoom/gork/internal/domain/player"
)

type FinishBattle struct {
	BattleID, WinnerID string
}

func FinishBattleHandler(br battle.Repository, pr player.Repository, calc battle.ScoreCalculator) func(c *FinishBattle) error {
	return func(c *FinishBattle) error {
		battleID := c.BattleID
		winnerID := c.WinnerID
		finishedAt := time.Now().UTC()
		b, err := br.GetByID(battleID)
		if err != nil {
			return err
		}
		b.Finish(battleID, winnerID, finishedAt, calc)

		player1ID := b.Player1ID
		player2ID := b.Player2ID
		player1, err := pr.GetByID(player1ID)
		if err != nil {
			return err
		}
		player2, err := pr.GetByID(player2ID)
		if err != nil {
			return err
		}

		err = pr.Update(player1)
		if err != nil {
			return err
		}
		err = pr.Update(player2)
		if err != nil {
			return err
		}

		return nil
	}
}
