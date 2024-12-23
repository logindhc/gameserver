package data

import (
	cherryError "gameserver/cherry/error"
	cherryLogger "gameserver/cherry/logger"
)

var HeroConfig = &heroConfig{}

type (
	HeroRow struct {
		Id              int     `json:"id"`              // 编号
		HeroName        int     `json:"heroName"`        // 英雄名称
		Name            string  `json:"name"`            // 英雄名称
		HeroGroup       int     `json:"heroGroup"`       // 英雄组
		HeroQuality     int     `json:"heroQuality"`     // 品质
		HeroLevel       int     `json:"heroLevel"`       // 等级
		UpCost          int     `json:"upCost"`          // 升级费用
		PieceId         int     `json:"pieceId"`         // 消耗碎片
		PieceCount      int     `json:"pieceCount"`      // 升级碎片数量
		ConveterCount   int     `json:"conveterCount"`   // 重复获取转碎片数量
		MaxHp           int     `json:"maxHp"`           // 最大生命
		MaxHpPer        int     `json:"maxHpPer"`        // 最大生命%（千
		Armor           int     `json:"armor"`           // 护甲
		ArmorPer        int     `json:"armorPer"`        // 护甲%（千
		AagicResist     int     `json:"aagicResist"`     // 魔抗
		AagicResistPer  int     `json:"aagicResistPer"`  // 魔抗%（千）
		AttackDamage    int     `json:"attackDamage"`    // 攻击
		AttackDamagePer int     `json:"attackDamagePer"` // 攻击%（千）
		AbilityPower    int     `json:"abilityPower"`    // 法攻
		AbilityPowerPer int     `json:"abilityPowerPer"` // 法攻%（千）
		AtkSpeed        float32 `json:"atkSpeed"`        // 攻速
		AtkSpeedPer     int     `json:"atkSpeedPer"`     // 攻速%(千）
		CritRate        int     `json:"critRate"`        // 暴击率（千）
		CritDamage      int     `json:"critDamage"`      // 暴击伤害（千）
		AtkRange        int     `json:"atkRange"`        // 攻击距离
		IniMana         int     `json:"iniMana"`         // 初始法力
		MaxMana         int     `json:"maxMana"`         // 最大法力
		MoveSpeed       int     `json:"moveSpeed"`       // 移动速度
		HpRegen         int     `json:"hpRegen"`         // 生命回复
		ManaRegen       int     `json:"manaRegen"`       // 法力回复
		ArPen           int     `json:"arPen"`           // 护甲穿透
		ArPenRate       int     `json:"arPenRate"`       // 护甲穿透率（千）
		ApPen           int     `json:"apPen"`           // 法术穿透
		ApPenRate       int     `json:"apPenRate"`       // 法术穿透率（千）
		AdLifeSteal     int     `json:"adLifeSteal"`     // 物理吸血%（千）
		ApLifeSteal     int     `json:"apLifeSteal"`     // 法术吸血%（千）
		ResPer          int     `json:"resPer"`          // 韧性%（千）
		NormalAttack    int     `json:"normalAttack"`    // 普通攻击
		BigSkill        int     `json:"bigSkill"`        // 大招
		RoleIcon        int     `json:"roleIcon"`        // 头像
		RolePic         int     `json:"rolePic"`         // 角色图
		Size            int     `json:"size"`            // 尺寸
		Grid            int     `json:"grid"`            // 格子
	}

	heroConfig struct {
		maps  map[int]*HeroRow
		maps2 map[[2]int]*HeroRow
	}
)

func (c *heroConfig) Init() {
	c.maps = make(map[int]*HeroRow)
}

func (c *heroConfig) OnLoad(maps interface{}, _ bool) (int, error) {
	list, ok := maps.([]interface{}) // map结构：maps.(map[string]interface{})
	if !ok {
		return 0, cherryError.Error("maps convert to map[string]interface{} error.")
	}

	loadMaps := make(map[int]*HeroRow)
	loadMap2s := make(map[[2]int]*HeroRow)
	for index, data := range list {
		loadConfig := &HeroRow{}
		err := DecodeData(data, loadConfig)
		if err != nil {
			cherryLogger.Warnf("decode error. [id = %v, %v], err = %s", index, loadConfig, err)
			continue
		}
		loadMaps[loadConfig.Id] = loadConfig
		loadMap2s[[2]int{loadConfig.HeroGroup, loadConfig.HeroLevel}] = loadConfig
	}
	c.maps = loadMaps
	c.maps2 = loadMap2s

	return len(list), nil
}

func (c *heroConfig) Name() string {
	return "HeroConfig"
}

func (c *heroConfig) OnAfterLoad(_ bool) {}

func (c *heroConfig) Get(key int) (*HeroRow, bool) {
	row, ok := c.maps[key]
	return row, ok
}
func (c *heroConfig) GetByGroupLevel(groupId, level int) (*HeroRow, bool) {
	row, ok := c.maps2[[2]int{groupId, level}]
	return row, ok
}
func (c *heroConfig) List() []*HeroRow {
	var list []*HeroRow
	for _, row := range c.maps {
		list = append(list, row)
	}
	return list
}
