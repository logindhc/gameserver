package facade

import (
	"gameserver/nodes/game/opcode"
)

type (
	IRes interface {
		BaseID() int
		CheckConfigId(id int) bool
		Get(playerId int64, id int) int
		Enough(playerId int64, id, count int) bool
		Change(playerId int64, id, count int) (int, [][]int)
	}

	IResManager interface {
		Register(res IRes)
		EnoughOne(playerId int64, baseId, id, count int, dup ...int) bool
		Enough(playerId int64, ress [][]int, dup ...int) bool
		AddOne(playerId int64, baseId, id, count int, opcode opcode.Type, remark string, dup ...int) ([][]int, error)
		Add(playerId int64, ress [][]int, opcode opcode.Type, remark string, dup ...int) ([][]int, error)
		DelOne(playerId int64, baseId, id, count int, opcode opcode.Type, remark string, dup ...int) ([][]int, error)
		Del(playerId int64, ress [][]int, opcode opcode.Type, remark string, dup ...int) ([][]int, error)
		GoldEnough(playerId int64, count int) bool
		GoldAdd(playerId int64, count int, opcode opcode.Type, remark string) ([][]int, error)
		GoldDel(playerId int64, count int, opcode opcode.Type, remark string) ([][]int, error)
		MoneyEnough(playerId int64, count int) bool
		MoneyAdd(playerId int64, count int, opcode opcode.Type, remark string) ([][]int, error)
		MoneyDel(playerId int64, count int, opcode opcode.Type, remark string) ([][]int, error)
		DiamondEnough(playerId int64, count int) bool
		DiamondAdd(playerId int64, count int, opcode opcode.Type, remark string) ([][]int, error)
		DiamondDel(playerId int64, count int, opcode opcode.Type, remark string) ([][]int, error)
	}
)
