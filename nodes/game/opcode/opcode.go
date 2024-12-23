package opcode

type Type int

const (
	Gm Type = iota + 101 //GM操作
)
const (
	PlayerInit Type = iota + 10001 //初始化
)

const (
	ItemUse Type = iota + 10101 //使用道具
)
const (
	HeroUp Type = iota + 10201 //英雄升级
)
