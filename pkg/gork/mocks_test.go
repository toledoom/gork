package gork_test

import (
	"errors"
	"reflect"

	"github.com/toledoom/gork/pkg/gork"
)

type dumbCommand struct {
	ID string
}

func dumbCommandHandler(dc *dumbCommand) error { return nil }

type dumbQuery struct{}

func dumbQueryHandler(dc *dumbQuery) (string, error) { return "a value", nil }

func dumbUseCase(cr *gork.CommandRegistry, qr *gork.QueryRegistry) gork.UseCase[dumbUseCaseInput, dumbUseCaseOutput] {
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

type dumbUnitOfWork struct {
	gork.Worker

	entities []gork.Entity
}

func (uow *dumbUnitOfWork) RegisterNew(newEntity gork.Entity) error {
	uow.entities = append(uow.entities, newEntity)
	return nil
}

func (uow *dumbUnitOfWork) FetchOne(t reflect.Type, id string) (gork.Entity, error) {
	for _, v := range uow.entities {
		e := v.(*dumbEntity)
		if e.ID == id {
			return e, nil
		}
	}
	return nil, errors.New("entity not found")
}

func (uow *dumbUnitOfWork) Commit() error {
	return nil
}

func (uow *dumbUnitOfWork) DomainEvents() []gork.Event {
	return nil
}
