package leaderboard

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	domain "github.com/toledoom/gork/internal/domain/leaderboard"
)

type RedisRanking struct {
	c    *redis.Client
	name string
}

func NewRedisRanking(c *redis.Client, name string) *RedisRanking {
	return &RedisRanking{
		c:    c,
		name: name,
	}
}

func (r *RedisRanking) GetRank(playerID string) (uint64, error) {
	zRank := r.c.ZRank(context.TODO(), r.name, playerID)
	if zRank.Err() != redis.Nil {
		return 0, zRank.Err()
	}

	return uint64(zRank.Val()), nil
}

func (r *RedisRanking) GetTopPlayers(limit int64) ([]domain.Member, error) {
	zRevRange, err := r.c.ZRevRangeWithScores(context.TODO(), r.name, 0, limit).Result()
	if err != nil {
		return []domain.Member{}, err
	}

	var members []domain.Member
	for _, item := range zRevRange {
		m := domain.Member{
			PlayerID: fmt.Sprintf("%v", item.Member),
			Score:    int64(item.Score),
		}
		members = append(members, m)
	}

	return members, nil
}

func (r *RedisRanking) UpdateScore(playerID string, score int64) error {
	member := &redis.Z{
		Member: playerID,
		Score:  float64(score),
	}
	zAdd := r.c.ZAdd(context.TODO(), r.name, member)

	if zAdd.Err() != redis.Nil {
		return zAdd.Err()
	}

	return nil
}
