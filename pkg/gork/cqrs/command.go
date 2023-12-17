package cqrs

import "reflect"

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

func HandleCommand[T any](cr *CommandRegistry, c T) error {
	ch := cr.commandHandlers[reflect.TypeOf(c).String()].(CommandHandler[T])
	return ch(c)
}
