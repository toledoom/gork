package battle

import (
	"time"

	"github.com/toledoom/gork/pkg/gork"
)

type StartedEvent struct {
	gork.Event

	BattleID  string
	StartedAt time.Time
}

func NewStartedEvent(battleID string, startedAt time.Time) *StartedEvent {
	return &StartedEvent{
		BattleID:  battleID,
		StartedAt: startedAt,
	}
}

func (bse *StartedEvent) Name() string {
	return "BattleStartedEvent"
}

type FinishedEvent struct {
	gork.Event

	BattleID   string
	FinishedAt time.Time
}

func NewFinishedEvent(battleID string, finishedAt time.Time) *FinishedEvent {
	return &FinishedEvent{
		BattleID:   battleID,
		FinishedAt: finishedAt,
	}
}

func (bse *FinishedEvent) Name() string {
	return "BattleFinishedEvent"
}
