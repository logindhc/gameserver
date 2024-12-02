package girdpos

import (
	"container/heap"
	"fmt"
	"math"
)

var (
	// 所有格子坐标
	GridPos    [][]*Pos
	GridPosMap = map[string]*Pos{}
	// 障碍物列表
	Obstacles = make(map[Pos]bool)
)

// 定义格子结构
type Pos struct {
	X int
	Y int
}

func NewPos(x, y int) *Pos {
	return &Pos{
		X: x,
		Y: y,
	}
}

// 初始化菱形对称棋盘
func init() {
	layout := []int{7, 6, 7, 6, 7, 6, 7} // 单数行 7 个，双数行 6 个
	for i, numCells := range layout {
		row := make([]*Pos, numCells)
		index := 0
		if i == 2 || i == 3 { //第2,3行 X从-1开始
			index = -1
		} else if i == 4 || i == 5 { //第4,5行 X从-2开始
			index = -2
		} else if i == 6 { //最后一行从-3开始
			index = -3
		}
		for j := 0; j < numCells; j++ {
			pos := &Pos{X: index, Y: i}
			row[j] = pos
			GridPosMap[fmt.Sprintf("%v_%v", pos.X, pos.Y)] = pos
			index++
		}
		GridPos = append(GridPos, row)
	}

}

/*
获取格子的相邻格子
*/
func (p *Pos) GetAdjacentPos() []Pos {
	var poss []Pos
	//左侧格子
	pos := GridPosMap[fmt.Sprintf("%v_%v", p.X-1, p.Y)]
	if pos != nil {
		poss = append(poss, *pos)
	}
	//左上格子
	pos = GridPosMap[fmt.Sprintf("%v_%v", p.X, p.Y-1)]
	if pos != nil {
		poss = append(poss, *pos)
	}
	//右上格子
	pos = GridPosMap[fmt.Sprintf("%v_%v", p.X+1, p.Y-1)]
	if pos != nil {
		poss = append(poss, *pos)
	}
	//右测格子
	pos = GridPosMap[fmt.Sprintf("%v_%v", p.X+1, p.Y)]
	if pos != nil {
		poss = append(poss, *pos)
	}
	//右下格子
	pos = GridPosMap[fmt.Sprintf("%v_%v", p.X, p.Y+1)]
	if pos != nil {
		poss = append(poss, *pos)
	}
	//左下格子
	pos = GridPosMap[fmt.Sprintf("%v_%v", p.X-1, p.Y+1)]
	if pos != nil {
		poss = append(poss, *pos)
	}
	//会找到不存在格子外的位置
	//poss := []*Pos{
	//	//左侧格子
	//	{p.X - 1, p.Y},
	//	//左上格子
	//	{p.X, p.Y - 1},
	//	//右上格子
	//	{p.X + 1, p.Y - 1},
	//	//右测格子
	//	{p.X + 1, p.Y},
	//	//右下格子
	//	{p.X, p.Y + 1},
	//	//左下格子
	//	{p.X - 1, p.Y + 1},
	//}
	return poss
}
func (p *Pos) Heuristic(b Pos) int {
	dx := p.X - b.X
	dy := p.Y - b.Y
	return int(math.Max(math.Abs(float64(dx)), math.Abs(float64(dy+dx))))
}

func (p *Pos) ContainsAdjacentPos(target Pos) bool {
	adjacentPos := target.GetAdjacentPos()
	for _, pos := range adjacentPos {
		if pos.X == p.X && pos.Y == p.Y {
			return true
		}
	}
	return false
}

/*
AStar算法 寻路
*/
func (p *Pos) AStar(goal Pos) []Pos {
	openList := &PriorityQueue{}
	heap.Init(openList)
	heap.Push(openList, &Node{Pos: *p, G: 0, H: p.Heuristic(goal), F: p.Heuristic(goal)})

	visited := make(map[Pos]bool)

	for openList.Len() > 0 {
		current := heap.Pop(openList).(*Node)
		if current.Pos.X == goal.X && current.Pos.Y == goal.Y {
			var path []Pos
			for node := current; node != nil; node = node.Prev {
				if node.Pos == *p || node.Pos == goal {
					continue
				}
				path = append([]Pos{node.Pos}, path...)
			}
			return path
		}

		visited[current.Pos] = true
		adjacentPos := current.Pos.GetAdjacentPos()
		//fmt.Println(fmt.Sprintf("current:%v,adjacentPos:%v", current.Pos, adjacentPos))
		for _, neighbor := range adjacentPos {
			//目标不是路径，要排除掉
			if visited[neighbor] || (neighbor != goal && Obstacles[neighbor]) {
				continue
			}
			gCost := current.G + 1
			hCost := neighbor.Heuristic(goal)
			heap.Push(openList, &Node{
				Pos:  neighbor,
				G:    gCost,
				H:    hCost,
				F:    gCost + hCost,
				Prev: current,
			})
		}
	}
	//如果找不到，就找邻居的邻居
	for _, pos := range p.GetAdjacentPos() {
		for _, ap := range pos.GetAdjacentPos() {
			star := p.AStar(ap)
			if len(star) > 0 {
				return star
			}
		}
	}
	return nil
}

func (p *Pos) Copy() *Pos {
	return NewPos(p.X, p.Y)
}
