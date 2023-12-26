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

	dc := &dumbCommand{}
	err := gork.HandleCommand[*dumbCommand](app, dc)
	assert.IsType(&cqrs.CommandNotRegisteredError{}, err)

	dq := &dumbQuery{}
	_, err = gork.HandleQuery[*dumbQuery, string](app, dq)
	assert.IsType(&cqrs.QueryNotRegisteredError{}, err)
}
