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

type MutationFn func(e Entity) error
type FetchOneFn func(id string) (Entity, error)
type FetchManyFn func(filters ...Filter) ([]Entity, error)

type StorageMapper struct {
	mutationFns  map[string]MutationFn
	fetchOneFns  map[string]FetchOneFn
	fetchManyFns map[string]FetchManyFn
}

func newStorageMapper() *StorageMapper {
	return &StorageMapper{
		mutationFns:  make(map[string]MutationFn),
		fetchOneFns:  make(map[string]FetchOneFn),
		fetchManyFns: make(map[string]FetchManyFn),
	}
}

func (sm *StorageMapper) AddMutationFn(t reflect.Type, entityType int, fn MutationFn) {
	sm.mutationFns[fmt.Sprintf("%s-%d", t.String(), entityType)] = fn
}

func (sm *StorageMapper) GetMutationFn(t reflect.Type, entityType int) MutationFn {
	return sm.mutationFns[fmt.Sprintf("%s-%d", t.String(), entityType)]
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
