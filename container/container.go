package container

import "go.uber.org/dig"

var Default Container

type Container interface {
	Invoke(function interface{}) error
	MustInvoke(function interface{})
}

type Builder interface {
	Container
	Provide(constructor interface{}) error
	MustProvide(constructor interface{})
}

func New() Builder {
	return &container{
		Container: dig.New(),
	}
}

type container struct {
	*dig.Container
}

func (c *container) Invoke(function interface{}) error {
	return c.Container.Invoke(function)
}

func (c *container) MustInvoke(function interface{}) {
	if err := c.Invoke(function); err != nil {
		panic(err)
	}
}

func (c *container) Provide(constructor interface{}) error {
	return c.Container.Provide(constructor)
}

func (c *container) MustProvide(constructor interface{}) {
	if err := c.Provide(constructor); err != nil {
		panic(err)
	}
}

func Get[T any](container ...Container) (T, error) {
	if len(container) == 0 {
		container = []Container{Default}
	}

	var err error
	var val T
	for _, c := range container {
		err = c.Invoke(func(a T) {
			val = a
		})

		if err == nil {
			return val, nil
		}
	}
	return val, err
}

func MustGet[T any](container ...Container) (val T) {
	val, err := Get[T](container...)
	if err != nil {
		panic(err)
	}

	return val
}
