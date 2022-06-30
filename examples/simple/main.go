package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/mbict/go-webapp"
	"github.com/mbict/go-webapp/app"
	"log"
	"net/http"
)

func Empty(context.Context, *webapp.Empty) (interface{}, error) {
	return nil, nil
}

func Ping(context.Context, *webapp.Empty) (string, error) {
	return "pong", nil
}

func Test1(context.Context, *webapp.Empty) (string, error) {
	return "test1", nil
}

func Test2(context.Context, *webapp.Empty) (string, error) {
	return "test2", nil
}

func Test3(context.Context, *webapp.Empty) (string, error) {

	return "test3", nil
}

func Test4(context.Context, *webapp.Empty) (string, error) {

	return "test4", webapp.Error(errors.New("yup nicht gefunt"), http.StatusNotFound)
}

type GreetRequest struct {
	Name string `json:"name" query:"name" default:"world"` // Get name from JSON body.
	Age  int    `header:"X-User-Age"`                      // Get age from HTTP header.
}

type GreetResponse struct {
	Greeting string `json:"data"`
}

// Set a custom HTTP response code.
func (gr *GreetResponse) StatusCode() int {
	return http.StatusTeapot
}

// Add custom headers to the response.
func (gr *GreetResponse) Header() http.Header {
	header := http.Header{}
	header.Set("foo", "bar")
	return header
}

func Greet(ctx context.Context, req *GreetRequest) (*GreetResponse, error) {
	if req.Name == "" {
		return nil, webapp.ErrBadRequest
	}

	res := &GreetResponse{
		Greeting: fmt.Sprintf("Hello %s, you're %d years old.", req.Name, req.Age),
	}

	return res, nil
}

func main() {
	r := webapp.New(app.NewContainer())

	r.Use(func(next http.HandlerFunc) http.HandlerFunc {
		return func(rw http.ResponseWriter, req *http.Request) {
			fmt.Println("global middleware first")
			next(rw, req)
		}
	})

	r.Get("/", webapp.H(Greet))
	r.Post("/", webapp.H(Greet))
	r.Get("/test", webapp.H(Ping), func(next http.HandlerFunc) http.HandlerFunc {
		return func(rw http.ResponseWriter, req *http.Request) {
			fmt.Println("route middleware second")
			next(rw, req)
		}
	})
	r.Get("/test/@id:test", webapp.H(Test1))
	r.Get("/test/@id:test1", webapp.H(Test2))
	r.Get("/test/@id:test2", webapp.H(Test3))

	testGroup := r.Group("/test2/@id", func(next http.HandlerFunc) http.HandlerFunc {
		return func(rw http.ResponseWriter, req *http.Request) {
			fmt.Println("group middleware second")
			next(rw, req)
		}
	})
	{
		testGroup.Get("", webapp.H(Ping), func(next http.HandlerFunc) http.HandlerFunc {
			return func(rw http.ResponseWriter, req *http.Request) {
				fmt.Println("group route middleware last")
				next(rw, req)
			}
		})
		testGroup.Get(":test1", webapp.H(Test1))
		testGroup.Get(":test2", webapp.H(Test2))
		testGroup.Get(":test3", webapp.H(Test3))
		testGroup.Get(":test4", webapp.H(Test4))

		nestedGroup := testGroup.Group("/group", func(next http.HandlerFunc) http.HandlerFunc {
			return func(rw http.ResponseWriter, req *http.Request) {
				fmt.Println("nested group route middleware ")
				next(rw, req)
			}
		})
		{
			nestedGroup.Get("", webapp.H(Test1), func(next http.HandlerFunc) http.HandlerFunc {
				return func(rw http.ResponseWriter, req *http.Request) {
					fmt.Println("nested group path middleware last")
					next(rw, req)
				}
			})
		}
	}

	r.Get("/test/@id", webapp.H(Ping))

	r.Post("/created", webapp.H(func(ctx context.Context, empty *webapp.Empty) (webapp.CreatedResponse, error) {
		//return webapp.EmptyResponse, nil
		return webapp.NewCreatedResponse("/test"), nil
	}))

	r.Get("/notsocreated", webapp.H(func(ctx context.Context, empty *webapp.Empty) (webapp.CreatedResponse, error) {
		//return webapp.EmptyResponse, nil
		return webapp.NewCreatedResponse("/test"), errors.New("serious problem here. ;)")
	}))

	r.Get("/empty", webapp.H(func(ctx context.Context, empty *webapp.Empty) (*webapp.Empty, error) {
		return nil, nil
	}))

	log.Fatal(r.ListenAndServe(":8080"))
}
