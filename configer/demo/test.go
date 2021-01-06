package main

import (
	"fmt"
	"time"

	"github.com/orzzzli/goutils/configer"
)

func main() {
	cfg, err := configer.NewiniConfiger("./test.ini")
	if err != nil {
		fmt.Println(err)
		return
	}
	cfg.SetScanSec(1)
	err = cfg.Invoke()
	if err != nil {
		fmt.Println(err)
		return
	}
	go func() {
		err := cfg.GetHotLoadingErr()
		fmt.Println(err)
	}()
	fmt.Println(cfg.GetSection("mysql"))
	time.Sleep(10 * time.Second)
	fmt.Println(cfg.GetSection("mysql"))
	time.Sleep(10 * time.Second)
	fmt.Println(cfg.GetSection("mysql"))
}
