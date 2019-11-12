package convert

import (
	"fmt"
	"reflect"
	"testing"
)

func TestInt64to32(t *testing.T) {
	res, err := Int64to32(111111)
	fmt.Println(reflect.TypeOf(res), err)
	res, err = Int64to32(-111111)
	fmt.Println(reflect.TypeOf(res), err)
}

func TestInt32to64(t *testing.T) {
	res, err := Int32to64(666)
	fmt.Println(reflect.TypeOf(res), err)
	res, err = Int32to64(-666)
	fmt.Println(reflect.TypeOf(res), err)
}

func TestInt32to64bitStr(t *testing.T) {
	res := Int32to64bitStr(10000000)
	fmt.Println(res)
	res = Int32to64bitStr(60)
	fmt.Println(res)
	res = Int32to64bitStr(-10000000)
	fmt.Println(res)
	res = Int32to64bitStr(0)
	fmt.Println(res)
}

func TestInterfaceToStr(t *testing.T) {
	fmt.Println(InterfaceToStr("aaaa"))
	fmt.Println(InterfaceToStr(11111))
	var inta int64
	inta = 111111111111111111
	fmt.Println(InterfaceToStr(inta))
	var temp []string
	temp = append(temp,"a")
	temp = append(temp,"b")
	temp = append(temp,"c")
	fmt.Println(InterfaceToStr(temp))
	type Temp struct {
		aa string
		bb string
		cc int
	}
	fmt.Println(InterfaceToStr(&Temp{
		aa:"aaa",
		bb:"bbb",
		cc:1,
	}))
}
