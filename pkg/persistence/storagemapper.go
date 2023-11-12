package persistence

import (
	"fmt"
	"reflect"

	"github.com/toledoom/gork/pkg/entity"
)

const (
	CreationQuery = iota
	UpdateQuery
	DeletionQuery
	FetchOneQUery
	FetchManyQuery
)

type PersistenceFn func(e entity.Entity) error
type FetchFn func(id string) (entity.Entity, error)

type StorageMapper struct {
	persistenceFns map[string]PersistenceFn
	fetchFns       map[string]FetchFn
}

func NewStorageMapper() *StorageMapper {
	return &StorageMapper{
		persistenceFns: make(map[string]PersistenceFn),
	}
}

func (sm *StorageMapper) AddPersistenceFn(t reflect.Type, entityType int, fn PersistenceFn) {
	sm.persistenceFns[fmt.Sprintf("%s-%d", t.String(), entityType)] = fn
}

func (sm *StorageMapper) GetPersistenceFn(t reflect.Type, entityType int) PersistenceFn {
	return sm.persistenceFns[fmt.Sprintf("%s-%d", t.String(), entityType)]
}

func (sm *StorageMapper) AddFetchFn(t reflect.Type, fn FetchFn) {
	sm.fetchFns[t.String()] = fn
}

func (sm *StorageMapper) GetFetchFn(t reflect.Type) FetchFn {
	return sm.fetchFns[t.String()]
}
