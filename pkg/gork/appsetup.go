package gork

type ServicesSetup func(container *Container)
type UseCasesSetup func(useCaseRegistry *UseCaseRegistry, commandRegistry *CommandRegistry, queryRegistry *QueryRegistry)
type CommandHandlersSetup func(s *Scope, commandRegistry *CommandRegistry)
type QueryHandlersSetup func(s *Scope, queryRegistry *QueryRegistry)
