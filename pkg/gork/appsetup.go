package gork

type ServicesSetup func(container *Container)
type UseCasesSetup func(useCaseRegistry *UseCaseBuilderRegistry)
type CommandHandlersSetup func(scope *Scope, commandRegistry *CommandRegistry)
type QueryHandlersSetup func(scope *Scope, queryRegistry *QueryRegistry)
