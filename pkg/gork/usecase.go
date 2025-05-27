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

func RegisterUseCaseBuilder[I, O any](useCaseBuilderRegistry *UseCaseBuilderRegistry, useCaseBuilder UseCaseBuilder[I, O]) {
	var t I
	useCaseBuilderRegistry.useCaseBuilders[reflect.TypeOf(t).String()] = useCaseBuilder
}

type UseCaseBuilderNotRegisteredError struct {
	useCaseBuilder any
}

func (ucnre *UseCaseBuilderNotRegisteredError) Error() string {
	return fmt.Sprintf("use case builder not registered for use case: %s", reflect.TypeOf(ucnre.useCaseBuilder).String())
}

func ExecuteUseCase[I, O any](app *App, input I) (O, error) {
	tryUseCaseBuilder, ok := app.useCaseBuilderRegistry.useCaseBuilders[reflect.TypeOf(input).String()]
	if !ok {
		var r O
		return r, &UseCaseBuilderNotRegisteredError{useCaseBuilder: input}
	}

	useCaseBuilder := tryUseCaseBuilder.(UseCaseBuilder[I, O])
	scope := NewScope(app.container)

	queryRegistry := newQueryRegistry()
	commandRegistry := newCommandRegistry()
	app.queryHandlersSetup(scope, queryRegistry)
	app.commandHandlersSetup(scope, commandRegistry)
	useCase := useCaseBuilder(commandRegistry, queryRegistry)

	output, err := useCase(input)
	if err != nil {
		return output, err
	}

	unitOfWork := GetService[Worker](scope)
	unitOfWork.Commit()
	eventPublisher := GetService[*EventPublisher](scope)
	for _, ev := range unitOfWork.DomainEvents() {
		eventPublisher.publish(ev)
	}

	return output, err
}
