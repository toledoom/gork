package cqrs

import (
	"fmt"
	"reflect"
)

type CommandHandler[T any] func(T) error

type CommandRegistry struct {
	commandHandlers map[string]any
}

func NewCommandRegistry() *CommandRegistry {
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
