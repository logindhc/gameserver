package girdpos

import (
	"fmt"
	"testing"
	"time"
)

func TestBattle(t *testing.T) {
	PrintDiamondBoard()
	start := Pos{3, 1}
	goal := Pos{0, 1}
	goal2 := Pos{0, 0}
	Obstacles[goal2] = true
	Obstacles[Pos{1, 1}] = true
	Obstacles[Pos{-1, 2}] = true
	Obstacles[Pos{0, 2}] = true
	Obstacles[Pos{1, 0}] = true
	heuristic := start.Heuristic(goal)
	fmt.Println("Heuristic:", heuristic)
	heuristic2 := start.Heuristic(goal2)
	fmt.Println("heuristic2:", heuristic2)
	second := time.Now()
	path := start.AStar(goal)
	for _, p := range path {
		fmt.Printf("(%d, %d) -> ", p.X, p.Y)
		Obstacles[p] = true
	}
	fmt.Println("Goal reached!", time.Since(second))
	path2 := start.AStar(goal2)
	for _, p := range path2 {
		fmt.Printf("(%d, %d) -> ", p.X, p.Y)
	}
	fmt.Println("Goal2 reached!", time.Since(second))
}

// 打印棋盘
func PrintDiamondBoard() {
	for _, row := range GridPos {
		spaces := (7 - len(row)) // 菱形对齐的空格数
		for s := 0; s < spaces; s++ {
			fmt.Print("  ")
		}
		for _, cell := range row {
			fmt.Printf("(%d,%d)", cell.X, cell.Y)
		}
		fmt.Println()
	}
}
