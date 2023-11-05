package di

import "reflect"

type Builder[T any] func(*Container) T

type Container struct {
	services map[string]any
}

func NewContainer() *Container {
	return &Container{
		services: make(map[string]any),
	}
}

func AddService[T comparable](c *Container, builder Builder[T]) {
	var t T
	serviceID := reflect.TypeOf(t).String()
	c.services[serviceID] = builder
}

func GetService[T comparable](c *Container) T {
	var t T
	serviceID := reflect.TypeOf(t).String()
	builder := c.services[serviceID].(Builder[T])
	return builder(c)
}
