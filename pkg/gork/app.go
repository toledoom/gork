package gork

import (
	"net/http"

	"github.com/toledoom/gork/pkg/di"
	"github.com/toledoom/gork/pkg/event"
	"github.com/toledoom/gork/pkg/gork/cqrs"
	"github.com/toledoom/gork/pkg/persistence"
	grpcgork "github.com/toledoom/gork/pkg/ports/grpc"
	httpgork "github.com/toledoom/gork/pkg/ports/http"

	"google.golang.org/grpc"
)

type App struct {
	container         *di.Container
	storageMapper     *persistence.StorageMapper
	setupRepositories persistence.RepositoriesSetup
	eventPublisher    *event.Publisher

	commandRegistry *cqrs.CommandRegistry
	queryRegistry   *cqrs.QueryRegistry

	commandHandlersSetup CommandHandlersSetup
	queryHandlersSetup   QueryHandlersSetup
}

func NewApp(commandHandlersSetup CommandHandlersSetup, queryHandlersSetup QueryHandlersSetup) *App {
	container := di.NewContainer()
	storageMapper := persistence.NewStorageMapper()
	di.AddService[*event.Publisher](container, func(*di.Container) *event.Publisher { return event.NewPublisher() })
	eventPublisher := di.GetService[*event.Publisher](container)

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
	repositoriesSetup persistence.RepositoriesSetup,
	dataMapperSetup DataMapperSetup,
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
	interceptor := grpcgork.WithCommitAndNotifyInterceptor(app.container, app.setupRepositories, app.storageMapper)
	options = append(options, interceptor)
	s := grpc.NewServer(options...)

	return s
}

func (app *App) httpMiddleware(h http.Handler) http.Handler {
	middleware := httpgork.WithCommitAndNotifyMiddleware(app.container, app.setupRepositories, app.storageMapper)
	return middleware(h)
}

func (app *App) HttpListenAndServe(port string, h http.Handler) error {
	gorkHandler := app.httpMiddleware(h)
	return http.ListenAndServe(port, gorkHandler)
}
