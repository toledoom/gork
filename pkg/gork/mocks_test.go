package gork_test

import (
	"reflect"

	"github.com/toledoom/gork/pkg/gork"
)

type dumbCommand struct {
	ID string
}

func dumbCommandHandler(dc *dumbCommand) error { return nil }

type dumbQuery struct{}

func dumbQueryHandler(dc *dumbQuery) (string, error) { return "a value", nil }

func persistEntityCommandHandler(repository *dumbEntityUowRepository) func(dc *dumbCommand) error {
	return func(dc *dumbCommand) error {
		de := newDumbEntity(dc.ID)
		return repository.Add(de)
	}
}

type dumbEntityUowRepository struct {
	uow gork.Worker
}

func (der *dumbEntityUowRepository) FindByID(id string) (*dumbEntity, error) {
	entity, err := der.uow.FetchOne(reflect.TypeOf(&dumbEntity{}), id)
	if err != nil {
		return nil, err
	}
	d := entity.(*dumbEntity)

	return d, nil
}

func (der *dumbEntityUowRepository) Add(d *dumbEntity) error {
	return der.uow.RegisterNew(d)
}

type dumbEvent struct {
	gork.Event

	dumbID string
}

func (de *dumbEvent) Name() string {
	return "DumbEvent"
}

func newDumbEntity(id string) *dumbEntity {
	dEnt := &dumbEntity{
		ag: &gork.Aggregate{},
		ID: id,
	}
	dEnt.AddEvent(&dumbEvent{dumbID: id})
	return dEnt
}

func (de *dumbEntity) AddEvent(e gork.Event) {
	de.ag.Events = append(de.ag.Events, e)
}

func (a *dumbEntity) GetEvents() []gork.Event {
	return a.ag.Events
}

type dumbEntity struct {
	ag *gork.Aggregate

	ID string
}

func dumbUseCase(cr *gork.CommandRegistry, qr *gork.QueryRegistry) func(dumbUseCaseInput) (dumbUseCaseOutput, error) {
	return func(dumbUseCaseInput) (dumbUseCaseOutput, error) {
		dc := &dumbCommand{}

		err := gork.HandleCommand(cr, dc)
		if err != nil {
			return dumbUseCaseOutput{}, err
		}

		dq := &dumbQuery{}
		resp, err := gork.HandleQuery[*dumbQuery, string](qr, dq)
		if err != nil {
			return dumbUseCaseOutput{}, err
		}

		return dumbUseCaseOutput{
			response: resp,
		}, nil
	}
}

type dumbUseCaseInput struct{}

type dumbUseCaseOutput struct {
	response string
}
