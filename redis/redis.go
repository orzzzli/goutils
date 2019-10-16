package redis

import (
	"errors"
	"github.com/gomodule/redigo/redis"
	"github.com/orzzzli/goutils/convert"
	"time"
)

var GlobalRedisPool *redis.Pool

// NewRedisPool初始化连接池
func NewRedisPool(url string, password string, idle int, idleTime int) {
	GlobalRedisPool = &redis.Pool{
		MaxIdle:     idle,
		IdleTimeout: time.Duration(idleTime) * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.DialURL(url)
			if err != nil {
				panic(err)
			}
			if password != "" {
				//验证redis密码
				if _, authErr := c.Do("AUTH", password); authErr != nil {
					panic(err)
				}
			}
			return c, nil
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			if err != nil {
				panic(err)
			}
			return nil
		},
	}
}

//expire
func Expire(k string, ex int) error {
	if GlobalRedisPool == nil {
		return errors.New("redis pool is not init.")
	}
	conn := GlobalRedisPool.Get()
	defer conn.Close()
	var err error
	_, err = conn.Do("EXPIRE", k, ex)
	return err
}

//string
func Set(k string, v string, ex int) error {
	if GlobalRedisPool == nil {
		return errors.New("redis pool is not init.")
	}
	conn := GlobalRedisPool.Get()
	defer conn.Close()
	var err error
	if ex <= 0 {
		_, err = conn.Do("SET", k, v)
	} else {
		_, err = conn.Do("SET", k, v, "EX", ex)
	}
	return err
}
func SetNx(k string, v string) (bool,error) {
	if GlobalRedisPool == nil {
		return false,errors.New("redis pool is not init.")
	}
	conn := GlobalRedisPool.Get()
	defer conn.Close()
	var err error
	var res interface{}
	res, err = conn.Do("SETNX", k, v)
	if err != nil {
		return false, err
	}
	if res.(int64) == 1 {
		return true,nil
	}else{
		return false,nil
	}
}
func Get(k string) (string, bool, error) {
	if GlobalRedisPool == nil {
		return "", false, errors.New("redis pool is not init.")
	}
	conn := GlobalRedisPool.Get()
	defer conn.Close()
	res, err := conn.Do("GET", k)
	resOp := ""
	find := false
	if res != nil {
		resOp = string(res.([]uint8))
		find = true
	}
	return resOp, find, err
}

//sortedSet
func ZAdd(key string, k string, v float32) error {
	if GlobalRedisPool == nil {
		return errors.New("redis pool is not init.")
	}
	conn := GlobalRedisPool.Get()
	defer conn.Close()
	var err error
	_, err = conn.Do("ZADD", key, v, k)
	return err
}
func ZRem(key string,k string) error {
	if GlobalRedisPool == nil {
		return errors.New("redis pool is not init.")
	}
	conn := GlobalRedisPool.Get()
	defer conn.Close()
	var err error
	_, err = conn.Do("ZREM", key, k)
	return err
}
func ZRevRank(key string, k string) (int, bool, error) {
	if GlobalRedisPool == nil {
		return 0, false, errors.New("redis pool is not init.")
	}
	conn := GlobalRedisPool.Get()
	defer conn.Close()
	resultOp := 0
	find := false
	result, err := conn.Do("ZREVRANK", key, k)
	if result != nil {
		resultOp, _ = convert.Int64to32(result.(int64))
		find = true
	}
	return resultOp, find, err
}
func ZRange(key string, start int, end int, desc bool, withScore bool) ([]string, error) {
	if GlobalRedisPool == nil {
		return nil,errors.New("redis pool is not init.")
	}
	conn := GlobalRedisPool.Get()
	defer conn.Close()
	command := ""
	if desc {
		command = "ZREVRANGE"
	}else{
		command = "ZRANGE"
	}
	var res interface{}
	var err error
	if withScore {
		res, err = conn.Do(command, key, start, end,"WITHSCORES")
	}else{
		res, err = conn.Do(command, key, start, end)
	}
	var resOp []string
	for _,v := range res.([]interface{}) {
		resOp = append(resOp,string(v.([]uint8)))
	}
	return resOp, err
}

//list
func LPush(k string, v string) error {
	if GlobalRedisPool == nil {
		return errors.New("redis pool is not init.")
	}
	conn := GlobalRedisPool.Get()
	defer conn.Close()
	var err error
	_, err = conn.Do("LPUSH", k, v)
	return err
}
func RPop(k string) (string, error) {
	if GlobalRedisPool == nil {
		return "", errors.New("redis pool is not init.")
	}
	conn := GlobalRedisPool.Get()
	defer conn.Close()
	res, err := conn.Do("RPOP", k)
	resOp := ""
	if res != nil {
		resOp = string(res.([]uint8))
	}
	return resOp, err
}
func LRange(k string, start int, end int) ([]string, error) {
	if GlobalRedisPool == nil {
		return nil, errors.New("redis pool is not init.")
	}
	conn := GlobalRedisPool.Get()
	defer conn.Close()
	res, err := conn.Do("LRANGE", k, start, end)
	var resOp []string
	for _,v := range res.([]interface{}) {
		resOp = append(resOp,string(v.([]uint8)))
	}
	return resOp, err
}

//del
func Del(k string) error {
	if GlobalRedisPool == nil {
		return errors.New("redis pool is not init.")
	}
	conn := GlobalRedisPool.Get()
	defer conn.Close()
	var err error
	_, err = conn.Do("DEL", k)
	return err
}
