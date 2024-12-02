package girdpos

import (
	"fmt"
	"math"
	"math/rand"
	"testing"
	"time"
)

type Faction int

const (
	Red Faction = iota
	Blue
)

type UnitType int

const (
	Melee UnitType = iota
	Ranged
)

type Unit struct {
	ID         int
	Name       string
	HP         int
	Attack     int
	Defense    int
	Range      int // 远程攻击范围
	Type       UnitType
	Skills     []Skill
	Buffs      []Buff
	Position   Pos
	IsAlive    bool
	Faction    Faction
	Cooldown   uint64 // 间隔时间
	LastAction uint64 // 上次行动时间
	target     *Unit
}

type Skill struct {
	Name       string
	Damage     int
	Cooldown   time.Duration
	CurrentCD  time.Duration
	EffectFunc func(attacker, target *Unit, report *BattleReport)
}

type Buff struct {
	Name       string
	Duration   int
	EffectFunc func(target *Unit, report *BattleReport)
}

type BattleReport struct {
	Actions []string
}

type Battle struct {
	Units  []Unit
	Report BattleReport
	Time   uint64 // 当前时间轮
}

func NewBattle(units []Unit) *Battle {
	return &Battle{
		Units:  units,
		Report: BattleReport{},
	}
}

func (b *Battle) Start() {
	for b.isBattleOngoing() {
		b.executeRound()
	}
	fmt.Println("Battle finished!")
	b.printReport()
}

func (b *Battle) executeRound() {
	//更新时间轮
	b.Time++
	for i := range b.Units {
		unit := &b.Units[i]
		if unit.IsAlive && b.timeToAct(unit) {
			b.processAction(unit)
		}
	}

}

// 判断单位是否满足行动条件
func (b *Battle) timeToAct(unit *Unit) bool {
	unit.LastAction++
	// 判断单位的技能冷却是否完成（如果有冷却时间），以及是否处于行动状态
	if unit.LastAction >= unit.Cooldown {
		return true
	}
	return false
}

func (b *Battle) isBattleOngoing() bool {
	redAlive, blueAlive := false, false
	for _, unit := range b.Units {
		if unit.IsAlive {
			if unit.Faction == Red {
				redAlive = true
			} else {
				blueAlive = true
			}
		}
	}
	return redAlive && blueAlive
}

func (b *Battle) processAction(unit *Unit) {
	// 更新单位的最后行动时间
	unit.LastAction = 0
	var target *Unit
	if unit.target == nil || !unit.target.IsAlive {
		// 查找敌方目标
		target = b.selectEnemyTarget(unit)
		if target == nil {
			return
		}
		fmt.Println(fmt.Sprintf("%d selects %s -> target %s", b.Time, unit.Name, target.Name))
		unit.target = target
	} else {
		target = unit.target
	}
	// 判断距离是否合适，远程单位与近战单位有不同策略
	if unit.Type == Melee {
		// 近战单位：攻击范围内直接攻击
		if b.calculateDistance(unit.Position, target.Position) == 1 {
			unit.AttackTarget(target)
			b.addEvent(unit.ID, fmt.Sprintf("attacks %s Hp %d  target %s  Hp %d", unit.Name, unit.HP, target.Name, target.HP))
		} else {
			// 否则移动到目标附近
			b.moveUnitToAttack(unit, target)
		}
	} else if unit.Type == Ranged {
		// 远程单位：判断是否在攻击范围内
		if b.calculateDistance(unit.Position, target.Position) <= unit.Range {
			unit.AttackTarget(target)
			b.addEvent(unit.ID, fmt.Sprintf("attacks distance %s Hp %d  target %s  Hp %d", unit.Name, unit.HP, target.Name, target.HP))
		} else {
			// 否则移动到攻击范围内
			b.moveUnitToAttack(unit, target)
		}
	}
}

func (b *Battle) selectEnemyTarget(unit *Unit) *Unit {
	var targets []*Unit
	for i := range b.Units {
		if b.Units[i].IsAlive && b.Units[i].Faction != unit.Faction {
			targets = append(targets, &b.Units[i])
		}
	}
	if len(targets) == 0 {
		return nil
	}

	// 选择敌人的策略：最近、最远或随机选择
	// 1. 最近敌人
	var target *Unit
	switch rand.Intn(3) { // 随机选择1到3之间的策略
	case 0:
		target = b.selectNearestEnemy(unit, targets)
	case 1:
		target = b.selectFarthestEnemy(unit, targets)
	default:
		target = b.selectRandomEnemy(targets)
	}
	return target
}

