package goUtils

import (
	"errors"
	"github.com/gomodule/redigo/redis"
	"time"
)

var RedisPool *redis.Pool

// NewRedisPool初始化连接池
func NewRedisPool(url string, password string, idle int, idleTime int) {
	RedisPool = &redis.Pool{
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
	if RedisPool == nil {
		return errors.New("redis pool is not init.")
	}
	conn := RedisPool.Get()
	defer conn.Close()
	var err error
	_,err = conn.Do("EXPIRE",k,ex)
	return err
}
//string
func Set(k string, v string, ex int) error {
	if RedisPool == nil {
		return errors.New("redis pool is not init.")
	}
	conn := RedisPool.Get()
	defer conn.Close()
	var err error
	if ex <= 0 {
		_,err = conn.Do("SET",k,v)
	}else{
		_,err = conn.Do("SET",k,v,"EX",ex)
	}
	return err
}
func Get(k string) (string,error) {
	if RedisPool == nil {
		return "",errors.New("redis pool is not init.")
	}
	conn := RedisPool.Get()
	defer conn.Close()
	res,err := conn.Do("GET",k)
	return res.(string),err
}
//SortedSet
func ZAdd(key string,k string, v string) error {
	if RedisPool == nil {
		return errors.New("redis pool is not init.")
	}
	conn := RedisPool.Get()
	defer conn.Close()
	var err error
	_,err = conn.Do("ZADD",key,k,v)
	return err
}
func ZRevRank(key string,k string) (int64,error) {
	if RedisPool == nil {
		return 0,errors.New("redis pool is not init.")
	}
	conn := RedisPool.Get()
	defer conn.Close()
	result,err := conn.Do("ZREVRANK",key,k)
	return result.(int64),err
}