package gork

import (
	"net/http"

	"github.com/toledoom/gork/pkg/gork/cqrs"

	"google.golang.org/grpc"
)

type App struct {
	container         *Container
	storageMapper     *StorageMapper
	repositoriesSetup RepositoriesSetup
	eventPublisher    *EventPublisher

	commandRegistry *cqrs.CommandRegistry
	queryRegistry   *cqrs.QueryRegistry

	commandHandlersSetup CommandHandlersSetup
	queryHandlersSetup   QueryHandlersSetup
}

func NewApp(commandHandlersSetup CommandHandlersSetup, queryHandlersSetup QueryHandlersSetup) *App {
	container := newContainer()
	storageMapper := newStorageMapper()
	AddService[*EventPublisher](container, func(*Container) *EventPublisher { return newPublisher() })
	eventPublisher := GetService[*EventPublisher](container)

	return &App{
		container:      container,
		storageMapper:  storageMapper,
		eventPublisher: eventPublisher,

		commandHandlersSetup: commandHandlersSetup,
		queryHandlersSetup:   queryHandlersSetup,
	}
}

func (app *App) Start(
	servicesSetup ServicesSetup,
	repositoriesSetup RepositoriesSetup,
	dataMapperSetup StorageMapperSetup,
	eventPublisherSetup EventPublisherSetup) {

	servicesSetup(app.container)
	dataMapperSetup(app.storageMapper, app.container)
	eventPublisherSetup(app.eventPublisher, app.container)

	app.queryRegistry = cqrs.NewQueryRegistry()
	app.queryHandlersSetup(app.container, app.queryRegistry)
	app.commandRegistry = cqrs.NewCommandRegistry()
	app.commandHandlersSetup(app.container, app.commandRegistry)
}

func HandleCommand[T any](app *App, c T) error {
	return cqrs.HandleCommand[T](app.commandRegistry, c)
}

func HandleQuery[Q, R any](app *App, q Q) (R, error) {
	return cqrs.HandleQuery[Q, R](app.queryRegistry, q)
}

func (app *App) GrpcServer(options ...grpc.ServerOption) *grpc.Server {
	interceptor := withCommitAndNotifyInterceptor(app.container, app.repositoriesSetup, app.storageMapper)
	options = append(options, interceptor)
	s := grpc.NewServer(options...)

	return s
}

func (app *App) HttpListenAndServe(port string, h http.Handler) error {
	middleware := withCommitAndNotifyMiddleware(app.container, app.repositoriesSetup, app.storageMapper)
	gorkHandler := middleware(h)
	return http.ListenAndServe(port, gorkHandler)
}
