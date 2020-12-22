package server

import (
	"context"
	"errors"
	"fmt"
	"net"
	"reflect"
	"strings"

	"google.golang.org/grpc"
)

type grpcServer struct {
	Server
}

func newGrpcServer(name string, address string) *grpcServer {
	t := &grpcServer{}
	t.name = name
	t.address = address
	t.serviceMap = make(map[string]interface{}, 1)
	return t
}

func (g *grpcServer) RegisterService(url string, handler ...interface{}) error {
	_, exist := g.serviceMap[url]
	if exist {
		return errors.New("service already register")
	}
	if len(handler) < 2 {
		return errors.New("grpc generate server and implement server must send")
	}
	g.serviceMap[url] = nil
	g.analyseService(handler[0], handler[1])
	return nil
}

func (g *grpcServer) Serve() error {
	l, err := net.Listen("tcp", g.address)
	if err != nil {
		return err
	}
	err = g.gServer.Serve(l)
	if err != nil {
		return err
	}
	return nil
}

func (g *grpcServer) analyseService(grpcServer interface{}, implementServer interface{}) {
	sType := reflect.TypeOf(grpcServer).Elem()
	serverIndex := strings.Index(sType.String(), "Server")
	if serverIndex == -1 {
		fmt.Println("missing")
		return
	}
	methodList := make([]grpc.MethodDesc, 3)
	funcNumbers := sType.NumMethod()
	for i := 0; i < funcNumbers; i++ {
		method := sType.Method(i)
		var reuqestValue reflect.Value
		//get request
		for j := 0; j < sType.Method(i).Type.NumIn(); j++ {
			if sType.Method(i).Type.In(j).Kind() == reflect.Ptr {
				reuqestValue = reflect.New(sType.Method(i).Type.In(j).Elem())
				break
			}
		}
		methodList = append(methodList, grpc.MethodDesc{
			MethodName: sType.Method(i).Name,
			Handler: func(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (i interface{}, e error) {
				in := reuqestValue.Interface()
				if err := dec(in); err != nil {
					return nil, err
				}
				params := make([]reflect.Value, 2)
				params[0] = reflect.ValueOf(ctx)
				params[1] = reflect.ValueOf(in)
				funcRes := reflect.ValueOf(srv).MethodByName(method.Name).Call(params)
				if funcRes[1].IsNil() {
					return funcRes[0].Interface(), nil
				}
				return funcRes[0].Interface(), funcRes[1].Interface().(error)
			},
		})
	}
	g.gServer.RegisterService(&grpc.ServiceDesc{
		ServiceName: sType.String()[:serverIndex],
		HandlerType: grpcServer,
		Methods:     methodList,
		Streams:     []grpc.StreamDesc{},
		Metadata:    sType.PkgPath() + ".proto",
	}, implementServer)

}
