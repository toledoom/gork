package gork_test

import (
	"reflect"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/toledoom/gork/pkg/gork"
)

func TestConcurrentUseCase(t *testing.T) {
	assert := assert.New(t)

	commandHandlerSetup := func(s *gork.Scope, commandRegistry *gork.CommandRegistry) {
		gork.RegisterCommandHandler(commandRegistry, persistEntityCommandHandler(
			gork.GetService[*dumbEntityUowRepository](s), gork.GetService[*dumbNotThreadSafeService](s)))
	}
	queryHandlersSetup := func(s *gork.Scope, queryRegistry *gork.QueryRegistry) {
		gork.RegisterQueryHandler(queryRegistry, getEntityQueryHandler(gork.GetService[*dumbEntityUowRepository](s)))
	}
	servicesSetup := func(container *gork.Container) {
		gork.RegisterService(container, func(s *gork.Scope) gork.Worker {
			return &dumbUnitOfWork{}
		}, gork.USECASE)
		gork.RegisterService(container, func(s *gork.Scope) *dumbEntityUowRepository {
			return &dumbEntityUowRepository{
				uow: gork.GetService[gork.Worker](s),
			}
		}, gork.USECASE)
		gork.RegisterService(container, func(s *gork.Scope) *gork.EventPublisher { return gork.NewPublisher() }, gork.USECASE)
		gork.RegisterService(container, func(s *gork.Scope) *dumbNotThreadSafeService { return &dumbNotThreadSafeService{} }, gork.SINGLETON)
	}
	useCaseSetup := func(ucr *gork.UseCaseBuilderRegistry) {
		gork.RegisterUseCaseBuilder(ucr, concurrentUseCase)
	}

	app := gork.NewApp(useCaseSetup, commandHandlerSetup, queryHandlersSetup)
	app.Start(servicesSetup)

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		output, err := gork.ExecuteUseCase[concurrentInput, concurrentOutput](app, concurrentInput{
			entityID: "1",
			field1:   "random value 1",
		})
		assert.NoError(err)
		assert.Equal("1", output.de.ID)
		assert.Equal("random value 1", output.de.field1)
	}()
	go func() {
		defer wg.Done()
		output, err := gork.ExecuteUseCase[concurrentInput, concurrentOutput](app, concurrentInput{
			entityID: "2",
			field1:   "random value 2",
		})
		assert.NoError(err)
		assert.Equal("2", output.de.ID)
		assert.Equal("random value 2", output.de.field1)

	}()
	wg.Wait()
}

type dumbNotThreadSafeService struct{}

type persistEntityCommand struct {
	entityID, field1 string
}

func persistEntityCommandHandler(repository *dumbEntityUowRepository, dumbService *dumbNotThreadSafeService) func(c *persistEntityCommand) error {
	return func(c *persistEntityCommand) error {
		de := newDumbEntity(c.entityID, c.field1)
		return repository.Add(de)
	}
}

type getEntityQuery struct {
	entityID string
}

type getEntityQueryResponse struct {
	de *dumbEntity
}

func getEntityQueryHandler(repository *dumbEntityUowRepository) func(q *getEntityQuery) (*getEntityQueryResponse, error) {
	return func(q *getEntityQuery) (*getEntityQueryResponse, error) {
		entity, _ := repository.FindByID(q.entityID)
		return &getEntityQueryResponse{de: entity}, nil

	}
}

type dumbEntityUowRepository struct {
	uow gork.Worker
}

func (der *dumbEntityUowRepository) FindByID(id string) (*dumbEntity, error) {
	entity, err := der.uow.FetchOne(reflect.TypeOf(&dumbEntity{}), id)
	if err != nil {
		return nil, err
	}
	d := entity.(*dumbEntity)

	return d, nil
}

func (der *dumbEntityUowRepository) Add(d *dumbEntity) error {
	return der.uow.RegisterNew(d)
}

type dumbEvent struct {
	gork.Event

	dumbID string
}

func (de *dumbEvent) Name() string {
	return "DumbEvent"
}

func newDumbEntity(id, field1 string) *dumbEntity {
	dEnt := &dumbEntity{
		ag:     &gork.Aggregate{},
		ID:     id,
		field1: field1,
	}
	dEnt.AddEvent(&dumbEvent{dumbID: id})
	return dEnt
}

func (de *dumbEntity) AddEvent(e gork.Event) {
	de.ag.Events = append(de.ag.Events, e)
}

func (de *dumbEntity) GetEvents() []gork.Event {
	return de.ag.Events
}

func (de *dumbEntity) Field1() string {
	return de.field1
}

type dumbEntity struct {
	ag *gork.Aggregate

	ID, field1 string
}

type concurrentInput struct {
	entityID, field1 string
}

type concurrentOutput struct {
	de *dumbEntity
}

func concurrentUseCase(cr *gork.CommandRegistry, qr *gork.QueryRegistry) gork.UseCase[concurrentInput, concurrentOutput] {
	return func(ci concurrentInput) (concurrentOutput, error) {
		command := &persistEntityCommand{
			entityID: ci.entityID,
			field1:   ci.field1,
		}
		gork.HandleCommand(cr, command)

		query := &getEntityQuery{
			entityID: ci.entityID,
		}
		resp, _ := gork.HandleQuery[*getEntityQuery, *getEntityQueryResponse](qr, query)

		output := concurrentOutput{
			de: resp.de,
		}

		return output, nil
	}
}
