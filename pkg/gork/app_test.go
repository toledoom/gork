package gork_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/toledoom/gork/pkg/gork"
	"github.com/toledoom/gork/pkg/gork/cqrs"
)

func TestAppHandlesCommandsAndQueries(t *testing.T) {
	assert := assert.New(t)

	commandHandlerSetup := func(container *gork.Container, commandRegistry *cqrs.CommandRegistry) {
		cqrs.RegisterCommandHandler[*dumbCommand](commandRegistry, dumbCommandHandler)
	}
	queryHandlersSetup := func(container *gork.Container, commandRegistry *cqrs.QueryRegistry) {
		cqrs.RegisterQueryHandler[*dumbQuery, string](commandRegistry, dumbQueryHandler)
	}
	servicesSetup := func(container *gork.Container) {}
	repositoriesSetup := func(container *gork.Container, uow gork.Worker) {}
	storageMapperSetup := func(datamapper *gork.StorageMapper, container *gork.Container) {}
	eventPublisherSetup := func(eventPublisher *gork.EventPublisher, container *gork.Container) {}

	app := gork.NewApp(commandHandlerSetup, queryHandlersSetup)
	app.Start(servicesSetup, repositoriesSetup, storageMapperSetup, eventPublisherSetup)
	app.SetupCommandsAndQueries(&dumbUnitOfWork{})

	dc := &dumbCommand{}
	err := gork.HandleCommand[*dumbCommand](app, dc)
	assert.Nil(err)

	dq := &dumbQuery{}
	resp, err := gork.HandleQuery[*dumbQuery, string](app, dq)
	assert.Nil(err)
	assert.Equal("a value", resp)
}

func TestAppErrorsWhenHandlingUnregisteredCommandsAndQueries(t *testing.T) {
	assert := assert.New(t)

	commandHandlerSetup := func(container *gork.Container, commandRegistry *cqrs.CommandRegistry) {}
	queryHandlersSetup := func(container *gork.Container, commandRegistry *cqrs.QueryRegistry) {}
	servicesSetup := func(container *gork.Container) {}
	repositoriesSetup := func(container *gork.Container, uow gork.Worker) {}
	storageMapperSetup := func(datamapper *gork.StorageMapper, container *gork.Container) {}
	eventPublisherSetup := func(eventPublisher *gork.EventPublisher, container *gork.Container) {}

	app := gork.NewApp(commandHandlerSetup, queryHandlersSetup)
	app.Start(servicesSetup, repositoriesSetup, storageMapperSetup, eventPublisherSetup)
	app.SetupCommandsAndQueries(&dumbUnitOfWork{})

	dc := &dumbCommand{}
	err := gork.HandleCommand[*dumbCommand](app, dc)
	assert.IsType(&cqrs.CommandNotRegisteredError{}, err)

	dq := &dumbQuery{}
	_, err = gork.HandleQuery[*dumbQuery, string](app, dq)
	assert.IsType(&cqrs.QueryNotRegisteredError{}, err)
}

func TestAppUseUow(t *testing.T) {
	assert := assert.New(t)

	commandHandlerSetup := func(container *gork.Container, commandRegistry *cqrs.CommandRegistry) {
		cqrs.RegisterCommandHandler[*dumbCommand](commandRegistry, persistEntityCommandHandler(gork.GetService[*dumbEntityUowRepository](container)))
	}
	queryHandlersSetup := func(container *gork.Container, commandRegistry *cqrs.QueryRegistry) {}
	servicesSetup := func(container *gork.Container) {}
	repositoriesSetup := func(container *gork.Container, uow gork.Worker) {
		gork.AddService[*dumbEntityUowRepository](container, func(*gork.Container) *dumbEntityUowRepository {
			return &dumbEntityUowRepository{
				uow: uow,
			}
		})
	}
	storageMapperSetup := func(storageMapper *gork.StorageMapper, container *gork.Container) {}
	eventPublisherSetup := func(eventPublisher *gork.EventPublisher, container *gork.Container) {}

	dumbUow := &dumbUnitOfWork{}
	app := gork.NewApp(commandHandlerSetup, queryHandlersSetup)
	app.Start(servicesSetup, repositoriesSetup, storageMapperSetup, eventPublisherSetup)
	app.SetupCommandsAndQueries(dumbUow)

	dc := &dumbCommand{}
	err := gork.HandleCommand[*dumbCommand](app, dc)
	assert.NoError(err)
	assert.Len(dumbUow.entities, 1)
}

type dumbUnitOfWork struct {
	gork.Worker

	entities []gork.Entity
}

func (uow *dumbUnitOfWork) RegisterNew(newEntity gork.Entity) error {
	uow.entities = append(uow.entities, newEntity)
	return nil
}
