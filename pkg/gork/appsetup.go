package gork

import "github.com/toledoom/gork/pkg/gork/cqrs"

type ServicesSetup func(container *Container)
type CommandHandlersSetup func(container *Container, commandRegistry *cqrs.CommandRegistry)
type QueryHandlersSetup func(container *Container, queryRegistry *cqrs.QueryRegistry)
type StorageMapperSetup func(datamapper *StorageMapper, container *Container)
type EventPublisherSetup func(eventPublisher *EventPublisher, container *Container)
