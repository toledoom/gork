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
	c interface{}
}

func (cnre *CommandNotRegisteredError) Error() string {
	return fmt.Sprintf("command handler not registered for command %s", reflect.TypeOf(cnre.c).String())
}

func HandleCommand[T any](cr *CommandRegistry, c T) error {
	tryCommandHandlerh, ok := cr.commandHandlers[reflect.TypeOf(c).String()]
	if !ok {
		return &CommandNotRegisteredError{c: c}
	}
	ch := tryCommandHandlerh.(CommandHandler[T])
	return ch(c)
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

func RegisterQueryHandler[T, R any](qr *QueryRegistry, qh QueryHandler[T, R]) {
	var t T
	qr.queryHandlers[reflect.TypeOf(t).String()] = qh
}

type QueryNotRegisteredError struct {
	q interface{}
}

func (qnre *QueryNotRegisteredError) Error() string {
	return fmt.Sprintf("query handler not registered for query %s", reflect.TypeOf(qnre.q).String())
}

func HandleQuery[T, R any](qr *QueryRegistry, q T) (R, error) {
	tryQueryHandler, ok := qr.queryHandlers[reflect.TypeOf(q).String()]
	if !ok {
		var r R
		return r, &QueryNotRegisteredError{q: q}
	}

	qh := tryQueryHandler.(QueryHandler[T, R])
	return qh(q)
}
