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
		cqrs.RegisterCommandHandler[*DumbCommand](commandRegistry, dumbCommandHandler)
	}
	queryHandlersSetup := func(container *gork.Container, commandRegistry *cqrs.QueryRegistry) {
		cqrs.RegisterQueryHandler[*DumbQuery, string](commandRegistry, dumbQueryHandler)
	}
	servicesSetup := func(container *gork.Container) {}
	repositoriesSetup := func(container *gork.Container, uow gork.Worker) {}
	storageMapperSetup := func(datamapper *gork.StorageMapper, container *gork.Container) {}
	eventPublisherSetup := func(eventPublisher *gork.EventPublisher, container *gork.Container) {}

	app := gork.NewApp(commandHandlerSetup, queryHandlersSetup)
	app.Start(servicesSetup, repositoriesSetup, storageMapperSetup, eventPublisherSetup)

	dc := &DumbCommand{}
	err := gork.HandleCommand[*DumbCommand](app, dc)
	assert.Nil(err)

	dq := &DumbQuery{}
	resp, err := gork.HandleQuery[*DumbQuery, string](app, dq)
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

	dc := &DumbCommand{}
	err := gork.HandleCommand[*DumbCommand](app, dc)
	assert.IsType(&cqrs.CommandNotRegisteredError{}, err)

	dq := &DumbQuery{}
	_, err = gork.HandleQuery[*DumbQuery, string](app, dq)
	assert.IsType(&cqrs.QueryNotRegisteredError{}, err)
}

type DumbCommand struct{}

func dumbCommandHandler(dc *DumbCommand) error { return nil }

type DumbQuery struct{}

func dumbQueryHandler(dc *DumbQuery) (string, error) { return "a value", nil }
