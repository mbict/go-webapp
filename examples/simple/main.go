package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/mbict/webapp"
	"github.com/mbict/webapp/app"
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
	Name string `json:"name"`         // Get name from JSON body.
	Age  int    `header:"X-User-Age"` // Get age from HTTP header.
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

	r.Get("/test", webapp.H(Ping))
	r.Get("/test/@id:test", webapp.H(Test1))
	r.Get("/test/@id:test1", webapp.H(Test2))
	r.Get("/test/@id:test2", webapp.H(Test3))

	testGroup := r.Group("/test2/@id")
	{
		testGroup.Get("", webapp.H(Ping))
		testGroup.Get(":test1", webapp.H(Test1))
		testGroup.Get(":test2", webapp.H(Test2))
		testGroup.Get(":test3", webapp.H(Test3))
		testGroup.Get(":test4", webapp.H(Test4))
	}

	r.Get("/test/@id", webapp.H(Ping))

	log.Fatal(r.ListenAndServe(":8080"))
}
