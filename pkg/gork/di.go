package gork

import (
	"math/rand"
	"reflect"
	"sync"
)

type Builder[T any] func(*Scope) T

type LifeTime int32

const (
	SINGLETON LifeTime = 0
	USECASE   LifeTime = 1
	TRANSIENT LifeTime = 2
)

type Container struct {
	serviceCollection   map[string]any
	serviceLifetimeList map[string]LifeTime

	mutex             sync.RWMutex
	singletonServices map[string]any
}

type Scope struct {
	container       *Container
	useCaseServices map[string]any
	id              uint64
}

func NewScope(container *Container) *Scope {
	return &Scope{
		container:       container,
		useCaseServices: make(map[string]any),
		id:              rand.Uint64(),
	}
}

func RegisterService[T comparable](container *Container, builder Builder[T], lifeTime LifeTime) {
	serviceID := reflect.TypeOf((*T)(nil)).String()
	container.serviceCollection[serviceID] = builder
	container.serviceLifetimeList[serviceID] = lifeTime
}

func GetService[T comparable](scope *Scope) T {
	serviceID := reflect.TypeOf((*T)(nil)).String()
	lifeTime := scope.container.serviceLifetimeList[serviceID]

	if lifeTime == SINGLETON {
		scope.container.mutex.RLock()
		service, ok := scope.container.singletonServices[serviceID].(T)
		scope.container.mutex.RUnlock()
		if ok {
			return service
		}
		builder := scope.container.serviceCollection[serviceID].(Builder[T])
		scope.container.mutex.Lock()
		scope.container.singletonServices[serviceID] = builder(scope)
		scope.container.mutex.Unlock()

		return scope.container.singletonServices[serviceID].(T)
	}

	if lifeTime == USECASE {
		service, ok := scope.useCaseServices[serviceID].(T)
		if ok {
			return service

		}
		builder := scope.container.serviceCollection[serviceID].(Builder[T])
		scope.useCaseServices[serviceID] = builder(scope)

		return scope.useCaseServices[serviceID].(T)
	}

	builder := scope.container.serviceCollection[serviceID].(Builder[T])
	return builder(scope)
}

func newContainer() *Container {
	return &Container{
		serviceCollection:   make(map[string]any),
		singletonServices:   make(map[string]any),
		serviceLifetimeList: make(map[string]LifeTime),
		mutex:               sync.RWMutex{},
	}
}
