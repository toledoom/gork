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

	app.queryRegistry = newQueryRegistry()
	app.commandRegistry = newCommandRegistry()
	app.useCaseRegistry = newUseCaseRegistry()
	app.useCasesSetup(app.useCaseRegistry, app.commandRegistry, app.queryRegistry)
}
