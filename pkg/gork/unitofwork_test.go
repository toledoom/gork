package gork_test

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/toledoom/gork/pkg/gork"
)

func TestFetchMany(t *testing.T) {
	a := assert.New(t)
	storageMapper := gork.NewStorageMapper()
	fetchMany := func(filters ...gork.Filter) ([]gork.Entity, error) {
		inMemoryPersistence := map[string]*dumbEntity{
			"1": newDumbEntity("1", "value 1"),
			"2": newDumbEntity("2", "value 2"),
			"3": newDumbEntity("3", "value 3"),
		}
		var response []gork.Entity
		for _, f := range filters {
			id := f.(string)
			if de, ok := inMemoryPersistence[id]; ok {
				response = append(response, de)
			}
		}

		return response, nil
	}
	storageMapper.AddFetchManyFn(reflect.TypeOf(dumbEntity{}), fetchMany)

	sut := gork.NewUnitOfWork(storageMapper)

	result, err := sut.FetchMany(reflect.TypeOf(dumbEntity{}), "1")

	a.NoError(err)
	a.Len(result, 1)
	entity, ok := result[0].(*dumbEntity)
	a.True(ok)
	a.Equal("1", entity.ID)
	a.Equal("value 1", entity.Field1())
}
