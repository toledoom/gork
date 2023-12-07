package application

import (
	"github.com/toledoom/gork/pkg/application/cqrs"
	"github.com/toledoom/gork/pkg/di"
	"github.com/toledoom/gork/pkg/event"
	"github.com/toledoom/gork/pkg/persistence"
)

type ServicesSetup func(container *di.Container)
type CommandHandlersSetup func(container *di.Container, commandRegistry *cqrs.CommandRegistry)
type QueryHandlersSetup func(container *di.Container, queryRegistry *cqrs.QueryRegistry)
type DataMapperSetup func(datamapper *persistence.StorageMapper, container *di.Container)
type EventPublisherSetup func(eventPublisher *event.Publisher, container *di.Container)
