package application

import (
	"github.com/toledoom/gork/pkg/cqrs"
	"github.com/toledoom/gork/pkg/di"
	"github.com/toledoom/gork/pkg/event"
	"github.com/toledoom/gork/pkg/persistence"
)

type ServicesSetup func(container *di.Container)
type CommandHandlersSetup func(container *di.Container) []cqrs.CommandHandler
type QueryHandlersSetup func(container *di.Container) []cqrs.QueryHandler
type DataMapperSetup func(datamapper *persistence.DataMapper, container *di.Container)
type EventPublisherSetup func(eventPublisher *event.Publisher, container *di.Container)
