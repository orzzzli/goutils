package server

import (
	"google.golang.org/grpc"
)

const (
	ServeTypeHttp = 0
	ServeTypeGRPC = 1
)

type ServerInterface interface {
	Serve() error
	RegisterService(url string, handler ...interface{}) error
}

type Server struct {
	name       string
	address    string
	serviceMap map[string]interface{}
	gServer    *grpc.Server
}

//工厂
func NewServer(sType int, name string, address string) ServerInterface {
	if sType == ServeTypeHttp {
		return newHttpServer(name, address)
	}
	if sType == ServeTypeGRPC {
		s := newGrpcServer(name, address)
		s.gServer = grpc.NewServer()
		return s
	}
	return nil
}
