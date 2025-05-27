package gork

import (
	"fmt"
	"reflect"
)

type CommandHandler[T any] func(T) error

type CommandRegistry struct {
	commandHandlers map[string]any
}

func newCommandRegistry() *CommandRegistry {
	return &CommandRegistry{
		commandHandlers: make(map[string]any),
	}
}

func RegisterCommandHandler[T any](cr *CommandRegistry, ch CommandHandler[T]) {
	var t T
	cr.commandHandlers[reflect.TypeOf(t).String()] = ch
}

type CommandNotRegisteredError struct {
	command any
}

func (cnre *CommandNotRegisteredError) Error() string {
	return fmt.Sprintf("command handler not registered for command %s", reflect.TypeOf(cnre.command).String())
}

func HandleCommand[T any](commandRegistry *CommandRegistry, command T) error {
	tryCommandHandlerh, ok := commandRegistry.commandHandlers[reflect.TypeOf(command).String()]
	if !ok {
		return &CommandNotRegisteredError{command: command}
	}
	commandHandler := tryCommandHandlerh.(CommandHandler[T])
	return commandHandler(command)
}

type QueryHandler[T, R any] func(T) (R, error)

type QueryRegistry struct {
	queryHandlers map[string]any
}

func newQueryRegistry() *QueryRegistry {
	return &QueryRegistry{
		queryHandlers: make(map[string]any),
	}
}

func RegisterQueryHandler[T, R any](queryRegistry *QueryRegistry, queryHandler QueryHandler[T, R]) {
	var t T
	queryRegistry.queryHandlers[reflect.TypeOf(t).String()] = queryHandler
}

type QueryNotRegisteredError struct {
	query any
}

func (qnre *QueryNotRegisteredError) Error() string {
	return fmt.Sprintf("query handler not registered for query %s", reflect.TypeOf(qnre.query).String())
}

func HandleQuery[T, R any](queryRegistry *QueryRegistry, query T) (R, error) {
	tryQueryHandler, ok := queryRegistry.queryHandlers[reflect.TypeOf(query).String()]
	if !ok {
		var r R
		return r, &QueryNotRegisteredError{query: query}
	}

	queryHandler := tryQueryHandler.(QueryHandler[T, R])
	return queryHandler(query)
}
