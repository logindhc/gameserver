package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	ctime "gameserver/cherry/extend/time"
	clog "gameserver/cherry/logger"
	pomeloClient "gameserver/cherry/net/parser/pomelo/client"
)

var (
	maxRobotNum = 5000                   // 运行x个机器人
	url         = "http://0.0.0.0:10000" // web node
	channel     = "101"                  // 测试的渠道
	platform    = "3"                    // 测试的平台
	printLog    = false                  // 是否输出详细日志
)

func main() {
	//if len(os.Args) > 0 {
	//	client(os.Args[0])
	//	return
	//}

	client("dhc33")

	//runRobot()
}

func runRobot() {
	wg := sync.WaitGroup{}
	wg.Add(1)

	accounts := make(map[string]string)
	for i := 1; i <= maxRobotNum; i++ {
		key := fmt.Sprintf("test_%d", i)
		accounts[key] = key
	}

	for userName, _ := range accounts {
		time.Sleep(time.Duration(rand.Int31n(20)) * time.Millisecond)
		go RunRobot(url, userName, printLog)
	}

	wg.Wait()
}

func RunRobot(url, userName string, printLog bool) *Robot {

	// 创建客户端
	cli := New(
		pomeloClient.New(
			pomeloClient.WithRequestTimeout(10*time.Second),
			pomeloClient.WithErrorBreak(true),
		),
	)
	cli.On("currencyInfo", cli.CurrencyInfo)
	cli.On("heroInfo", cli.HeroInfo)
	cli.On("itemInfo", cli.ItemInfo)
	cli.On("resUpdate", cli.ResUpdate)
	cli.PrintLog = printLog

	// 登录获取token
	if err := cli.GetServerInfo(url, userName, channel, platform); err != nil {
		clog.Error(err)
		return nil
	}

	//split := strings.Split(cli.address, ":")
	//port, _ := strconv.Atoi(split[1])
	//address := fmt.Sprintf("%s:%d", split[0], port+1)
	//// 根据地址连接网关
	//if err := cli.ConnectToTCP(address); err != nil {
	//	clog.Error(err)
	//	return nil
	//}
	address := cli.address
	// 根据地址连接网关
	if err := cli.ConnectToWS(address, ""); err != nil {
		clog.Error(err)
		return nil
	}

	if cli.PrintLog {
		clog.Infof("connect %s is ok", address)
	}

	// 随机休眠
	cli.RandSleep()

	// 用户登录到网关节点
	err := cli.UserLogin()
	if err != nil {
		clog.Warn(err)
		return nil
	}

	if cli.PrintLog {
		clog.Infof("user login is ok. [user = %s]", userName)
	}

	//cli.RandSleep()

	// 角色进入游戏
	err = cli.ActorEnter()
	if err != nil {
		clog.Warn(err)
		return nil
	}

	elapsedTime := cli.StartTime.DiffInMillisecond(ctime.Now())
	clog.Debugf("[%s] is enter to game. elapsed time:%dms", cli.TagName, elapsedTime)

	err = cli.GetItemInfo()
	if err != nil {

		return nil
	}
	//cli.Disconnect()

	return cli
}
