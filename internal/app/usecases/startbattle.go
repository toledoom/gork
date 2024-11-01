package usecases

import (
	"github.com/toledoom/gork/internal/app/command"
	"github.com/toledoom/gork/pkg/gork/cqrs"
)

type StartBattleInput struct {
	BattleID, Player1ID, Player2ID string
}

type StartBattleOutput struct {
}

func StartBattle(cr *cqrs.CommandRegistry) func(sbi StartBattleInput) (StartBattleOutput, error) {
	return func(sbi StartBattleInput) (StartBattleOutput, error) {
		startBattleCommand := command.StartBattle{
			BattleID:  sbi.BattleID,
			Player1ID: sbi.Player1ID,
			Player2ID: sbi.Player2ID,
		}
		err := cqrs.HandleCommand(cr, &startBattleCommand)
		if err != nil {
			return StartBattleOutput{}, nil
		}

		return StartBattleOutput{}, nil
	}
}
