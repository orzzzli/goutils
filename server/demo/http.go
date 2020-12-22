package main

import (
	"fmt"
	"net/http"

	server2 "github.com/orzzzli/goutils/server"
)

func main() {
	s := server2.NewServer(server2.ServeTypeHttp, "test", "localhost:8888")
	err := s.RegisterService("/", func(request *http.Request) string {
		return "helloWorld"
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	err = s.RegisterService("/hello", func(request *http.Request) string {
		return "hello"
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	err = s.Serve()
	if err != nil {
		fmt.Println(err)
	}
}
