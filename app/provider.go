package app

import (
	"github.com/mbict/go-webapp/container"
)

var Container container.Container

func NewContainer() container.Container {
	c := container.New()

	//c.MustProvide(providesArgumentsDecoderBuilder)

	return c
}

//func providesArgumentsDecoderBuilder() webapp.ArgumentsDecoderBuilder {
//	return webapp.ArgumentsDecoderBuilder(webapp.NewDecoderBuilder(
//		decoder.NewHeaderDecoderBuilder(),
//		decoder.NewQueryDecoderBuilder(),
//		decoder.NewPathDecoderBuilder(),
//	))
//}
