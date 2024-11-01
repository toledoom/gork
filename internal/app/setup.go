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
	gork.AddService(container, func(c *gork.Container) gork.Worker {
		return gork.NewUnitOfWork(gork.GetService[*gork.StorageMapper](c))
	}, gork.USECASE_SCOPE)

	gork.AddService(container, func(c *gork.Container) playerdomain.Repository {
		return player.NewUowRepository(gork.GetService[gork.Worker](c))
	}, gork.USECASE_SCOPE)

	gork.AddService(container, func(c *gork.Container) battledomain.Repository {
		return battle.NewUowRepository(gork.GetService[gork.Worker](c))
	}, gork.USECASE_SCOPE)

	gork.AddService(container, func(c *gork.Container) *gork.EventPublisher {
		eventPublisher := gork.NewPublisher()
		r := gork.GetService[leaderboarddomain.Ranking](c)
		eventPublisher.Subscribe(leaderboarddomain.NewPlayerScoreUpdatedEventHandler(r), &playerdomain.ScoreUpdatedEvent{})
		return eventPublisher
	}, gork.USECASE_SCOPE)

	gork.AddService(container, func(c *gork.Container) *gork.StorageMapper {
		storageMapper := gork.NewStorageMapper()
		storageMapper.AddMutationFn(reflect.TypeOf(battledomain.Battle{}), gork.CreationQuery, gork.GetService[*battle.DynamoStorage](c).Add)
		storageMapper.AddMutationFn(reflect.TypeOf(battledomain.Battle{}), gork.UpdateQuery, gork.GetService[*battle.DynamoStorage](c).Update)
		storageMapper.AddFetchOneFn(reflect.TypeOf(battledomain.Battle{}), gork.GetService[*battle.DynamoStorage](c).GetByID)

		storageMapper.AddMutationFn(reflect.TypeOf(playerdomain.Player{}), gork.CreationQuery, gork.GetService[*player.DynamoStorage](c).Add)
		storageMapper.AddMutationFn(reflect.TypeOf(playerdomain.Player{}), gork.UpdateQuery, gork.GetService[*player.DynamoStorage](c).Update)
		storageMapper.AddFetchOneFn(reflect.TypeOf(playerdomain.Player{}), gork.GetService[*player.DynamoStorage](c).GetByID)
		return storageMapper
	}, gork.SERVER_SCOPE)

	gork.AddService(container, func(*gork.Container) *redis.Client {
		return createRedisLocalClient(os.Getenv("REDIS_ADDR"))
	}, gork.SERVER_SCOPE)
	gork.AddService(container, func(c *gork.Container) leaderboarddomain.Ranking {
		redisClient := gork.GetService[*redis.Client](c)
		return leaderboard.NewRedisRanking(redisClient, "my-ranking")
	}, gork.SERVER_SCOPE)
	gork.AddService(container, func(*gork.Container) *dynamodb.Client {
		return createDynamoDBLocalClient(os.Getenv("DYNAMO_ADDR"))
	}, gork.SERVER_SCOPE)
	gork.AddService(container, func(*gork.Container) battledomain.ScoreCalculator {
		// A better idea would be to retrieve these next values from a config repository, since they may vary
		// depending on several factors (e.g. players levels). In that case, the solution would be to create a
		// a config repository and inject it as a dependency into the score calculator
		k := int64(20)
		s := int64(400)
		return battledomain.NewEloScoreCalculator(k, s)
	}, gork.SERVER_SCOPE)
	gork.AddService(container, func(c *gork.Container) *battle.DynamoStorage {
		return battle.NewDynamoStorage(gork.GetService[*dynamodb.Client](c))
	}, gork.SERVER_SCOPE)
	gork.AddService(container, func(c *gork.Container) *player.DynamoStorage {
		return player.NewDynamoStorage(gork.GetService[*dynamodb.Client](c))
	}, gork.SERVER_SCOPE)
}

func SetupCommandHandlers(container *gork.Container, cr *gork.CommandRegistry) {
	gork.RegisterCommandHandler(
		cr, command.CreatePlayerHandler(gork.GetService[playerdomain.Repository](container)),
	)
	gork.RegisterCommandHandler(
		cr, command.StartBattleHandler(
			gork.GetService[battledomain.Repository](container),
			gork.GetService[playerdomain.Repository](container),
		),
	)
	gork.RegisterCommandHandler(
		cr, command.FinishBattleHandler(
			gork.GetService[battledomain.Repository](container),
			gork.GetService[playerdomain.Repository](container),
			gork.GetService[battledomain.ScoreCalculator](container),
		),
	)
}

func SetupQueryHandlers(container *gork.Container, qr *gork.QueryRegistry) {
	gork.RegisterQueryHandler(qr, query.GetRankHandler(gork.GetService[leaderboarddomain.Ranking](container)))
	gork.RegisterQueryHandler(qr, query.GetTopPlayersHandler(gork.GetService[leaderboarddomain.Ranking](container)))
	gork.RegisterQueryHandler(qr, query.GetPlayerByIDHandler(gork.GetService[playerdomain.Repository](container)))
	gork.RegisterQueryHandler(qr, query.GetBattleResultHandler(gork.GetService[battledomain.Repository](container)))
}

func SetupUseCases(ucr *gork.UseCaseRegistry, cr *gork.CommandRegistry, qr *gork.QueryRegistry) {
	gork.RegisterUseCase(ucr, usecases.CreatePlayer(cr, qr))
	gork.RegisterUseCase(ucr, usecases.FinishBattle(cr, qr))
	gork.RegisterUseCase(ucr, usecases.GetPlayerByID(qr))
	gork.RegisterUseCase(ucr, usecases.GetRank(qr))
	gork.RegisterUseCase(ucr, usecases.GetTopPlayers(qr))
	gork.RegisterUseCase(ucr, usecases.StartBattle(cr))
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
