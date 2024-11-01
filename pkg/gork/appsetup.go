package gork

import "github.com/toledoom/gork/pkg/gork/cqrs"

type ServicesSetup func(container *Container)
type UseCasesSetup func(useCaseRegistry *UseCaseRegistry, commandRegistry *cqrs.CommandRegistry, queryRegistry *cqrs.QueryRegistry)
type CommandHandlersSetup func(container *Container, commandRegistry *cqrs.CommandRegistry)
type QueryHandlersSetup func(container *Container, queryRegistry *cqrs.QueryRegistry)
