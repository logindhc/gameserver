package facade

type IDrop interface {
	Drop(playerId int64, dropGroupId int, count int) ([][]int, error)
	DropRand(playerId int64, dropGroupId int, min, max int) ([][]int, error)
}
