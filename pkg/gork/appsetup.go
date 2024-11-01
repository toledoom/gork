package gork

type ServicesSetup func(container *Container)
type UseCasesSetup func(useCaseRegistry *UseCaseRegistry, commandRegistry *CommandRegistry, queryRegistry *QueryRegistry)
type CommandHandlersSetup func(container *Container, commandRegistry *CommandRegistry)
type QueryHandlersSetup func(container *Container, queryRegistry *QueryRegistry)
