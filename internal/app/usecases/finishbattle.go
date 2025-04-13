package usecases

import (
	"github.com/toledoom/gork/internal/app/command"
	"github.com/toledoom/gork/internal/app/query"
	"github.com/toledoom/gork/pkg/gork"
)

type FinishBattleInput struct {
	BattleID, WinnerID string
}

type FinishBattleOutput struct {
	Player1Score, Player2Score int64
}

func FinishBattle(cr *gork.CommandRegistry, qr *gork.QueryRegistry) gork.UseCase[FinishBattleInput, FinishBattleOutput] {
	return func(fbi FinishBattleInput) (FinishBattleOutput, error) {
		finishBattleCommand := command.FinishBattle{
			BattleID: fbi.BattleID,
			WinnerID: fbi.WinnerID,
		}

		err := gork.HandleCommand(cr, &finishBattleCommand)
		if err != nil {
			return FinishBattleOutput{}, err
		}

		getBattleResultQuery := query.GetBattleResult{
			BattleID: fbi.BattleID,
		}

		queryResult, err := gork.HandleQuery[*query.GetBattleResult, *query.GetBattleResultResponse](qr, &getBattleResultQuery)
		if err != nil {
			return FinishBattleOutput{}, err
		}

		return FinishBattleOutput{
			Player1Score: queryResult.Player1Score,
			Player2Score: queryResult.Player2Score,
		}, nil
	}
}
