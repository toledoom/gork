package battle

import "math"

type ScoreCalculator interface {
	Calculate(score1, score2 int64, battleResult Result) int64
}

type EloScoreCalculator struct {
	K, S int64 // K normally varies depending on several factors: battle result (win, lose), player score... Here I am simplifying
}

func NewEloScoreCalculator(k, s int64) *EloScoreCalculator {
	return &EloScoreCalculator{
		K: k,
		S: s,
	}
}

func (e *EloScoreCalculator) Calculate(score1, score2 int64, battleResult Result) int64 {
	diffRatio := float64((score1 - score2) / e.S)
	x := math.Pow(10, diffRatio)
	expectedScore := 1 / (1 + x)
	br := float64(int64(battleResult))
	delta := float64(e.K) * (br - expectedScore)
	finalScore := float64(score1) + delta

	return int64(finalScore)
}
