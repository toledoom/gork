package gork

import "reflect"

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
	singletonServices   map[string]any
}

func newContainer() *Container {
	return &Container{
		serviceCollection:   make(map[string]any),
		singletonServices:   make(map[string]any),
		serviceLifetimeList: make(map[string]LifeTime),
	}
}

func RegisterService[T comparable](c *Container, builder Builder[T], l LifeTime) {
	serviceID := reflect.TypeOf((*T)(nil)).String()
	c.serviceCollection[serviceID] = builder
	c.serviceLifetimeList[serviceID] = l
}

func GetService[T comparable](s *Scope) T {
	serviceID := reflect.TypeOf((*T)(nil)).String()
	lt := s.c.serviceLifetimeList[serviceID]

	if lt == SINGLETON {
		_, ok := s.c.singletonServices[serviceID].(Builder[T])
		if !ok {
			s.c.singletonServices[serviceID] = s.c.serviceCollection[serviceID]
		}
		b := s.c.singletonServices[serviceID].(Builder[T])
		return b(s)
	}

	if lt == USECASE {
		_, ok := s.useCaseServices[serviceID].(Builder[T])
		if !ok {
			s.useCaseServices[serviceID] = s.c.serviceCollection[serviceID]

		}
		b := s.useCaseServices[serviceID].(Builder[T])
		return b(s)
	}

	b := s.c.serviceCollection[serviceID].(Builder[T])
	return b(s)
}

type Scope struct {
	c               *Container
	useCaseServices map[string]any
}

func NewScope(c *Container) *Scope {
	return &Scope{
		c:               c,
		useCaseServices: make(map[string]any),
	}
}
