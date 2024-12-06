package token

import (
	"encoding/json"
	"fmt"
	"gameserver/internal/code"

	cherryCrypto "gameserver/cherry/extend/crypto"
	cherryTime "gameserver/cherry/extend/time"
	cherryLogger "gameserver/cherry/logger"
)

const (
	hashFormat      = "uid:%d,open_id:%s,channel:%d,platform:%d,server_id:%d,timestamp:%d"
	tokenExpiredDay = 3
)

type Token struct {
	UID       int64  `gorm:"column:uid;comment:玩家ID" json:"uid"`
	OpenId    string `json:"open_id"`
	Channel   int32  `json:"channel"`
	Platform  int32  `json:"platform"`
	ServerId  int32  `json:"server_id"`
	Timestamp int64  `json:"tt"`
	Hash      string `json:"hash"`
}

func New(uid int64, openId string, channel int32, platform int32, serverId int32, appKey string) *Token {
	token := &Token{
		UID:       uid,
		OpenId:    openId,
		Channel:   channel,
		Platform:  platform,
		ServerId:  serverId,
		Timestamp: cherryTime.Now().ToMillisecond(),
	}

	token.Hash = BuildHash(token, appKey)
	return token
}

func (t *Token) ToBase64() string {
	bytes, _ := json.Marshal(t)
	return cherryCrypto.Base64Encode(string(bytes))
}

func DecodeToken(base64Token string) (*Token, bool) {
	if len(base64Token) < 1 {
		return nil, false
	}

	token := &Token{}
	bytes, err := cherryCrypto.Base64DecodeBytes(base64Token)
	if err != nil {
		cherryLogger.Warnf("base64Token = %s, validate error = %v", base64Token, err)
		return nil, false
	}

	err = json.Unmarshal(bytes, token)
	if err != nil {
		cherryLogger.Warnf("base64Token = %s, unmarshal error = %v", base64Token, err)
		return nil, false
	}

	return token, true
}

func Validate(token *Token, appKey string) (int32, bool) {
	now := cherryTime.Now()
	now.AddDays(tokenExpiredDay)

	if token.Timestamp > now.ToMillisecond() {
		cherryLogger.Warnf("token is expired, token = %s", token)
		return code.AccountTokenValidateFail, false
	}

	newHash := BuildHash(token, appKey)
	if newHash != token.Hash {
		cherryLogger.Warnf("hash validate fail. newHash = %s, token = %s", token)
		return code.AccountTokenValidateFail, false
	}

	return code.OK, true
}

func BuildHash(t *Token, appKey string) string {
	value := fmt.Sprintf(hashFormat, t.UID, t.OpenId, t.Channel, t.Platform, t.ServerId, t.Timestamp)
	return cherryCrypto.MD5(value + appKey)
}
