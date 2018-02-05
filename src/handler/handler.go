package ginkgo

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"
)

//go:generate counterfeiter . Client

type Handler interface {
	Handle(string) error
}

type Client interface {
	Send(string)
}

type ActualClient struct {
	url string
}

func (ac *ActualClient) Send(input string) {
	log.Printf("Actual client is sending '%s' to %s", input, ac.url)
}

type RequestHandler struct {
	client Client
	name   string
}

func (rh *RequestHandler) Handle(input string) error {
	if len(input) == 0 {
		return errors.New("Empty input!")
	}
	if strings.HasSuffix(input, "?") {
		time.Sleep(3 * time.Second) // Do some complex stuff
	}
	rh.client.Send(fmt.Sprintf("%s handled by %s", input, rh.name))
	return nil
}

func NewHandler(name string, client Client) *RequestHandler {
	return &RequestHandler{
		name:   name,
		client: client,
	}
}

type AsyncHandler struct {
	handler Handler
	errChan chan error
}

func NewAsyncHandler(handler Handler, ch chan error) *AsyncHandler {
	return &AsyncHandler{
		handler: handler,
		errChan: ch,
	}
}

func (ah *AsyncHandler) Handle(input string) {
	go func() {
		time.Sleep(500 * time.Millisecond)
		ah.errChan <- ah.handler.Handle(input)
	}()
}
