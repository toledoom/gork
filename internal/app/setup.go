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
	"github.com/toledoom/gork/internal/app/usecases"
	battledomain "github.com/toledoom/gork/internal/domain/battle"
	leaderboarddomain "github.com/toledoom/gork/internal/domain/leaderboard"
	playerdomain "github.com/toledoom/gork/internal/domain/player"
	"github.com/toledoom/gork/internal/storage/battle"
	"github.com/toledoom/gork/internal/storage/leaderboard"
	"github.com/toledoom/gork/internal/storage/player"
	"github.com/toledoom/gork/pkg/gork"
)

func SetupServices(container *gork.Container) {
	gork.RegisterService(container, func(s *gork.Scope) gork.Worker {
		return gork.NewUnitOfWork(gork.GetService[*gork.StorageMapper](s))
	}, gork.USECASE)

	gork.RegisterService(container, func(s *gork.Scope) playerdomain.Repository {
		return player.NewUowRepository(gork.GetService[gork.Worker](s))
	}, gork.USECASE)

	gork.RegisterService(container, func(s *gork.Scope) battledomain.Repository {
		return battle.NewUowRepository(gork.GetService[gork.Worker](s))
	}, gork.USECASE)

	gork.RegisterService(container, func(s *gork.Scope) *gork.EventPublisher {
		eventPublisher := gork.NewPublisher()
		r := gork.GetService[leaderboarddomain.Ranking](s)
		eventPublisher.Subscribe(leaderboarddomain.NewPlayerScoreUpdatedEventHandler(r), &playerdomain.ScoreUpdatedEvent{})
		return eventPublisher
	}, gork.USECASE)

	gork.RegisterService(container, func(s *gork.Scope) *gork.StorageMapper {
		storageMapper := gork.NewStorageMapper()
		storageMapper.AddMutationFn(reflect.TypeOf(battledomain.Battle{}), gork.CreationQuery, gork.GetService[*battle.DynamoStorage](s).Add)
		storageMapper.AddMutationFn(reflect.TypeOf(battledomain.Battle{}), gork.UpdateQuery, gork.GetService[*battle.DynamoStorage](s).Update)
		storageMapper.AddFetchOneFn(reflect.TypeOf(battledomain.Battle{}), gork.GetService[*battle.DynamoStorage](s).GetByID)

		storageMapper.AddMutationFn(reflect.TypeOf(playerdomain.Player{}), gork.CreationQuery, gork.GetService[*player.DynamoStorage](s).Add)
		storageMapper.AddMutationFn(reflect.TypeOf(playerdomain.Player{}), gork.UpdateQuery, gork.GetService[*player.DynamoStorage](s).Update)
		storageMapper.AddFetchOneFn(reflect.TypeOf(playerdomain.Player{}), gork.GetService[*player.DynamoStorage](s).GetByID)
		return storageMapper
	}, gork.SINGLETON)

	gork.RegisterService(container, func(*gork.Scope) *redis.Client {
		return createRedisLocalClient(os.Getenv("REDIS_ADDR"))
	}, gork.SINGLETON)
	gork.RegisterService(container, func(s *gork.Scope) leaderboarddomain.Ranking {
		redisClient := gork.GetService[*redis.Client](s)
		return leaderboard.NewRedisRanking(redisClient, "my-ranking")
	}, gork.SINGLETON)
	gork.RegisterService(container, func(*gork.Scope) *dynamodb.Client {
		return createDynamoDBLocalClient(os.Getenv("DYNAMO_ADDR"))
	}, gork.SINGLETON)
	gork.RegisterService(container, func(*gork.Scope) battledomain.ScoreCalculator {
		// A better idea would be to retrieve these next values from a config repository, since they may vary
		// depending on several factors (e.g. players levels). In that case, the solution would be to create a
		// a config repository and inject it as a dependency into the score calculator
		k := int64(20)
		s := int64(400)
		return battledomain.NewEloScoreCalculator(k, s)
	}, gork.SINGLETON)
	gork.RegisterService(container, func(s *gork.Scope) *battle.DynamoStorage {
		return battle.NewDynamoStorage(gork.GetService[*dynamodb.Client](s))
	}, gork.SINGLETON)
	gork.RegisterService(container, func(s *gork.Scope) *player.DynamoStorage {
		return player.NewDynamoStorage(gork.GetService[*dynamodb.Client](s))
	}, gork.SINGLETON)
}

func SetupCommandHandlers(s *gork.Scope, cr *gork.CommandRegistry) {
	gork.RegisterCommandHandler(
		cr, command.CreatePlayerHandler(gork.GetService[playerdomain.Repository](s)),
	)
	gork.RegisterCommandHandler(
		cr, command.StartBattleHandler(
			gork.GetService[battledomain.Repository](s),
			gork.GetService[playerdomain.Repository](s),
		),
	)
	gork.RegisterCommandHandler(
		cr, command.FinishBattleHandler(
			gork.GetService[battledomain.Repository](s),
			gork.GetService[playerdomain.Repository](s),
			gork.GetService[battledomain.ScoreCalculator](s),
		),
	)
}

func SetupQueryHandlers(s *gork.Scope, qr *gork.QueryRegistry) {
	gork.RegisterQueryHandler(qr, query.GetRankHandler(gork.GetService[leaderboarddomain.Ranking](s)))
	gork.RegisterQueryHandler(qr, query.GetTopPlayersHandler(gork.GetService[leaderboarddomain.Ranking](s)))
	gork.RegisterQueryHandler(qr, query.GetPlayerByIDHandler(gork.GetService[playerdomain.Repository](s)))
	gork.RegisterQueryHandler(qr, query.GetBattleResultHandler(gork.GetService[battledomain.Repository](s)))
}

func SetupUseCases(ucbr *gork.UseCaseBuilderRegistry) {
	gork.RegisterUseCaseBuilder(ucbr, usecases.CreatePlayer)
	gork.RegisterUseCaseBuilder(ucbr, usecases.FinishBattle)
	gork.RegisterUseCaseBuilder(ucbr, usecases.GetPlayerByID)
	gork.RegisterUseCaseBuilder(ucbr, usecases.GetRank)
	gork.RegisterUseCaseBuilder(ucbr, usecases.GetTopPlayers)
	gork.RegisterUseCaseBuilder(ucbr, usecases.StartBattle)
}

func SetupEventPublisher(s *gork.Scope, eventPublisher *gork.EventPublisher) {
	r := gork.GetService[leaderboarddomain.Ranking](s)
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