func (b *Battle) selectNearestEnemy(unit *Unit, targets []*Unit) *Unit {
	var closestTarget *Unit
	minDistance := math.MaxInt32
	for _, target := range targets {
		distance := b.calculateDistance(unit.Position, target.Position)
		if distance < minDistance {
			closestTarget = target
			minDistance = distance
		}
	}
	return closestTarget
}

func (b *Battle) selectFarthestEnemy(unit *Unit, targets []*Unit) *Unit {
	var farthestTarget *Unit
	maxDistance := 0
	for _, target := range targets {
		distance := b.calculateDistance(unit.Position, target.Position)
		if distance > maxDistance {
			farthestTarget = target
			maxDistance = distance
		}
	}
	return farthestTarget
}

func (b *Battle) selectRandomEnemy(targets []*Unit) *Unit {
	return targets[rand.Intn(len(targets))]
}

func (b *Battle) calculateDistance(pos1, pos2 Pos) int {
	heuristic := pos1.Heuristic(pos2)
	return heuristic
}

func (b *Battle) moveUnitToAttack(unit, target *Unit) {
	// 计算并移动到目标附近的格子
	// 这里只是简单示意，实际上可以使用路径规划算法（如A*）来寻找最短路径
	newPos := unit.Position.AStar(target.Position)
	if len(newPos) < 1 {
		return
	}
	old := unit.Position
	delete(Obstacles, old)
	// 更新单位位置
	unit.Position = newPos[0]
	Obstacles[unit.Position] = true
	b.addEvent(unit.ID, fmt.Sprintf("%s moves (%d,%d) to (%d, %d)", unit.Name, old.X, old.Y, unit.Position.X, unit.Position.Y))
	//fmt.Println(fmt.Sprintf("%s moves (%d,%d) to position (%d, %d)", unit.Name, old.X, old.Y, unit.Position.X, unit.Position.Y))
}

func (b *Battle) addEvent(unitID int, action string) {
	b.Report.Actions = append(b.Report.Actions, fmt.Sprintf("time %d -> %s", b.Time, action))
}

func (b *Battle) printReport() {
	fmt.Println("\n--- Battle Report ---")
	for i, action := range b.Report.Actions {
		fmt.Println(i, action)
	}
}

func (u *Unit) AttackTarget(target *Unit) {
	// 执行攻击
	damage := u.Attack - target.Defense
	if damage < 0 {
		damage = 0
	}
	target.HP -= damage
	if target.HP <= 0 {
		target.HP = 0
		target.IsAlive = false
		fmt.Println(fmt.Sprintf("%s killed %s \n", u.Name, target.Name))
	}
}

func Test_main(t *testing.T) {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	// 初始化红方和蓝方的单位
	units := []Unit{
		{ID: 1, Name: "红1", HP: 1000, Attack: 200, Defense: 5, Faction: Red, IsAlive: true, Position: Pos{X: 0, Y: 1}, Type: Melee, Cooldown: 30, Range: 1},
		{ID: 2, Name: "红2", HP: 800, Attack: 150, Defense: 2, Faction: Red, IsAlive: true, Position: Pos{X: 0, Y: 2}, Type: Ranged, Cooldown: 20, Range: 2},
		{ID: 3, Name: "红3", HP: 600, Attack: 205, Defense: 3, Faction: Red, IsAlive: true, Position: Pos{X: 0, Y: 3}, Type: Ranged, Cooldown: 50, Range: 3},
		{ID: 4, Name: "蓝1", HP: 1000, Attack: 200, Defense: 5, Faction: Blue, IsAlive: true, Position: Pos{X: 3, Y: 1}, Type: Melee, Cooldown: 30, Range: 1},
		{ID: 5, Name: "蓝2", HP: 800, Attack: 150, Defense: 2, Faction: Blue, IsAlive: true, Position: Pos{X: 3, Y: 2}, Type: Ranged, Cooldown: 20, Range: 2},
		{ID: 6, Name: "蓝3", HP: 600, Attack: 205, Defense: 3, Faction: Blue, IsAlive: true, Position: Pos{X: 3, Y: 3}, Type: Ranged, Cooldown: 50, Range: 3},
	}

	// 开始战斗
	battle := NewBattle(units)
	battle.Start()
}
