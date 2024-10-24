package usecases

import (
	"github.com/toledoom/gork/internal/app/command"
	"github.com/toledoom/gork/pkg/gork/cqrs"
)

type StartBattleUseCase struct {
	cr *cqrs.CommandRegistry
}

type StartBattleInput struct {
	BattleID, Player1ID, Player2ID string
}

type StartBattleOutput struct {
}

func (sbuc *StartBattleUseCase) Execute(sbi StartBattleInput) (StartBattleOutput, error) {
	startBattleCommand := command.StartBattle{
		BattleID:  sbi.BattleID,
		Player1ID: sbi.Player1ID,
		Player2ID: sbi.Player2ID,
	}
	err := cqrs.HandleCommand(sbuc.cr, &startBattleCommand)
	if err != nil {
		return StartBattleOutput{}, nil
	}

	return StartBattleOutput{}, nil
}
