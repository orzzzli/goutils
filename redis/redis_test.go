package redis

import (
	"fmt"
	"reflect"
	"testing"
)

const (
	RedisURL      = "redis://127.0.0.1:6379"
	RedisPass     = ""
	MaxIdleNumber = 1
	MaxIdleTime   = 60
)

func TestSet(t *testing.T) {
	NewRedisPool(RedisURL, RedisPass, MaxIdleNumber, MaxIdleTime)
	err := Set("a", "b", 60)
	fmt.Println(err)
	err = Set("a", "b", 0)
	fmt.Println(err)
	err = Set("c", "5", 0)
	fmt.Println(err)
}

func TestGet(t *testing.T) {
	NewRedisPool(RedisURL, RedisPass, MaxIdleNumber, MaxIdleTime)
	res, find, err := Get("a")
	fmt.Println(res, reflect.TypeOf(res), find, err)
	res, find, err = Get("b")
	fmt.Println(res, reflect.TypeOf(res), find, err)
	res, find, err = Get("c")
	fmt.Println(res, reflect.TypeOf(res), find, err)
}

func TestZAdd(t *testing.T) {
	NewRedisPool(RedisURL, RedisPass, MaxIdleNumber, MaxIdleTime)
	err := ZAdd("d", "b", 5)
	fmt.Println(err)
	err = ZAdd("d", "b", 5)
	fmt.Println(err)
	err = ZAdd("d", "bb", 50)
	fmt.Println(err)
}

func TestZRevRank(t *testing.T) {
	NewRedisPool(RedisURL, RedisPass, MaxIdleNumber, MaxIdleTime)
	res, find, err := ZRevRank("d", "b")
	fmt.Println(res, reflect.TypeOf(res), find, err)
	res, find, err = ZRevRank("d", "bb")
	fmt.Println(res, reflect.TypeOf(res), find, err)
	res, find, err = ZRevRank("d", "cc")
	fmt.Println(res, reflect.TypeOf(res), find, err)
}

func TestExpire(t *testing.T) {
	NewRedisPool(RedisURL, RedisPass, MaxIdleNumber, MaxIdleTime)
	err := Expire("a", 60)
	fmt.Println(err)
	err = Expire("a", 0)
	fmt.Println(err)
	err = Expire("c", 0)
	fmt.Println(err)
}
