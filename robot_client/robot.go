package main

import (
	"fmt"
	"gameserver/cherry/net/parser/pomelo/message"
	"gameserver/internal/code"
	"gameserver/internal/pb"
	"math/rand"
	"time"

	cherryError "gameserver/cherry/error"
	cherryHttp "gameserver/cherry/extend/http"
	cherryTime "gameserver/cherry/extend/time"
	cherryLogger "gameserver/cherry/logger"
	cherryClient "gameserver/cherry/net/parser/pomelo/client"
	jsoniter "github.com/json-iterator/go"
)

type (
	// Robot client robot
	Robot struct {
		*cherryClient.Client
		PrintLog   bool
		Token      string
		PID        int32
		CID        int32
		OpenId     string
		PlayerId   int64
		PlayerName string
		StartTime  cherryTime.CherryTime
		address    string
	}
)

func New(client *cherryClient.Client) *Robot {
	return &Robot{
		Client: client,
	}
}

// GetserverInfo  http登录获取token对象
// http://172.16.124.137/serverInfo?pid=2126003&account=test1&password=test1
func (p *Robot) GetServerInfo(url, pCode, channel, platform string) error {
	// http登陆获取token json对象
	requestURL := fmt.Sprintf("%s/api/serverInfo", url)
	jsonBytes, _, err := cherryHttp.GET(requestURL, map[string]string{
		"code":     pCode, //帐号名
		"channel":  channel,
		"platform": platform, //平台id
		"version":  "1.0",
	})

	if err != nil {
		return err
	}
	// 转换json对象
	rsp := code.Result{}
	if err = jsoniter.Unmarshal(jsonBytes, &rsp); err != nil {
		return err
	}

	if code.IsFail(rsp.Code) {
		return cherryError.Errorf("GetServerInfo fail. [code = %v]", rsp.Code)
	}
	maps := rsp.Data.(map[string]interface{})
	p.address = fmt.Sprintf("%v:%v", maps["server_host"].(string), maps["server_port"].(string))
	p.Token = maps["token"].(string)
	p.TagName = fmt.Sprintf("%s_%s", channel, pCode)
	p.StartTime = cherryTime.Now()
	p.Debugf("GetServerInfo success. %v", rsp)
	return nil
}

// UserLogin 用户登录对某游戏服
func (p *Robot) UserLogin() error {
	route := "gate.user.login"

	p.Debugf("[%s] [UserLogin] request route = %s", p.TagName, route)

	msg, err := p.Request(route, &pb.C2SLogin{
		Token:  p.Token,
		Params: nil,
	})

	if err != nil {
		return err
	}

	rsp := new(pb.S2CLogin)
	err = p.Serializer().Unmarshal(msg.Data, rsp)
	if err != nil {
		return err
	}
	p.PlayerId = rsp.Uid
	p.Debugf("[%s] [UserLogin] response = %+v", p.TagName, rsp)
	return nil
}

// ActorEnter 角色进入游戏
func (p *Robot) ActorEnter() error {
	route := "game.player.enter"
	req := &pb.C2SPlayerEnter{}

	msg, err := p.Request(route, req)
	if err != nil {
		return err
	}

	rsp := &pb.S2CPlayerEnter{}
	err = p.Serializer().Unmarshal(msg.Data, rsp)
	if err != nil {
		return err
	}
	player := rsp.GetPlayer()

	p.PlayerName = player.PlayerName

	p.Debugf("[%s] [PlayerEnter] response %v", p.TagName, rsp)
	return nil
}

// GetItemInfo 角色获取道具信息
func (p *Robot) GetItemInfo() error {
	route := "game.player.itemInfo"
	req := &pb.None{}

	msg, err := p.Request(route, req)
	if err != nil {
		return err
	}
	rsp := &pb.S2CItemInfo{}
	err = p.Serializer().Unmarshal(msg.Data, rsp)
	if err != nil {
		return err
	}
	p.Debugf("[%s] [getItemInfo] response ret = %v", p.TagName, rsp)
	return nil
}

func (p *Robot) UseItem(itemId, count int32) error {
	route := "game.player.itemUse"
	req := &pb.C2SItemUse{
		ItemId: itemId,
		Count:  count,
	}

	msg, err := p.Request(route, req)
	if err != nil {
		return err
	}
	rsp := &pb.S2CItemUse{}
	err = p.Serializer().Unmarshal(msg.Data, rsp)
	if err != nil {
		return err
	}
	p.Debugf("[%s] [%v] response ret = %v", p.TagName, route, rsp)
	return nil
}

func (p *Robot) HeroUp(heroId int32) error {
	route := "game.player.heroUp"
	req := &pb.C2SHeroUp{
		HeroId: heroId,
	}

	msg, err := p.Request(route, req)
	if err != nil {
		return err
	}
	rsp := &pb.S2CHeroUp{}
	err = p.Serializer().Unmarshal(msg.Data, rsp)
	if err != nil {
		return err
	}
	p.Debugf("[%s] [%v] response ret = %v", p.TagName, route, rsp)
	return nil
}

func (p *Robot) Gm(cmd, args string) error {
	route := "game.player.gm"
	req := &pb.C2SPlayerGM{
		Cmd:  cmd,
		Args: args,
	}
	msg, err := p.Request(route, req)
	if err != nil {
		return err
	}
	rsp := &pb.S2CPlayerGM{}
	err = p.Serializer().Unmarshal(msg.Data, rsp)
	if err != nil {
		return err
	}
	p.Debugf("[%s] [%v] response ret = %v", p.TagName, route, rsp)
	return nil
}

func (p *Robot) RandSleep() {
	time.Sleep(time.Duration(rand.Int31n(10)) * time.Millisecond)
}

func (p *Robot) Debug(args ...interface{}) {
	if p.PrintLog {
		cherryLogger.Debug(args...)
	}

}

func (p *Robot) Debugf(template string, args ...interface{}) {
	if p.PrintLog {
		cherryLogger.Debugf(template, args...)
	}
}

func (p *Robot) ResUpdate(msg *pomeloMessage.Message) {
	rsp := &pb.S2CResUpdate{}
	err := p.Serializer().Unmarshal(msg.Data, rsp)
	if err != nil {
		p.Debugf("ResUpdate fail. [err = %v]", err)
		return
	}
	p.Debugf("[%s] [%v] push ret = %v", p.TagName, msg.Route, rsp)
}

func (p *Robot) CurrencyInfo(msg *pomeloMessage.Message) {
	rsp := &pb.S2CCurrencyInfo{}
	err := p.Serializer().Unmarshal(msg.Data, rsp)
	if err != nil {
		p.Debugf("CurrencyInfo fail. [err = %v]", err)
		return
	}
	p.Debugf("[%s] [%v] push ret = %v", p.TagName, msg.Route, rsp)
}

func (p *Robot) HeroInfo(msg *pomeloMessage.Message) {
	rsp := &pb.S2CHeroInfo{}
	err := p.Serializer().Unmarshal(msg.Data, rsp)
	if err != nil {
		p.Debugf("HeroInfo fail. [err = %v]", err)
		return
	}
	p.Debugf("[%s] [%v] push ret = %v", p.TagName, msg.Route, rsp)
}

func (p *Robot) ItemInfo(msg *pomeloMessage.Message) {
	rsp := &pb.S2CItemInfo{}
	err := p.Serializer().Unmarshal(msg.Data, rsp)
	if err != nil {
		p.Debugf("ItemInfo fail. [err = %v]", err)
		return
	}
	p.Debugf("[%s] [%v] push ret = %v", p.TagName, msg.Route, rsp)
}
