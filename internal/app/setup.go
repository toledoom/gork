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
	"github.com/toledoom/gork/pkg/gork"
	"github.com/toledoom/gork/pkg/gork/cqrs"
)

func SetupServices(container *gork.Container) {
	gork.AddService[*redis.Client](container, func(*gork.Container) *redis.Client {
		return createRedisLocalClient(os.Getenv("REDIS_ADDR"))
	})
	gork.AddService[leaderboarddomain.Ranking](container, func(*gork.Container) leaderboarddomain.Ranking {
		redisClient := gork.GetService[*redis.Client](container)
		return leaderboard.NewRedisRanking(redisClient, "my-ranking")
	})
	gork.AddService[*dynamodb.Client](container, func(*gork.Container) *dynamodb.Client {
		return createDynamoDBLocalClient(os.Getenv("DYNAMO_ADDR"))
	})
	gork.AddService[battledomain.ScoreCalculator](container, func(*gork.Container) battledomain.ScoreCalculator {
		// A better idea would be to retrieve these next values from a config repository, since they may vary
		// depending on several factors (e.g. players levels). In that case, the solution would be to create a
		// a config repository and inject it as a dependency into the score calculator
		k := int64(20)
		s := int64(400)
		return battledomain.NewEloScoreCalculator(k, s)
	})
	gork.AddService[*battle.DynamoStorage](container, func(*gork.Container) *battle.DynamoStorage {
		return battle.NewDynamoStorage(gork.GetService[*dynamodb.Client](container))
	})
	gork.AddService[*player.DynamoStorage](container, func(*gork.Container) *player.DynamoStorage {
		return player.NewDynamoStorage(gork.GetService[*dynamodb.Client](container))
	})
}

func SetupRepositories(container *gork.Container, uow gork.Worker) {
	gork.AddService[playerdomain.Repository](container, func(*gork.Container) playerdomain.Repository {
		return player.NewUowRepository(uow)
	})
	gork.AddService[battledomain.Repository](container, func(*gork.Container) battledomain.Repository {
		return battle.NewUowRepository(uow)
	})
}

func SetupCommandHandlers(container *gork.Container, cr *cqrs.CommandRegistry) {
	cqrs.RegisterCommandHandler[*command.CreatePlayer](
		cr, command.CreatePlayerHandler(gork.GetService[playerdomain.Repository](container)),
	)
	cqrs.RegisterCommandHandler[*command.StartBattle](
		cr, command.StartBattleHandler(
			gork.GetService[battledomain.Repository](container),
			gork.GetService[playerdomain.Repository](container),
		),
	)
	cqrs.RegisterCommandHandler[*command.FinishBattle](
		cr, command.FinishBattleHandler(
			gork.GetService[battledomain.Repository](container),
			gork.GetService[playerdomain.Repository](container),
			gork.GetService[battledomain.ScoreCalculator](container),
		),
	)
}

func SetupQueryHandlers(container *gork.Container, qr *cqrs.QueryRegistry) {
	cqrs.RegisterQueryHandler[*query.GetRank, *query.GetRankResponse](
		qr, query.GetRankHandler(gork.GetService[leaderboarddomain.Ranking](container)),
	)
	cqrs.RegisterQueryHandler[*query.GetTopPlayers, *query.GetTopPlayersResponse](
		qr, query.GetTopPlayersHandler(gork.GetService[leaderboarddomain.Ranking](container)),
	)
	cqrs.RegisterQueryHandler[*query.GetPlayerByID, *query.GetPlayerByIDResponse](
		qr, query.GetPlayerByIDHandler(gork.GetService[playerdomain.Repository](container)),
	)
}

func SetupStorageMapper(storageMapper *gork.StorageMapper, container *gork.Container) {
	storageMapper.AddMutationFn(reflect.TypeOf(battledomain.Battle{}), gork.CreationQuery, gork.GetService[*battle.DynamoStorage](container).Add)
	storageMapper.AddMutationFn(reflect.TypeOf(battledomain.Battle{}), gork.UpdateQuery, gork.GetService[*battle.DynamoStorage](container).Update)
	storageMapper.AddFetchOneFn(reflect.TypeOf(battledomain.Battle{}), gork.GetService[*battle.DynamoStorage](container).GetByID)

	storageMapper.AddMutationFn(reflect.TypeOf(playerdomain.Player{}), gork.CreationQuery, gork.GetService[*player.DynamoStorage](container).Add)
	storageMapper.AddMutationFn(reflect.TypeOf(playerdomain.Player{}), gork.UpdateQuery, gork.GetService[*player.DynamoStorage](container).Update)
	storageMapper.AddFetchOneFn(reflect.TypeOf(playerdomain.Player{}), gork.GetService[*player.DynamoStorage](container).GetByID)
}

func SetupEventPublisher(eventPublisher *gork.EventPublisher, container *gork.Container) {
	r := gork.GetService[leaderboarddomain.Ranking](container)
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
