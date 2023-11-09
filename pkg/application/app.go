package application

import (
	"net/http"

	"github.com/toledoom/gork/pkg/cqrs"
	"github.com/toledoom/gork/pkg/di"
	"github.com/toledoom/gork/pkg/event"
	"github.com/toledoom/gork/pkg/persistence"
	grpcgork "github.com/toledoom/gork/pkg/ports/grpc"
	httpgork "github.com/toledoom/gork/pkg/ports/http"

	"google.golang.org/grpc"
)

type App struct {
	container      *di.Container
	dataMapper     *persistence.DataMapper
	eventPublisher *event.Publisher

	commandBus    *cqrs.CommandBus
	queryRegistry *cqrs.QueryRegistry

	commandHandlersSetup CommandHandlersSetup
	queryHandlersSetup   QueryHandlersSetup
}

func New(commandHandlersSetup CommandHandlersSetup, queryHandlersSetup QueryHandlersSetup) *App {
	container := di.NewContainer()
	dataMapper := persistence.NewDataMapper()
	eventPublisher := event.NewPublisher()

	return &App{
		container:      container,
		dataMapper:     dataMapper,
		eventPublisher: eventPublisher,

		commandHandlersSetup: commandHandlersSetup,
		queryHandlersSetup:   queryHandlersSetup,
	}
}

func (app *App) Start(servicesSetup ServicesSetup, dataMapperSetup DataMapperSetup, eventPublisherSetup EventPublisherSetup) {
	servicesSetup(app.container)
	di.AddService[persistence.Worker](app.container, func(*di.Container) persistence.Worker { return persistence.NewUnitOfWork(app.dataMapper) })

	dataMapperSetup(app.dataMapper, app.container)
	eventPublisherSetup(app.eventPublisher, app.container)
	di.AddService[*event.Publisher](app.container, func(*di.Container) *event.Publisher { return event.NewPublisher() })

	app.queryRegistry = cqrs.NewQueryRegistry()
	app.queryHandlersSetup(app.container, app.queryRegistry)
	app.commandBus = cqrs.NewCommandBus(app.commandHandlersSetup(app.container))
}

func (app *App) HandleCommand(c cqrs.Command) error {
	return app.commandBus.Handle(c)
}

func (app *App) HandleQuery(q any) (any, error) {
	return cqrs.HandleQuery[any, any](app.queryRegistry, q)
}

func (app *App) GrpcServer(options ...grpc.ServerOption) *grpc.Server {
	interceptor := grpcgork.WithCommitAndNotifyInterceptor(app.container)
	options = append(options, interceptor)
	s := grpc.NewServer(options...)

	return s
}

func (app *App) HttpMiddleware(h http.Handler) http.Handler {
	middleware := httpgork.WithCommitAndNotifyMiddleware(app.container)
	return middleware(h)
}

func (app *App) HttpListenAndServe(port string, h http.Handler) error {
	gorkHandler := app.HttpMiddleware(h)
	return http.ListenAndServe(port, gorkHandler)
}
