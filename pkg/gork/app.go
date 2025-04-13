package gork

type App struct {
	container              *Container
	useCaseBuilderRegistry *UseCaseBuilderRegistry

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

	app.useCaseBuilderRegistry = newUseCaseBuilderRegistry()
	app.useCasesSetup(app.useCaseBuilderRegistry)
}
