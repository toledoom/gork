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
	c               *Container
	useCaseServices map[string]any
	id              uint64
}

func NewScope(c *Container) *Scope {
	return &Scope{
		c:               c,
		useCaseServices: make(map[string]any),
		id:              rand.Uint64(),
	}
}

func RegisterService[T comparable](c *Container, builder Builder[T], l LifeTime) {
	serviceID := reflect.TypeOf((*T)(nil)).String()
	c.serviceCollection[serviceID] = builder
	c.serviceLifetimeList[serviceID] = l
}

func GetService[T comparable](s *Scope) T {
	serviceID := reflect.TypeOf((*T)(nil)).String()
	lifeTime := s.c.serviceLifetimeList[serviceID]

	if lifeTime == SINGLETON {
		s.c.mutex.RLock()
		t, ok := s.c.singletonServices[serviceID].(T)
		s.c.mutex.RUnlock()
		if ok {
			return t
		}
		builder := s.c.serviceCollection[serviceID].(Builder[T])
		s.c.mutex.Lock()
		s.c.singletonServices[serviceID] = builder(s)
		s.c.mutex.Unlock()

		return s.c.singletonServices[serviceID].(T)
	}

	if lifeTime == USECASE {
		t, ok := s.useCaseServices[serviceID].(T)
		if ok {
			return t

		}
		b := s.c.serviceCollection[serviceID].(Builder[T])
		s.useCaseServices[serviceID] = b(s)

		return s.useCaseServices[serviceID].(T)
	}

	b := s.c.serviceCollection[serviceID].(Builder[T])
	return b(s)
}

func newContainer() *Container {
	return &Container{
		serviceCollection:   make(map[string]any),
		singletonServices:   make(map[string]any),
		serviceLifetimeList: make(map[string]LifeTime),
		mutex:               sync.RWMutex{},
	}
}
