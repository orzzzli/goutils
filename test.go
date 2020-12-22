package main

import "fmt"

type TestI interface {
	Serve()
}

type Test struct {
	name string
}

func (*Test) Serve() {
	fmt.Println("test serve")
}

type TestA struct {
}

func (*TestA) Serve() {
	fmt.Println("testa serve")
}

func NewTest(sType int, name string) TestI {
	if sType == 1 {
		return &TestA{}
	}

	return &Test{
		name: name,
	}
}

func main() {
	s := NewTest(0, "aaa")
	s.Serve()
}
