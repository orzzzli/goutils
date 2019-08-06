package convert

import (
	"fmt"
	"reflect"
	"testing"
)

func TestInt64to32(t *testing.T) {
	res, err := Int64to32(111111)
	fmt.Println(reflect.TypeOf(res), err)
}

func TestInt32to64(t *testing.T) {
	res, err := Int32to64(666)
	fmt.Println(reflect.TypeOf(res), err)
}
