package persistence

import (
	"fmt"
	"reflect"

	"github.com/toledoom/gork/pkg/entity"
)

const (
	EntityNew = iota
	EntityDirty
	EntityDeleted
)

type DataMapper struct {
	persistenceFns map[string]PersistenceFn
	fetchFns       map[string]FetchFn
}

func (dm *DataMapper) AddPersistenceFn(t reflect.Type, entityType int, fn PersistenceFn) {
	dm.persistenceFns[fmt.Sprintf("%s-%d", t.String(), entityType)] = fn
}

func (dm *DataMapper) GetPersistenceFn(t reflect.Type, entityType int) PersistenceFn {
	return dm.persistenceFns[fmt.Sprintf("%s-%d", t.String(), entityType)]
}

func (dm *DataMapper) AddFetchFn(t reflect.Type, fn FetchFn) {
	dm.fetchFns[t.String()] = fn
}

func (dm *DataMapper) GetFetchFn(t reflect.Type) FetchFn {
	return dm.fetchFns[t.String()]
}

type PersistenceFn func(e entity.Entity) error
type FetchFn func(id string) (entity.Entity, error)

func NewDataMapper() *DataMapper {
	return &DataMapper{
		persistenceFns: make(map[string]PersistenceFn),
	}
}
