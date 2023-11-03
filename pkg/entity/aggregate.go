package entity

import "github.com/toledoom/gork/pkg/event"

type Entity interface {
	AddEvent(e event.Event)
	GetEvents() []event.Event
}

type Aggregate struct {
	Events []event.Event
}

func (a *Aggregate) AddEvent(e event.Event) {
	a.Events = append(a.Events, e)
}

func (a *Aggregate) GetEvents() []event.Event {
	return a.Events
}
