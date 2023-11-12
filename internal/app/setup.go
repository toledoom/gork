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
	di.AddService[*redis.Client](container, func(*di.Container) *redis.Client {
		return createRedisLocalClient(os.Getenv("REDIS_ADDR"))
	})
	di.AddService[leaderboarddomain.Ranking](container, func(*di.Container) leaderboarddomain.Ranking {
		redisClient := di.GetService[*redis.Client](container)
		return leaderboard.NewRedisRanking(redisClient, "my-ranking")
	})
	di.AddService[*dynamodb.Client](container, func(*di.Container) *dynamodb.Client {
		return createDynamoDBLocalClient(os.Getenv("DYNAMO_ADDR"))
	})
	di.AddService[playerdomain.Repository](container, func(*di.Container) playerdomain.Repository {
		return player.NewUowRepository(di.GetService[persistence.Worker](container))
	})
	di.AddService[battledomain.Repository](container, func(*di.Container) battledomain.Repository {
		return battle.NewUowRepository(di.GetService[persistence.Worker](container))
	})
	di.AddService[battledomain.ScoreCalculator](container, func(*di.Container) battledomain.ScoreCalculator {
		// A better idea would be to retrieve these next values from a config repository, since they may vary
		// depending on several factors (e.g. players levels). In that case, the solution would be to create a
		// a config repository and inject it as a dependency into the score calculator
		k := int64(20)
		s := int64(400)
		return battledomain.NewEloScoreCalculator(k, s)
	})
}

func SetupCommandHandlers(container *di.Container, cr *cqrs.CommandRegistry) {
	cqrs.RegisterCommandHandler[*command.CreatePlayer](
		cr, command.CreatePlayerHandler(di.GetService[playerdomain.Repository](container)),
	)
	cqrs.RegisterCommandHandler[*command.StartBattle](
		cr, command.StartBattleHandler(
			di.GetService[battledomain.Repository](container),
			di.GetService[playerdomain.Repository](container),
		),
	)
	cqrs.RegisterCommandHandler[*command.FinishBattle](
		cr, command.FinishBattleHandler(
			di.GetService[battledomain.Repository](container),
			di.GetService[playerdomain.Repository](container),
			di.GetService[battledomain.ScoreCalculator](container),
		),
	)
}

func SetupQueryHandlers(container *di.Container, qr *cqrs.QueryRegistry) {
	cqrs.RegisterQueryHandler[*query.GetRank, *query.GetRankResponse](
		qr, query.GetRankHandler(di.GetService[leaderboarddomain.Ranking](container)),
	)
	cqrs.RegisterQueryHandler[*query.GetTopPlayers, *query.GetTopPlayersResponse](
		qr, query.GetTopPlayersHandler(di.GetService[leaderboarddomain.Ranking](container)),
	)
	cqrs.RegisterQueryHandler[*query.GetPlayerByID, *query.GetPlayerByIDResponse](
		qr, query.GetPlayerByIDHandler(di.GetService[playerdomain.Repository](container)),
	)
}

func SetupDataMapper(dataMapper *persistence.StorageMapper, container *di.Container) {
	dataMapper.AddPersistenceFn(reflect.TypeOf(battledomain.Battle{}), persistence.EntityNew, func(e entity.Entity) error {
		b := e.(*battledomain.Battle)
		bdr := battle.NewDynamoStorage(di.GetService[*dynamodb.Client](container))
		return bdr.Add(b)
	})

	dataMapper.AddPersistenceFn(reflect.TypeOf(battledomain.Battle{}), persistence.EntityDirty, func(e entity.Entity) error {
		b := e.(*battledomain.Battle)
		bdr := battle.NewDynamoStorage(di.GetService[*dynamodb.Client](container))
		return bdr.Update(b)
	})
}

func SetupEventPublisher(eventPublisher *event.Publisher, container *di.Container) {
	r := di.GetService[leaderboarddomain.Ranking](container)
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
