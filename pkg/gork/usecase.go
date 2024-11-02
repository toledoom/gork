package gork

import (
	"fmt"
	"reflect"
)

type UseCase[I, O any] func(I) (O, error)

type UseCaseRegistry struct {
	useCases map[string]any
}

func newUseCaseRegistry() *UseCaseRegistry {
	return &UseCaseRegistry{
		useCases: make(map[string]any),
	}
}

func RegisterUseCase[I, O any](ucr *UseCaseRegistry, uc UseCase[I, O]) {
	var t I
	ucr.useCases[reflect.TypeOf(t).String()] = uc
}

type UseCaseNotRegisteredError struct {
	uc interface{}
}

func (ucnre *UseCaseNotRegisteredError) Error() string {
	return fmt.Sprintf("use case handler not registered for use case: %s", reflect.TypeOf(ucnre.uc).String())
}

func ExecuteUseCase[I, O any](app *App, input I) (O, error) {
	tryUseCase, ok := app.useCaseRegistry.useCases[reflect.TypeOf(input).String()]
	if !ok {
		var r O
		return r, &UseCaseNotRegisteredError{uc: input}
	}

	uc := tryUseCase.(UseCase[I, O])
	s := NewScope(app.container)

	app.queryHandlersSetup(s, app.queryRegistry)
	app.commandHandlersSetup(s, app.commandRegistry)

	output, err := uc(input)
	if err != nil {
		return output, err
	}

	uow := GetService[Worker](s)
	uow.Commit()
	eventPublisher := GetService[*EventPublisher](s)
	for _, ev := range uow.DomainEvents() {
		eventPublisher.publish(ev)
	}

	return output, err
}
