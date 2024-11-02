package gork_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/toledoom/gork/pkg/gork"
)

func TestAppExecutesUseCase(t *testing.T) {
	assert := assert.New(t)

	commandHandlerSetup := func(s *gork.Scope, commandRegistry *gork.CommandRegistry) {
		gork.RegisterCommandHandler(commandRegistry, dumbCommandHandler)
	}
	queryHandlersSetup := func(s *gork.Scope, commandRegistry *gork.QueryRegistry) {
		gork.RegisterQueryHandler(commandRegistry, dumbQueryHandler)
	}
	useCaseSetup := func(ucr *gork.UseCaseRegistry, cr *gork.CommandRegistry, qr *gork.QueryRegistry) {
		gork.RegisterUseCase(ucr, dumbUseCase(cr, qr))
	}
	servicesSetup := func(container *gork.Container) {
		gork.RegisterService(container, func(s *gork.Scope) gork.Worker { return &dumbUnitOfWork{} }, gork.USECASE)
		gork.RegisterService(container, func(s *gork.Scope) *gork.EventPublisher { return gork.NewPublisher() }, gork.USECASE)
	}

	app := gork.NewApp(useCaseSetup, commandHandlerSetup, queryHandlersSetup)
	app.Start(servicesSetup)

	dumbOutput, err := gork.ExecuteUseCase[dumbUseCaseInput, dumbUseCaseOutput](app, dumbUseCaseInput{})

	assert.Nil(err)
	assert.Equal("a value", dumbOutput.response)
}

func TestAppErrorsWhenHandlingUnregisteredCommands(t *testing.T) {
	assert := assert.New(t)

	commandHandlerSetup := func(s *gork.Scope, commandRegistry *gork.CommandRegistry) {}
	queryHandlersSetup := func(s *gork.Scope, commandRegistry *gork.QueryRegistry) {}
	servicesSetup := func(container *gork.Container) {
		gork.RegisterService(container, func(s *gork.Scope) gork.Worker { return &dumbUnitOfWork{} }, gork.USECASE)
		gork.RegisterService(container, func(s *gork.Scope) *gork.EventPublisher { return gork.NewPublisher() }, gork.USECASE)
	}
	useCaseSetup := func(ucr *gork.UseCaseRegistry, cr *gork.CommandRegistry, qr *gork.QueryRegistry) {
		gork.RegisterUseCase(ucr, dumbUseCase(cr, qr))
	}

	app := gork.NewApp(useCaseSetup, commandHandlerSetup, queryHandlersSetup)
	app.Start(servicesSetup)

	_, err := gork.ExecuteUseCase[dumbUseCaseInput, dumbUseCaseOutput](app, dumbUseCaseInput{})
	assert.IsType(&gork.CommandNotRegisteredError{}, err)
}

func TestAppErrorsWhenHandlingUnregisteredQueries(t *testing.T) {
	assert := assert.New(t)

	commandHandlerSetup := func(s *gork.Scope, commandRegistry *gork.CommandRegistry) {
		gork.RegisterCommandHandler(commandRegistry, persistEntityCommandHandler(gork.GetService[*dumbEntityUowRepository](s)))
	}
	queryHandlersSetup := func(s *gork.Scope, commandRegistry *gork.QueryRegistry) {}
	servicesSetup := func(container *gork.Container) {
		gork.RegisterService(container, func(s *gork.Scope) gork.Worker { return &dumbUnitOfWork{} }, gork.USECASE)
		gork.RegisterService(container, func(s *gork.Scope) *dumbEntityUowRepository {
			return &dumbEntityUowRepository{
				uow: gork.GetService[gork.Worker](s),
			}
		}, gork.USECASE)
		gork.RegisterService(container, func(s *gork.Scope) *gork.EventPublisher { return gork.NewPublisher() }, gork.USECASE)
	}
	useCaseSetup := func(ucr *gork.UseCaseRegistry, cr *gork.CommandRegistry, qr *gork.QueryRegistry) {
		gork.RegisterUseCase(ucr, dumbUseCase(cr, qr))
	}

	app := gork.NewApp(useCaseSetup, commandHandlerSetup, queryHandlersSetup)
	app.Start(servicesSetup)

	_, err := gork.ExecuteUseCase[dumbUseCaseInput, dumbUseCaseOutput](app, dumbUseCaseInput{})
	assert.Error(err)
	assert.IsType(&gork.QueryNotRegisteredError{}, err)
}

type dumbUnitOfWork struct {
	gork.Worker

	entities []gork.Entity
}

func (uow *dumbUnitOfWork) RegisterNew(newEntity gork.Entity) error {
	uow.entities = append(uow.entities, newEntity)
	return nil
}

func (uow *dumbUnitOfWork) Commit() error {
	return nil
}

func (uow *dumbUnitOfWork) DomainEvents() []gork.Event {
	return nil
}
