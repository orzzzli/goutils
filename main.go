package main

import (
	"fmt"
	"goutils/redis"
)

func main() {
	//init redis pool
	redis.NewRedisPool("127.0.0.1","",4,300)

	err := redis.Set("a","b",0)
	if err != nil {
		fmt.Println(err)
	}
	res, err := redis.Get("a")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res)
}