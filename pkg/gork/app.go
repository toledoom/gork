package gork

import (
	"github.com/toledoom/gork/pkg/gork/cqrs"
)

type App struct {
	container *Container

	commandRegistry *cqrs.CommandRegistry
	queryRegistry   *cqrs.QueryRegistry
	useCaseRegistry *UseCaseRegistry

	commandHandlersSetup CommandHandlersSetup
	queryHandlersSetup   QueryHandlersSetup
	useCasesSetup        UseCasesSetup
}

func NewApp(
	useCasesSetup UseCasesSetup,
	commandHandlersSetup CommandHandlersSetup,
	queryHandlersSetup QueryHandlersSetup) *App {
	container := newContainer()

	return &App{
		container:            container,
		commandHandlersSetup: commandHandlersSetup,
		queryHandlersSetup:   queryHandlersSetup,
		useCasesSetup:        useCasesSetup,
	}
}

func (app *App) Start(servicesSetup ServicesSetup) {

	servicesSetup(app.container)

	app.queryRegistry = cqrs.NewQueryRegistry()
	app.commandRegistry = cqrs.NewCommandRegistry()
	app.useCaseRegistry = NewUseCaseRegistry()

	app.useCasesSetup(app.useCaseRegistry, app.commandRegistry, app.queryRegistry)
}

func (app *App) SetupCommandsAndQueries(unitOfWork Worker) {
	app.queryHandlersSetup(app.container, app.queryRegistry)
	app.commandHandlersSetup(app.container, app.commandRegistry)
}

func HandleCommand[T any](app *App, c T) error {
	return cqrs.HandleCommand[T](app.commandRegistry, c)
}

func HandleQuery[Q, R any](app *App, q Q) (R, error) {
	return cqrs.HandleQuery[Q, R](app.queryRegistry, q)
}
