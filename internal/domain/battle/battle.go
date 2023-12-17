package battle

import (
	"time"

	"github.com/toledoom/gork/pkg/gork"
)

type Result int64

const (
	Lose Result = iota
	Win
)

type Battle struct {
	ag *gork.Aggregate

	ID                                         string
	Player1ID                                  string
	Player2ID                                  string
	OriginalPlayer1Score, OriginalPlayer2Score int64
	FinalPlayer1Score, FinalPlayer2Score       int64
	StartedAt, FinishedAt                      time.Time
}

func (b *Battle) AddEvent(e gork.Event) {
	b.ag.AddEvent(e)
}

func (b *Battle) GetEvents() []gork.Event {
	return b.ag.GetEvents()
}

func New(battleID, player1ID, player2ID string, originalPlayer1Score, originalPlayer2Score int64, startedAt time.Time) *Battle {
	b := &Battle{
		ag: &gork.Aggregate{},

		ID:                   battleID,
		Player1ID:            player1ID,
		Player2ID:            player2ID,
		OriginalPlayer1Score: originalPlayer1Score,
		OriginalPlayer2Score: originalPlayer2Score,
		StartedAt:            startedAt,
	}
	b.AddEvent(NewStartedEvent(battleID, startedAt))
	return b
}

func (b *Battle) Finish(battleID, winnerID string, finishedAt time.Time, calculator ScoreCalculator) {
	b.FinishedAt = finishedAt
	if winnerID == b.Player1ID {
		b.FinalPlayer1Score = calculator.Calculate(b.OriginalPlayer1Score, b.OriginalPlayer2Score, Win)
		b.FinalPlayer2Score = calculator.Calculate(b.OriginalPlayer2Score, b.OriginalPlayer1Score, Lose)
	} else {
		b.FinalPlayer1Score = calculator.Calculate(b.OriginalPlayer1Score, b.OriginalPlayer2Score, Lose)
		b.FinalPlayer2Score = calculator.Calculate(b.OriginalPlayer2Score, b.OriginalPlayer1Score, Win)
	}
	b.AddEvent(NewFinishedEvent(battleID, finishedAt))
}

var _ gork.Entity = (*Battle)(nil)

type Repository interface {
	Add(b *Battle) error
	Update(b *Battle) error
	GetByID(id string) (*Battle, error)
}
