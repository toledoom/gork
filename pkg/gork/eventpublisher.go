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

func NewPublisher() *EventPublisher {
	return &EventPublisher{
		handlers: make(map[string][]eventHandler),
	}
}

func (ep *EventPublisher) Subscribe(handler eventHandler, events ...Event) {
	for _, event := range events {
		handlers := ep.handlers[event.Name()]
		handlers = append(handlers, handler)
		ep.handlers[event.Name()] = handlers
	}
}

func (ep *EventPublisher) publish(event Event) error {
	var multipleError error
	eventName := event.Name()
	for _, handler := range ep.handlers[eventName] {
		err := handler.Handle(event)
		if err != nil {
			multipleError = multierror.Append(multipleError, err)
		}
	}
	return multipleError
}
