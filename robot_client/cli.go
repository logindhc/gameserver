package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

func client(userName string) {
	ctx, cancel := context.WithCancel(context.Background())
	// 监听关闭信号
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		cancel() // 取消上下文
	}()
	robot := RunRobot(url, userName, true)
	go scanner(robot)
	// 等待上下文被取消
	<-ctx.Done()
}

func scanner(cli *Robot) {
	// 从标准输入流中接收输入数据
	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		line := input.Text()
		split := strings.Split(line, " ")
		if split[0] == "help" {
			fmt.Println("item use")
			continue
		} else if split[0] == "item" {
			err := cli.GetItemInfo()
			if err != nil {
				continue
			}
			continue
		} else if split[0] == "use" {
			err := cli.UseItem()
			if err != nil {
				continue
			}
			continue
		}
	}
}
