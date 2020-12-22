package server

import (
	"errors"
	"fmt"
	"net"
	"net/http"
)

type httpServer struct {
	Server
}

func newHttpServer(name string, address string) *httpServer {
	t := &httpServer{}
	t.name = name
	t.address = address
	t.serviceMap = make(map[string]interface{}, 5)
	return t
}

func (h *httpServer) Serve() error {
	l, err := net.Listen("tcp", h.address)
	if err != nil {
		return err
	}
	err = http.Serve(l, h)
	if err != nil {
		return err
	}
	return nil
}

func (h *httpServer) RegisterService(url string, handler ...interface{}) error {
	_, exist := h.serviceMap[url]
	if exist {
		return errors.New("service already register")
	}
	if len(handler) == 0 {
		return errors.New("handler must send")
	}
	t, ok := handler[0].(func(r *http.Request) string)
	if !ok {
		return errors.New("handler type must be [ func(request *http.Request) string ]")
	}
	h.serviceMap[url] = t
	return nil
}

func (h *httpServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	service, ok := h.serviceMap[r.URL.Path]
	if !ok {
		fmt.Fprintln(w, "not found")
		return
	}
	handler, _ := service.(func(r *http.Request) string)
	res := handler(r)
	fmt.Fprintln(w, res)
}
