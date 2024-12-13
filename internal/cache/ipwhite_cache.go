package cache

import (
	"fmt"
	cherryredis "gameserver/cherry/components/redis"
	jsoniter "github.com/json-iterator/go"
)

var (
	RedisIpWhite = "admin:ipWhite:"
)

type IpWhite struct {
	Ip     int32  `json:"ip"`     //ip地址
	Remark string `json:"remark"` // 备注
	Enable int    `json:"enable"` // 是否启用
}

func GetIpWhite(ip int32) (*IpWhite, error) {
	key := fmt.Sprintf("%s%v", RedisIpWhite, ip)
	res, err := cherryredis.GetRds().Get(ctx, key).Result()
	ipWhite := new(IpWhite)
	if err != nil {
		return nil, err
	}
	err = jsoniter.Unmarshal([]byte(res), &ipWhite)
	if err != nil {
		return nil, err
	}
	return ipWhite, nil
}

func UpdateIpWhite(ipWhite *IpWhite) error {
	key := fmt.Sprintf("%s%v", RedisIpWhite, ipWhite.Ip)
	data, err := jsoniter.Marshal(ipWhite)
	if err != nil {
		return err
	}
	tx := cherryredis.GetRds().Set(ctx, key, data, 0)
	if tx.Err() != nil {
		return err
	}
	return nil
}

func DelIpWhite(ip int32) error {
	key := fmt.Sprintf("%s%v", RedisIpWhite, ip)
	tx := cherryredis.GetRds().Del(ctx, key)
	if tx.Err() != nil {
		return tx.Err()
	}
	return nil
}
