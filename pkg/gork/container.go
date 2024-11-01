package gork

import "reflect"

type Builder[T any] func(*Container) T

type Scope int32

const (
	SERVER_SCOPE  Scope = 0
	USECASE_SCOPE Scope = 1
)

type Container struct {
	serverServices  map[string]any
	useCaseServices map[string]any
	cache           map[string]any
}

func newContainer() *Container {
	return &Container{
		serverServices:  make(map[string]any),
		useCaseServices: make(map[string]any),
		cache:           make(map[string]any),
	}
}

func AddService[T comparable](c *Container, builder Builder[T], s Scope) {
	serviceID := reflect.TypeOf((*T)(nil)).String()
	if s == SERVER_SCOPE {
		c.serverServices[serviceID] = builder
		return
	}
	c.useCaseServices[serviceID] = builder
}

func GetService[T comparable](c *Container) T {
	var t T
	serviceID := reflect.TypeOf((*T)(nil)).String()

	serviceBuilder, ok := c.useCaseServices[serviceID]
	if ok {
		ucServiceBuilder := serviceBuilder.(Builder[T])
		t = ucServiceBuilder(c)
		return t
	}

	t, ok = c.cache[serviceID].(T)
	if ok {
		return t
	}
	serverServiceBuilder := c.serverServices[serviceID].(Builder[T])
	t = serverServiceBuilder(c)
	c.cache[serviceID] = t
	return t
}
