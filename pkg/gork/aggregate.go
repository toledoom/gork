package gork

type Entity interface {
	AddEvent(e Event)
	GetEvents() []Event
}

type Aggregate struct {
	Events []Event
}

func (a *Aggregate) AddEvent(e Event) {
	a.Events = append(a.Events, e)
}

func (a *Aggregate) GetEvents() []Event {
	return a.Events
}
