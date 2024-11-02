package leaderboard

type Member struct {
	PlayerID string
	Score    int64
}

type Ranking interface {
	GetRank(playerID string) (uint64, error)
	GetTopPlayers(limit int64) ([]Member, error)
	UpdateScore(playerID string, score int64) error
}
