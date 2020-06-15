package main

import (
	"fmt"
	"os"
	"os/signal"
	"socket_zero_copy"
	"syscall"
	"time"
)

func main() {
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	server := socket_zero_copy.NewSoketServer(10000)
	go server.Start() //启动服务
	fmt.Println("start success")

	for s := range c {
		switch s {
		case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
			fmt.Println("退出", s)
			server.Stop()
			time.Sleep(time.Second * 3) //等待三秒
			os.Exit(0)
		default:
			fmt.Println("other", s)
		}
	}

}
