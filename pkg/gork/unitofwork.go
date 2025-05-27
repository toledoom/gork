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

type Filter any

type UnitOfWork struct {
	newEntities, dirtyEntities, deletedEntities []Entity
	storageMapper                               *StorageMapper
}

func NewUnitOfWork(storagemapper *StorageMapper) *UnitOfWork {
	return &UnitOfWork{
		storageMapper: storagemapper,
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
	fn := uow.storageMapper.GetFetchOneFn(t)
	return fn(id)
}

func (uow *UnitOfWork) FetchMany(t reflect.Type, filters ...Filter) ([]Entity, error) {
	fn := uow.storageMapper.GetFetchManyFn(t)
	return fn(filters...)
}

func (uow *UnitOfWork) Commit() error {
	for _, newEntity := range uow.newEntities {
		fn := uow.storageMapper.GetMutationFn(reflect.TypeOf(newEntity), CreationQuery)
		err := fn(newEntity)
		if err != nil {
			return err
		}
	}

	for _, dirtyEntity := range uow.dirtyEntities {
		fn := uow.storageMapper.GetMutationFn(reflect.TypeOf(dirtyEntity), UpdateQuery)
		err := fn(dirtyEntity)
		if err != nil {
			return err
		}
	}

	for _, deletedEntity := range uow.deletedEntities {
		fn := uow.storageMapper.GetMutationFn(reflect.TypeOf(deletedEntity), DeletionQuery)
		err := fn(deletedEntity)
		if err != nil {
			return err
		}
	}

	return nil
}

func (uow *UnitOfWork) DomainEvents() []Event {
	var events []Event
	for _, entity := range uow.newEntities {
		events = append(events, entity.GetEvents()...)
	}
	for _, entity := range uow.dirtyEntities {
		events = append(events, entity.GetEvents()...)
	}
	for _, entity := range uow.deletedEntities {
		events = append(events, entity.GetEvents()...)
	}

	return events
}
