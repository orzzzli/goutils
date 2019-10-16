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
	err = Set("b", "5", 0)
	fmt.Println(err)
}

func TestSetNx(t *testing.T) {
	NewRedisPool(RedisURL, RedisPass, MaxIdleNumber, MaxIdleTime)
	res,err := SetNx("a", "b")
	fmt.Println(res,err)
	res,err = SetNx("c", "b")
	fmt.Println(res,err)
	res,err = SetNx("d", "5")
	fmt.Println(res,err)
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
	err := ZAdd("zset", "a", 5)
	fmt.Println(err)
	err = ZAdd("zset", "b", 5)
	fmt.Println(err)
	err = ZAdd("zset", "c", 50)
	fmt.Println(err)
	err = ZAdd("zset", "d", 50)
	fmt.Println(err)
}

func TestZRevRank(t *testing.T) {
	NewRedisPool(RedisURL, RedisPass, MaxIdleNumber, MaxIdleTime)
	res, find, err := ZRevRank("zset", "b")
	fmt.Println(res, reflect.TypeOf(res), find, err)
	res, find, err = ZRevRank("zset", "c")
	fmt.Println(res, reflect.TypeOf(res), find, err)
	res, find, err = ZRevRank("zset", "d")
	fmt.Println(res, reflect.TypeOf(res), find, err)
}

func TestZRem(t *testing.T) {
	NewRedisPool(RedisURL, RedisPass, MaxIdleNumber, MaxIdleTime)
	err := ZRem("zset","d")
	fmt.Println(err)
}

func TestZRange(t *testing.T) {
	NewRedisPool(RedisURL, RedisPass, MaxIdleNumber, MaxIdleTime)
	res, err := ZRange("zset",0,-1,false,true)
	fmt.Println(res,err)
	res, err = ZRange("zset",0,-1,true,false)
	fmt.Println(res,err)
}

func TestLPush(t *testing.T) {
	NewRedisPool(RedisURL, RedisPass, MaxIdleNumber, MaxIdleTime)
	err := LPush("list","a")
	fmt.Println(err)
	err = LPush("list","b")
	fmt.Println(err)
	err = LPush("list","c")
	fmt.Println(err)
}

func TestRPop(t *testing.T) {
	NewRedisPool(RedisURL, RedisPass, MaxIdleNumber, MaxIdleTime)
	res,err := RPop("list")
	fmt.Println(res,err)
	res,err = RPop("list")
	fmt.Println(res,err)
}

func TestLRange(t *testing.T) {
	NewRedisPool(RedisURL, RedisPass, MaxIdleNumber, MaxIdleTime)
	res, err := LRange("list",0,-1)
	fmt.Println(res,err)
}

func TestExpire(t *testing.T) {
	NewRedisPool(RedisURL, RedisPass, MaxIdleNumber, MaxIdleTime)
	err := Expire("a", 60)
	fmt.Println(err)
	err = Expire("b", 0)
	fmt.Println(err)
	err = Expire("c", 0)
	fmt.Println(err)
}

func TestDel(t *testing.T) {
	NewRedisPool(RedisURL, RedisPass, MaxIdleNumber, MaxIdleTime)
	err := Del("d")
	fmt.Println(err)
}
