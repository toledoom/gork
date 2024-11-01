package gork

type App struct {
	container *Container

	commandRegistry *CommandRegistry
	queryRegistry   *QueryRegistry
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

	app.queryRegistry = NewQueryRegistry()
	app.commandRegistry = NewCommandRegistry()
	app.useCaseRegistry = NewUseCaseRegistry()

	app.useCasesSetup(app.useCaseRegistry, app.commandRegistry, app.queryRegistry)
}

func (app *App) SetupCommandsAndQueries(unitOfWork Worker) {
	app.queryHandlersSetup(app.container, app.queryRegistry)
	app.commandHandlersSetup(app.container, app.commandRegistry)
}
