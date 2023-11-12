package persistence

import (
	"reflect"

	"github.com/toledoom/gork/pkg/entity"
	"github.com/toledoom/gork/pkg/event"
)

type Worker interface {
	RegisterNew(newEntity entity.Entity) error
	RegisterDirty(modifiedEntity entity.Entity) error
	RegisterDeleted(deletedEntity entity.Entity) error
	FetchByID(t reflect.Type, id string) (entity.Entity, error)
	Commit() error
	DomainEvents() []event.Event
}

type UnitOfWork struct {
	newEntities, dirtyEntities, deletedEntities []entity.Entity
	dataMapper                                  *StorageMapper
}

func NewUnitOfWork(datamapper *StorageMapper) *UnitOfWork {
	return &UnitOfWork{
		dataMapper: datamapper,
	}
}

func (uow *UnitOfWork) RegisterNew(newEntity entity.Entity) error {
	uow.newEntities = append(uow.newEntities, newEntity)
	return nil
}

func (uow *UnitOfWork) RegisterDirty(modifiedEntity entity.Entity) error {
	uow.dirtyEntities = append(uow.dirtyEntities, modifiedEntity)
	return nil
}

func (uow *UnitOfWork) RegisterDeleted(deletedEntity entity.Entity) error {
	uow.deletedEntities = append(uow.deletedEntities, deletedEntity)
	return nil
}

func (uow *UnitOfWork) FetchByID(t reflect.Type, id string) (entity.Entity, error) {
	fn := uow.dataMapper.GetFetchFn(t)
	return fn(id)
}

func (uow *UnitOfWork) Commit() error {
	for _, en := range uow.newEntities {
		fn := uow.dataMapper.GetPersistenceFn(reflect.TypeOf(en), EntityNew)
		err := fn(en)
		if err != nil {
			return nil
		}
	}

	for _, en := range uow.dirtyEntities {
		fn := uow.dataMapper.GetPersistenceFn(reflect.TypeOf(en), EntityDirty)
		err := fn(en)
		if err != nil {
			return nil
		}
	}

	for _, en := range uow.deletedEntities {
		fn := uow.dataMapper.GetPersistenceFn(reflect.TypeOf(en), EntityDeleted)
		err := fn(en)
		if err != nil {
			return nil
		}
	}

	return nil
}

func (uow *UnitOfWork) DomainEvents() []event.Event {
	var events []event.Event
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
