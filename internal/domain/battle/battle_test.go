package battle_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/toledoom/gork/internal/domain/battle"
)

func TestELOScoreCalculator(t *testing.T) {
	assert := assert.New(t)

	c := &battle.EloScoreCalculator{
		K: 20,
		S: 400,
	}
	res := c.Calculate(10, 10, battle.Win)
	assert.Equal(int64(20), res)

	res = c.Calculate(10, 10, battle.Lose)
	assert.Equal(int64(0), res)

	res = c.Calculate(0, 0, battle.Lose)
	assert.Equal(int64(-10), res)
}
