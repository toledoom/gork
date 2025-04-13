package gork

import (
	"fmt"
	"reflect"
)

type UseCase[I, O any] func(I) (O, error)

type UseCaseBuilder[I, O any] func(cr *CommandRegistry, qr *QueryRegistry) UseCase[I, O]

type UseCaseBuilderRegistry struct {
	useCaseBuilders map[string]any
}

func newUseCaseBuilderRegistry() *UseCaseBuilderRegistry {
	return &UseCaseBuilderRegistry{
		useCaseBuilders: make(map[string]any),
	}
}

func RegisterUseCaseBuilder[I, O any](ucbr *UseCaseBuilderRegistry, ucb UseCaseBuilder[I, O]) {
	var t I
	ucbr.useCaseBuilders[reflect.TypeOf(t).String()] = ucb
}

type UseCaseBuilderNotRegisteredError struct {
	ucb interface{}
}

func (ucnre *UseCaseBuilderNotRegisteredError) Error() string {
	return fmt.Sprintf("use case builder not registered for use case: %s", reflect.TypeOf(ucnre.ucb).String())
}

func ExecuteUseCase[I, O any](app *App, input I) (O, error) {
	tryUseCaseBuilder, ok := app.useCaseBuilderRegistry.useCaseBuilders[reflect.TypeOf(input).String()]
	if !ok {
		var r O
		return r, &UseCaseBuilderNotRegisteredError{ucb: input}
	}

	ucb := tryUseCaseBuilder.(UseCaseBuilder[I, O])
	s := NewScope(app.container)

	qr := newQueryRegistry()
	cr := newCommandRegistry()
	app.queryHandlersSetup(s, qr)
	app.commandHandlersSetup(s, cr)
	uc := ucb(cr, qr)

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
