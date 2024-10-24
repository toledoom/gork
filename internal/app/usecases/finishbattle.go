package usecases

import (
	"github.com/toledoom/gork/internal/app/command"
	"github.com/toledoom/gork/internal/app/query"
	"github.com/toledoom/gork/pkg/gork/cqrs"
)

type FinishBattleUseCase struct {
	Cr *cqrs.CommandRegistry
	Qr *cqrs.QueryRegistry
}

type FinishBattleInput struct {
	BattleID, WinnerID string
}

type FinishBattleOutput struct {
	Player1Score, Player2Score int64
}

func (fbuc *FinishBattleUseCase) Execute(fbi FinishBattleInput) (FinishBattleOutput, error) {
	finishBattleCommand := command.FinishBattle{
		BattleID: fbi.BattleID,
		WinnerID: fbi.WinnerID,
	}

	err := cqrs.HandleCommand(fbuc.Cr, &finishBattleCommand)
	if err != nil {
		return FinishBattleOutput{}, err
	}

	getBattleResultQuery := query.GetBattleResult{
		BattleID: fbi.BattleID,
	}

	queryResult, err := cqrs.HandleQuery[*query.GetBattleResult, *query.GetBattleResultResponse](fbuc.Qr, &getBattleResultQuery)
	if err != nil {
		return FinishBattleOutput{}, err
	}

	return FinishBattleOutput{
		Player1Score: queryResult.Player1Score,
		Player2Score: queryResult.Player2Score,
	}, nil
}
