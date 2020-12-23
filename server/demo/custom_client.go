package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8888")
	if err != nil {
		fmt.Println(err)
		return
	}

	go func() {
		buf := make([]byte, 1000)
		for {
			l, err := conn.Read(buf)
			if err == io.EOF {
				break
			}
			fmt.Println("server send:" + string(buf[:l]))
		}
	}()

	inputReader := bufio.NewReader(os.Stdin)
	fmt.Println("URL?")
	input, _ := inputReader.ReadString('\n')
	trimmedInput := strings.Trim(input, "\r\n")
	// 给服务器发送信息直到程序退出：
	for {
		fmt.Println("message?")
		input2, err := inputReader.ReadString('\n')
		if err == io.EOF {
			break
		}
		trimmedInput2 := strings.Trim(input2, "\r\n")
		_, err = conn.Write([]byte(strconv.Itoa(len(trimmedInput2)) + trimmedInput + trimmedInput2))
	}
}
