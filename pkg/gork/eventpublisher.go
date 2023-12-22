package gork

import "github.com/hashicorp/go-multierror"

type Event interface {
	Name() string
}

type eventHandler interface {
	Handle(event Event) error
}

type EventPublisher struct {
	handlers map[string][]eventHandler
}

func newPublisher() *EventPublisher {
	return &EventPublisher{
		handlers: make(map[string][]eventHandler),
	}
}

func (e *EventPublisher) Subscribe(handler eventHandler, events ...Event) {
	for _, event := range events {
		handlers := e.handlers[event.Name()]
		handlers = append(handlers, handler)
		e.handlers[event.Name()] = handlers
	}
}

func (e *EventPublisher) publish(event Event) error {
	var multipleError error
	n := event.Name()
	for _, handler := range e.handlers[n] {
		err := handler.Handle(event)
		if err != nil {
			multipleError = multierror.Append(multipleError, err)
		}
	}
	return multipleError
}
