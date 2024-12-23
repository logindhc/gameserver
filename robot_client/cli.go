package main

import (
	"bufio"
	"context"
	"fmt"
	cstring "gameserver/cherry/extend/string"
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
			if len(split) < 3 {
				continue
			}
			itemId := split[1]
			count := split[2]
			err := cli.UseItem(cstring.ToInt32D(itemId), cstring.ToInt32D(count))
			if err != nil {
				cli.Debug(err)
				continue
			}
			continue
		} else if split[0] == "gm" {
			if len(split) < 3 {
				continue
			}
			cmd := split[1]
			args := split[2]
			err := cli.Gm(cmd, args)
			if err != nil {
				cli.Debug(err)
				continue
			}
			continue
		} else if split[0] == "heroUp" {
			if len(split) < 2 {
				continue
			}
			heroId := split[1]
			err := cli.HeroUp(cstring.ToInt32D(heroId))
			if err != nil {
				cli.Debug(err)
				continue
			}
			continue
		}
	}
}
