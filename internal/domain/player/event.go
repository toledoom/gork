package player

import "github.com/toledoom/gork/pkg/gork"

type ScoreUpdatedEvent struct {
	gork.Event

	PlayerID string
	OldScore int64
	NewScore int64
}

func NewScoreUpdatedEvent(playerID string, oldScore, newScore int64) *ScoreUpdatedEvent {
	return &ScoreUpdatedEvent{
		PlayerID: playerID,
		OldScore: oldScore,
		NewScore: newScore,
	}
}

func (evt *ScoreUpdatedEvent) Name() string {
	return "PlayerScoreUpdated"
}
