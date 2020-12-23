package main

import (
	"bufio"
	"fmt"
	"io"
	"strconv"

	server2 "github.com/orzzzli/goutils/server"
)

var cServer *server2.CustomServer

func main() {
	s := server2.NewServer(server2.ServeTypeCustom, "test", "localhost:8888")
	cServer, _ = s.(*server2.CustomServer)
	//第一个字节代表长度，为了方便直接使用字符串，也就是最大9个message长度
	//第二个字节代表地址
	//注册协议解析方法
	cServer.SetHandler(testRead)
	//注册功能
	err := cServer.RegisterService("a", echoFunc)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = cServer.Serve()
	if err != nil {
		fmt.Println(err)
		return
	}
}

func testRead(reader *bufio.Reader, connIndex int) {
	for {
		//消息长度
		tmp, err := reader.ReadByte()
		if err == io.EOF {
			break
		}
		msgLen, _ := strconv.Atoi(string(tmp))
		//消息路由
		tmp, err = reader.ReadByte()
		if err == io.EOF {
			break
		}
		url := string(tmp)

		msg := ""
		for i := 0; i < msgLen; i++ {
			tmp1, err := reader.ReadByte()
			if err == io.EOF {
				break
			}
			msg += string(tmp1)
		}
		//同步调用方法
		res, err := cServer.CallService(url, msg)
		if err != nil {
			fmt.Println(url, msg)
			fmt.Println("call service error" + err.Error())
			continue
		}
		//定点发给自己
		cServer.SendTo(connIndex, res)
		//广播通知所有人
		cServer.BroadcastTo(res)
	}
}

//模拟一个简单回射功能
func echoFunc(p interface{}) string {
	res, _ := p.(string)
	return res
}
