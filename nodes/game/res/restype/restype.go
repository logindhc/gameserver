package restype

const (
	Currency  = iota + 1 //货币
	Item                 //道具
	Hero                 //英雄
	BattleExp            //战斗经验(前端用)
	Box                  //宝箱
)

type CurrencyType int

const (
	Gold    CurrencyType = iota + 1 //金币
	Money                           //银币
	Diamond                         //钻石
)
