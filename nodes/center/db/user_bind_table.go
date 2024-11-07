package db

import (
	"fmt"
	cherrySnowflake "gameserver/cherry/extend/snowflake"
	cherryTime "gameserver/cherry/extend/time"
	"gameserver/internal/persistence"
	"gameserver/internal/persistence/repository"
)

// UserBindTable uid绑定第三方平台表
type UserBindTable struct {
	ID       string `gorm:"column:id;primary_key;autoIncrement:false;comment:用户唯一id" json:"id"`
	SdkId    int32  `gorm:"column:sdk_id;comment:sdk id" json:"sdkId"`
	PID      int32  `gorm:"column:pid;comment:平台id" json:"pid"`
	OpenId   string `gorm:"column:open_id;comment:平台帐号open_id" json:"openId"`
	UID      int64  `gorm:"column:uid;comment:游戏唯一ID" json:"uid"`
	BindTime int64  `gorm:"column:bind_time;comment:绑定时间" json:"bindTime"`
}

func (*UserBindTable) TableName() string {
	return "user_bind"
}

var UserBindRepository *repository.LruRepository[string, UserBindTable]

func (p *UserBindTable) InitRepository() {
	UserBindRepository = repository.NewLruRepository[string, UserBindTable](database.GetGameDB(), p.TableName())
	persistence.RegisterRepository(UserBindRepository)
}

func GetUID(pid int32, openId string) (int64, bool) {
	cacheKey := fmt.Sprintf(uidKey, pid, openId)
	val, found := uidCache.GetIfPresent(cacheKey)
	if found == false {
		bind := UserBindRepository.Get(cacheKey)
		if bind == nil {
			return 0, false
		}
		uidCache.Put(cacheKey, bind.UID)
		return bind.UID, true
	}
	return val.(int64), true
}

// BindUID 绑定UID
func BindUID(sdkId, pid int32, openId string) (int64, bool) {
	// TODO 根据 platformType的配置要求，决定查询UID的方式：
	// 条件1: platformType + openId查询，是否存在uid
	// 条件2: pid + openId查询，是否存在uid

	uid, ok := GetUID(pid, openId)
	if ok {
		return uid, true
	}

	userBind := &UserBindTable{
		ID:       fmt.Sprintf(uidKey, pid, openId),
		SdkId:    sdkId,
		PID:      pid,
		OpenId:   openId,
		UID:      cherrySnowflake.NextId(),
		BindTime: cherryTime.Now().ToMillisecond(),
	}

	cacheKey := fmt.Sprintf(uidKey, pid, openId)
	uidCache.Put(cacheKey, userBind.UID)
	add := UserBindRepository.Add(userBind)
	if add == nil {
		return 0, false
	}
	return userBind.UID, true
}
