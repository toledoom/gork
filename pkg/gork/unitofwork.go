package gork

import (
	"reflect"
)

type Worker interface {
	RegisterNew(newEntity Entity) error
	RegisterDirty(modifiedEntity Entity) error
	RegisterDeleted(deletedEntity Entity) error
	FetchOne(t reflect.Type, id string) (Entity, error)
	FetchMany(t reflect.Type, filters ...Filter) ([]Entity, error)
	Commit() error
	DomainEvents() []Event
}

type Filter interface { // TODO: Think about implementations of actual filters
	Condition() bool
}

type UnitOfWork struct {
	newEntities, dirtyEntities, deletedEntities []Entity
	dataMapper                                  *StorageMapper
}

func NewUnitOfWork(datamapper *StorageMapper) *UnitOfWork {
	return &UnitOfWork{
		dataMapper: datamapper,
	}
}

func (uow *UnitOfWork) RegisterNew(newEntity Entity) error {
	uow.newEntities = append(uow.newEntities, newEntity)
	return nil
}

func (uow *UnitOfWork) RegisterDirty(modifiedEntity Entity) error {
	uow.dirtyEntities = append(uow.dirtyEntities, modifiedEntity)
	return nil
}

func (uow *UnitOfWork) RegisterDeleted(deletedEntity Entity) error {
	uow.deletedEntities = append(uow.deletedEntities, deletedEntity)
	return nil
}

func (uow *UnitOfWork) FetchOne(t reflect.Type, id string) (Entity, error) {
	fn := uow.dataMapper.GetFetchOneFn(t)
	return fn(id)
}

func (uow *UnitOfWork) FetchMany(t reflect.Type, filters ...Filter) ([]Entity, error) {
	fn := uow.dataMapper.GetFetchManyFn(t)
	return fn(filters...)
}

func (uow *UnitOfWork) Commit() error {
	for _, en := range uow.newEntities {
		fn := uow.dataMapper.GetMutationFn(reflect.TypeOf(en), CreationQuery)
		err := fn(en)
		if err != nil {
			return nil
		}
	}

	for _, en := range uow.dirtyEntities {
		fn := uow.dataMapper.GetMutationFn(reflect.TypeOf(en), UpdateQuery)
		err := fn(en)
		if err != nil {
			return nil
		}
	}

	for _, en := range uow.deletedEntities {
		fn := uow.dataMapper.GetMutationFn(reflect.TypeOf(en), DeletionQuery)
		err := fn(en)
		if err != nil {
			return nil
		}
	}

	return nil
}

func (uow *UnitOfWork) DomainEvents() []Event {
	var events []Event
	for _, e := range uow.newEntities {
		events = append(events, e.GetEvents()...)
	}
	for _, e := range uow.dirtyEntities {
		events = append(events, e.GetEvents()...)
	}
	for _, e := range uow.deletedEntities {
		events = append(events, e.GetEvents()...)
	}

	return events
}
