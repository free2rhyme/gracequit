package main

import (
	"errors"
	"fmt"
)

type request struct {
	carriedReq interface{}
	reply      chan interface{}
}

type onEvent1Request struct {
	data interface{}
}

type onEvent2Request struct {
	data interface{}
}

type onEvent3Request struct {
	data interface{}
}

type Handler struct {
	requests chan *request
	quit     chan bool
	stopped  chan bool
	shutdown chan bool
}

func GetHandler() *Handler {
	h := &Handler{
		requests: make(chan *request, 64),
		quit:     make(chan bool, 1),
		stopped:  make(chan bool, 1),
		shutdown: make(chan bool, 1),
	}
	return h
}

func (h *Handler) Start() error {

	go h.run()

	return nil
}

func (h *Handler) Stop() {
	h.quit <- true
	<-h.shutdown
}

func (h *Handler) OnEvent1(req interface{}) error {

	intraReply := make(chan interface{}, 1)
	intraReq := &request{
		carriedReq: onEvent1Request{
			data: req,
		},
		reply: intraReply,
	}

	r, err := h.safeEnqueueRequest(intraReq, true)
	if nil != err {
		return err
	}

	//process r
	fmt.Println("OnEvent1: r is ", r)

	return nil
}

func (h *Handler) OnEvent2(req interface{}) error {
	intraReply := make(chan interface{}, 1)
	intraReq := &request{
		carriedReq: onEvent2Request{
			data: req,
		},
		reply: intraReply,
	}

	r, err := h.safeEnqueueRequest(intraReq, true)
	if nil != err {
		return err
	}

	//process r
	fmt.Println("OnEvent2: r is ", r)

	return nil
}

func (h *Handler) OnEvent3(req interface{}) error {
	intraReply := make(chan interface{}, 1)
	intraReq := &request{
		carriedReq: onEvent3Request{
			data: req,
		},
		reply: intraReply,
	}

	_, err := h.safeEnqueueRequest(intraReq, false)
	if nil != err {
		return err
	}

	//done
	fmt.Println("OnEvent3 finished")

	return nil
}

func (h *Handler) safeEnqueueRequest(req *request, wait4reply bool) (interface{}, error) {
	select {
	case <-h.stopped:
		return nil, errors.New("Service is down")
	case h.requests <- req:
	}

	if wait4reply {
		select {
		case <-h.stopped:
			return nil, errors.New("Service is down")
		case r, ok := <-req.reply:
			if ok {
				return r, nil
			}
			return nil, errors.New("Service is down")
		}
	} else {
		return nil, nil
	}
}

func (h *Handler) processEvent1(req *onEvent1Request) error {
	return nil
}

func (h *Handler) processEvent2(req *onEvent2Request) error {
	return nil
}

func (h *Handler) processEvent3(req *onEvent3Request) error {
	return nil
}

func setReply(req *request, r interface{}, err error) {
	if err != nil {
		req.reply <- err
	} else {
		req.reply <- r
	}
}

func (h *Handler) dispatch(req *request) {
	var empty struct{}
	switch v := req.carriedReq.(type) {
	case onEvent1Request:
		err := h.processEvent1(&v)
		setReply(req, empty, err)
	case onEvent2Request:
		err := h.processEvent2(&v)
		setReply(req, empty, err)
	case onEvent3Request:
		err := h.processEvent3(&v)
		setReply(req, empty, err)
	default:
		req.reply <- errors.New("Invalid request")
	}
}

func (h *Handler) run() {
	for {
		select {
		case req, ok := <-h.requests:
			if ok {
				h.dispatch(req)
			}
		case <-h.quit:
			goto OnCloseLabel
		}
	}
OnCloseLabel:
	close(h.stopped)
	close(h.shutdown)
}
