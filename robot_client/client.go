package main

import (
	"bufio"
	"context"
	"fmt"
	pomeloClient "gameserver/cherry/net/parser/pomelo/client"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

var (
	username, pwd string
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	// 监听关闭信号
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		cancel() // 取消上下文
	}()

	cli := New(pomeloClient.New(
		pomeloClient.WithRequestTimeout(10*time.Second),
		pomeloClient.WithErrorBreak(true),
	))

	cli.PrintLog = true
	err := cli.ConnectToWS("0.0.0.0:10010", "")
	defer cli.Disconnect()
	if err != nil {
		return
	}

	go scanner(cli)

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
			fmt.Println("register|token login select|create enter item")
			continue
		}
		if split[0] == "register" {
			username = split[1]
			pwd = split[2]
			maps := map[string]string{
				split[1]: split[2],
			}
			RegisterDevAccount(url, maps)
			continue
		} else if split[0] == "token" {
			if len(split) >= 2 {
				username = split[1]
			}
			if len(split) >= 3 {
				pwd = split[2]
			}
			err := cli.GetToken(url, pid, username, pwd)
			if err != nil {
				continue
			}
			continue
		} else if split[0] == "login" {
			err := cli.UserLogin(serverId)
			if err != nil {
				continue
			}
			continue
		} else if split[0] == "select" {
			err := cli.PlayerSelect()
			if err != nil {
				continue
			}
			continue
		} else if split[0] == "create" {
			err := cli.ActorCreate()
			if err != nil {
				continue
			}
			continue
		} else if split[0] == "enter" {
			err := cli.ActorEnter()
			if err != nil {
				continue
			}
			continue
		} else if split[0] == "item" {
			err := cli.GetItemInfo()
			if err != nil {
				continue
			}
			continue
		}
	}
}
