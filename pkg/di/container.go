package di

type Service func() any

type Container struct {
	services map[string]Service
}

func (c *Container) Add(name string, s Service) {
	c.services[name] = s
}

func (c *Container) Get(name string) Service {
	return c.services[name]
}

func NewContainer() *Container {
	return &Container{
		services: make(map[string]Service),
	}
}

func Add[T comparable](c *Container, name string, s Service) {

}

/////////////////////////////////////////////////////
