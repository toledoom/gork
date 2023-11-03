package command

import (
	"fmt"
	"time"

	"github.com/toledoom/gork/internal/domain/battle"
	"github.com/toledoom/gork/internal/domain/player"
	"github.com/toledoom/gork/pkg/cqrs"
	"github.com/toledoom/gork/pkg/event"
)

const FinishBattleCmdID = "FinishBattle"

type FinishBattle struct {
	BattleID, WinnerID string
}

func (fb *FinishBattle) CmdID() string {
	return FinishBattleCmdID
}

type FinishBattleHandler struct {
	br battle.Repository
	pr player.Repository
	c  battle.ScoreCalculator
	ep event.Publisher
}

func NewFinishBattleHandler(br battle.Repository,
	pr player.Repository,
	c battle.ScoreCalculator,
	ep event.Publisher) *FinishBattleHandler {
	return &FinishBattleHandler{
		br: br,
		pr: pr,
		c:  c,
		ep: ep,
	}
}

func (fbh *FinishBattleHandler) CmdID() string {
	return FinishBattleCmdID
}

func (fbh FinishBattleHandler) Handle(c cqrs.Command) error {
	fbc, ok := c.(*FinishBattle)
	if !ok {
		return fmt.Errorf("wrong command: %v", c)
	}
	battleID := fbc.BattleID
	winnerID := fbc.WinnerID
	finishedAt := time.Now().UTC()
	b, err := fbh.br.GetByID(battleID)
	if err != nil {
		return err
	}
	b.Finish(battleID, winnerID, finishedAt, fbh.c)

	player1ID := b.Player1ID
	player2ID := b.Player2ID
	player1, err := fbh.pr.GetByID(player1ID)
	if err != nil {
		return err
	}
	player2, err := fbh.pr.GetByID(player2ID)
	if err != nil {
		return err
	}

	err = fbh.pr.Update(player1)
	if err != nil {
		return err
	}
	err = fbh.pr.Update(player2)
	if err != nil {
		return err
	}

	// TODO: Delegate to the framework
	err = fbh.ep.Notify(player.NewScoreUpdatedEvent(player1ID, b.OriginalPlayer1Score, player1.Score))
	if err != nil {
		return err
	}
	err = fbh.ep.Notify(player.NewScoreUpdatedEvent(player2ID, b.OriginalPlayer2Score, player2.Score))
	if err != nil {
		return err
	}
	///////////////////////////////////

	return nil
}
