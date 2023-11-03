package command

import (
	"fmt"
	"time"

	"github.com/toledoom/gork/internal/domain/battle"
	"github.com/toledoom/gork/internal/domain/player"
	"github.com/toledoom/gork/pkg/cqrs"
)

const StartBattleCmdID = "StartBattle"

type StartBattle struct {
	BattleID             string
	Player1ID, Player2ID string
}

func (sbc *StartBattle) CmdID() string {
	return "StartBattle"
}

type StartBattleHandler struct {
	br battle.Repository
	pr player.Repository
}

func NewStartBattleHandler(br battle.Repository) *StartBattleHandler {
	return &StartBattleHandler{
		br: br,
	}
}

func (sbh *StartBattleHandler) CmdID() string {
	return StartBattleCmdID
}

func (sb StartBattleHandler) Handle(c cqrs.Command) error {
	sbc, ok := c.(*StartBattle)
	if !ok {
		return fmt.Errorf("wrong command: %v", c)
	}

	battleID := sbc.BattleID
	player1ID := sbc.Player1ID
	player2ID := sbc.Player2ID

	player1, err := sb.pr.GetByID(player1ID)
	if err != nil {
		return err
	}
	player2, err := sb.pr.GetByID(player2ID)
	if err != nil {
		return err
	}
	originalPlayer1Score := player1.Score
	originalPlayer2Score := player2.Score

	b := battle.New(battleID, player1ID, player2ID, originalPlayer1Score, originalPlayer2Score, time.Now().UTC())
	err = sb.br.Add(b)
	if err != nil {
		return err
	}

	return nil
}
