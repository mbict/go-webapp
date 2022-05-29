package main

import (
	"context"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/mbict/go-commandbus/v2"
	"github.com/mbict/go-querybus"
	"github.com/mbict/webapp"
	"log"
	"net/http"
)

var (
	queryBus   = querybus.New()
	commandBus = commandbus.New()
	validate   = validator.New()
)

//Example for the CommandBus
type CreateCommand struct {
	Id   uuid.UUID `json:"-" validate:"required"`
	Name string    `json:"name" validate:"required"`
}

func (c CreateCommand) CommandName() string {
	return "example.create.command"
}

type UpdateCommand struct {
	Id   uuid.UUID `path:"id" json:"-" validate:"required"`
	Name string    `json:"name" validate:"required"`
}

func (c UpdateCommand) CommandName() string {
	return "example.update.command"
}

//Example for the QueryBus
type QueryExample struct {
	Id string `path:"id"`

	WithName *string `query:"with:name"`

	Offset int `query:"offset" validate:"gte=0"`
	Size   int `query:"size" validate:"gte=1,lte=1000" default:"100"`
}

type Metadata struct {
	Offset int `json:"offset"`
	Size   int `json:"size"`
	Total  int `json:"total"`
}

type SubResults struct {
	Id   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type QueryResult struct {
	Metadata Metadata     `json:"metadata"`
	Data     []SubResults `json:"data"`
}

func (q QueryExample) QueryName() string {
	return "example.query"
}

func HandleCommand[T commandbus.Command]() http.HandlerFunc {
	return webapp.H(func(ctx context.Context, cmd T) (interface{}, error) {
		//validate
		if err := validate.Struct(cmd); err != nil {
			//validation failed
			return nil, err
		}

		//push the command onto the commandBus
		err := commandBus.Handle(ctx, cmd)
		if err != nil {
			//throw the error, command failed
			return nil, err
		}
		return webapp.EmptyResponse, nil
	})
}

func HandleCreateCommand[T commandbus.Command](gen func(cmd *T) string) http.HandlerFunc {

	return webapp.H(func(ctx context.Context, cmd T) (interface{}, error) {
		resourceLocation := gen(&cmd)

		//validate
		if err := validate.Struct(cmd); err != nil {
			//validation failed
			return nil, err
		}

		//push the command onto the commandBus
		err := commandBus.Handle(ctx, cmd)
		if err != nil {
			//throw the error, command failed
			return nil, err
		}
		return webapp.NewCreatedResponse(resourceLocation), nil
	})
}

func HandleQuery[T any]() http.HandlerFunc {
	return webapp.H(func(ctx context.Context, query T) (interface{}, error) {
		//validate
		if err := validate.Struct(query); err != nil {
			//validation failed
			return nil, err
		}

		//push the query into the querybus
		return queryBus.Handle(ctx, query)
	})

}

func main() {
	must(commandBus.Register(CreateCommand{}, commandbus.CommandHandlerFunc(func(ctx context.Context, command interface{}) error {
		return nil
	})))

	must(commandBus.Register(UpdateCommand{}, commandbus.CommandHandlerFunc(func(ctx context.Context, command interface{}) error {
		fmt.Println(command.(UpdateCommand).Id)
		return nil
	})))

	must(querybus.RegisterHandler(queryBus, func(ctx context.Context, query QueryExample) (QueryResult, error) {
		return QueryResult{
			Metadata: Metadata{
				Offset: 123,
				Size:   456,
				Total:  87359837,
			},
			Data: []SubResults{
				{
					Id:   uuid.New(),
					Name: "abc",
				},
				{
					Id:   uuid.New(),
					Name: "efg hijklmmn",
				},
			},
		}, nil
	}))

	r := webapp.New(nil)

	r.Get("/res/@id", HandleQuery[QueryExample]())

	//example of a creational command that will return a location header
	r.Post("/res", HandleCreateCommand[CreateCommand](func(cmd *CreateCommand) string {
		cmd.Id = uuid.New()
		return "http://localhost/res/" + cmd.Id.String()
	}))

	//normal commands that only process the request and do not return any information
	r.Put("/res/@id", HandleCommand[UpdateCommand]())

	log.Println(r.ListenAndServe(":8080"))
}

func must(e error) {
	if e != nil {
		panic(e)
	}
}
