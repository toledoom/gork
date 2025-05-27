package gork

type Entity interface {
	AddEvent(event Event)
	GetEvents() []Event
}

type Aggregate struct {
	Events []Event
}

func (aggregate *Aggregate) AddEvent(event Event) {
	aggregate.Events = append(aggregate.Events, event)
}

func (aggregate *Aggregate) GetEvents() []Event {
	return aggregate.Events
}
