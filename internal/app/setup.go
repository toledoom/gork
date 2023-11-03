package app

import (
	"context"
	"os"
	"reflect"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/go-redis/redis/v8"
	"github.com/toledoom/gork/internal/app/command"
	"github.com/toledoom/gork/internal/app/query"
	battledomain "github.com/toledoom/gork/internal/domain/battle"
	leaderboarddomain "github.com/toledoom/gork/internal/domain/leaderboard"
	playerdomain "github.com/toledoom/gork/internal/domain/player"
	"github.com/toledoom/gork/internal/storage/battle"
	"github.com/toledoom/gork/internal/storage/leaderboard"
	"github.com/toledoom/gork/internal/storage/player"
	"github.com/toledoom/gork/pkg/cqrs"
	"github.com/toledoom/gork/pkg/di"
	"github.com/toledoom/gork/pkg/entity"
	"github.com/toledoom/gork/pkg/event"
	"github.com/toledoom/gork/pkg/persistence"
)

func SetupServices(container *di.Container) {
	container.Add("redis-client", func() any { return createRedisLocalClient(os.Getenv("REDIS_ADDR")) })
	container.Add("ranking", func() any {
		return leaderboard.NewRedisRanking(container.Get("redis-client")().(*redis.Client), "my-ranking")
	})
	container.Add("dynamo-client", func() any { return createDynamoDBLocalClient(os.Getenv("DYNAMO_ADDR")) })
	container.Add("player-repository", func() any {
		return player.NewUowRepository(container.Get("uow")().(persistence.Worker))
	})
	container.Add("battle-repository", func() any {
		return battle.NewUowRepository(container.Get("uow")().(persistence.Worker))
	})
	container.Add("score-calculator", func() any {
		return battledomain.NewEloScoreCalculator(20, 400)
	})
}

func SetupCommandHandlers(container *di.Container) []cqrs.CommandHandler {
	commandHandlerList := []cqrs.CommandHandler{
		command.NewCreatePlayerHandler(container.Get("player-repository")().(playerdomain.Repository)),
		command.NewStartBattleHandler(container.Get("battle-repository")().(battledomain.Repository)),
		command.NewFinishBattleHandler(
			container.Get("battle-repository")().(battledomain.Repository),
			container.Get("player-repository")().(playerdomain.Repository),
			container.Get("score-calculator")().(battledomain.ScoreCalculator),
			*event.NewPublisher(),
		),
	}

	return commandHandlerList
}

func SetupQueryHandlers(container *di.Container) []cqrs.QueryHandler {
	queryHandlerList := []cqrs.QueryHandler{
		query.NewGetRankHandler(container.Get("ranking")().(leaderboarddomain.Ranking)),
		query.NewGetTopPlayersHandler(container.Get("ranking")().(leaderboarddomain.Ranking)),
	}

	return queryHandlerList
}

func SetupDataMapper(dataMapper *persistence.DataMapper, container *di.Container) {
	dataMapper.AddPersistenceFn(reflect.TypeOf(battledomain.Battle{}), persistence.EntityNew, func(e entity.Entity) error {
		b := e.(*battledomain.Battle)
		bdr := battle.NewDynamoStorage(container.Get("dynamo-client")().(*dynamodb.Client))
		return bdr.Add(b)
	})

	dataMapper.AddPersistenceFn(reflect.TypeOf(battledomain.Battle{}), persistence.EntityDirty, func(e entity.Entity) error {
		b := e.(*battledomain.Battle)
		bdr := battle.NewDynamoStorage(container.Get("dynamo-client")().(*dynamodb.Client))
		return bdr.Update(b)
	})
}

func SetupEventPublisher(eventPublisher *event.Publisher, container *di.Container) {
	r := container.Get("ranking")().(leaderboarddomain.Ranking)
	eventPublisher.Subscribe(leaderboarddomain.NewPlayerScoreUpdatedEventHandler(r), &playerdomain.ScoreUpdatedEvent{})
}

func createRedisLocalClient(redisAddr string) *redis.Client {
	if redisAddr == "" {
		redisAddr = "127.0.0.1:6379"
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return rdb
}

func createDynamoDBLocalClient(dynamoAddr string) *dynamodb.Client {
	if dynamoAddr == "" {
		dynamoAddr = "http://127.0.0.1:8000"
	}
	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion("us-east-1"),
		config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
			return aws.Endpoint{URL: dynamoAddr}, nil
		})),
		config.WithCredentialsProvider(credentials.StaticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID: "dummy", SecretAccessKey: "dummy", SessionToken: "dummy",
				Source: "Hard-coded credentials; values are irrelevant for local DynamoDB",
			},
		}),
	)

	if err != nil {
		panic(err)
	}

	return dynamodb.NewFromConfig(cfg)
}
