package gork

import "github.com/hashicorp/go-multierror"

type Event interface {
	Name() string
}

type Handler interface {
	Notify(event Event) error
}

type EventPublisher struct {
	handlers map[string][]Handler
}

func NewPublisher() *EventPublisher {
	return &EventPublisher{
		handlers: make(map[string][]Handler),
	}
}

func (e *EventPublisher) Subscribe(handler Handler, events ...Event) {
	for _, event := range events {
		handlers := e.handlers[event.Name()]
		handlers = append(handlers, handler)
		e.handlers[event.Name()] = handlers
	}
}

func (e *EventPublisher) Publish(event Event) error {
	var multipleError error
	n := event.Name()
	for _, handler := range e.handlers[n] {
		err := handler.Notify(event)
		if err != nil {
			multipleError = multierror.Append(multipleError, err)
		}
	}
	return multipleError
}
