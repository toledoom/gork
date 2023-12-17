package gork

import (
	"fmt"
	"reflect"
)

const (
	CreationQuery = iota
	UpdateQuery
	DeletionQuery
	FetchOneQUery
	FetchManyQuery
)

type RepositoriesSetup func(container *Container, uow Worker)

type PersistenceFn func(e Entity) error
type FetchOneFn func(id string) (Entity, error)
type FetchManyFn func(filters ...Filter) ([]Entity, error)

type StorageMapper struct {
	persistenceFns map[string]PersistenceFn
	fetchOneFns    map[string]FetchOneFn
	fetchManyFns   map[string]FetchManyFn
}

func NewStorageMapper() *StorageMapper {
	return &StorageMapper{
		persistenceFns: make(map[string]PersistenceFn),
		fetchOneFns:    make(map[string]FetchOneFn),
		fetchManyFns:   make(map[string]FetchManyFn),
	}
}

func (sm *StorageMapper) AddPersistenceFn(t reflect.Type, entityType int, fn PersistenceFn) {
	sm.persistenceFns[fmt.Sprintf("%s-%d", t.String(), entityType)] = fn
}

func (sm *StorageMapper) GetPersistenceFn(t reflect.Type, entityType int) PersistenceFn {
	return sm.persistenceFns[fmt.Sprintf("%s-%d", t.String(), entityType)]
}

func (sm *StorageMapper) AddFetchOneFn(t reflect.Type, fn FetchOneFn) {
	sm.fetchOneFns[t.String()] = fn
}

func (sm *StorageMapper) GetFetchOneFn(t reflect.Type) FetchOneFn {
	return sm.fetchOneFns[t.String()]
}

func (sm *StorageMapper) AddFetchManyFn(t reflect.Type, fn FetchManyFn) {
	sm.fetchManyFns[t.String()] = fn
}

func (sm *StorageMapper) GetFetchManyFn(t reflect.Type) FetchManyFn {
	return sm.fetchManyFns[t.String()]
}
