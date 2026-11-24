package main

import (
	"math"
	"strings"
	"testing"
)

func padGrid(rows []string) [][]byte {
	maxLen := 0
	for _, r := range rows {
		if len(r) > maxLen {
			maxLen = len(r)
		}
	}
	out := make([][]byte, len(rows))
	for i, r := range rows {
		row := []byte(r)
		for len(row) < maxLen {
			row = append(row, TileBrick)
		}
		out[i] = row
	}
	return out
}

func allLevels() [][]string {
	return [][]string{level1, level2, level3}
}

func TestLevelsAreWellFormed(t *testing.T) {
	for i, lv := range allLevels() {
		grid := padGrid(lv)
		w, h := len(grid[0]), len(grid)
		if w < 5 || h < 5 {
			t.Fatalf("level %d: too small (%dx%d)", i+1, w, h)
		}
		spawns, exits, marks := 0, 0, 0
		for _, row := range grid {
			if strings.Contains(string(row), " ") {
				t.Fatalf("level %d: blank character in map", i+1)
			}
			for _, c := range row {
				switch c {
				case TileSpawn:
					spawns++
				case TileExit:
					exits++
				case TileEnemy, TileHealth, TileAmmo:
					marks++
				}
			}
		}
		for x := 0; x < w; x++ {
			if !isSolid(grid[0][x]) || !isSolid(grid[h-1][x]) {
				t.Fatalf("level %d: open border at column %d", i+1, x)
			}
		}
		for y := 0; y < h; y++ {
			if !isSolid(grid[y][0]) || !isSolid(grid[y][w-1]) {
				t.Fatalf("level %d: open border at row %d", i+1, y)
			}
		}
		if spawns != 1 {
			t.Fatalf("level %d: want exactly 1 spawn, got %d", i+1, spawns)
		}
		if exits != 1 {
			t.Fatalf("level %d: want exactly 1 exit, got %d", i+1, exits)
		}
		if marks < 3 {
			t.Fatalf("level %d: too few content markers (%d)", i+1, marks)
		}
	}
}

func TestLevelsReachableExit(t *testing.T) {
	for i, lv := range allLevels() {
		grid := padGrid(lv)
		w := &World{Grid: grid, W: len(grid[0]), H: len(grid)}
		start, exit := findKind(w, TileSpawn), findKind(w, TileExit)
		if start == nil || exit == nil {
			t.Fatalf("level %d: missing spawn or exit", i+1)
		}
		if !floodReach(w, int(start.X), int(start.Y), int(exit.X), int(exit.Y)) {
			t.Fatalf("level %d: exit not reachable from spawn", i+1)
		}
	}
}

func findKind(w *World, kind byte) *Pickup {
	for y := 0; y < w.H; y++ {
		for x := 0; x < w.W; x++ {
			if w.Grid[y][x] == kind {
				return &Pickup{X: float64(x) + 0.5, Y: float64(y) + 0.5}
			}
		}
	}
	return nil
}

func floodReach(w *World, sx, sy, tx, ty int) bool {
	seen := map[[2]int]bool{}
	queue := [][2]int{{sx, sy}}
	target := [2]int{tx, ty}
	for len(queue) > 0 {
		c := queue[0]
		queue = queue[1:]
		if seen[c] {
			continue
		}
		seen[c] = true
		if c == target {
			return true
		}
		for _, d := range [4][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}} {
			nx, ny := c[0]+d[0], c[1]+d[1]
			if w.solidAt(nx, ny) || seen[[2]int{nx, ny}] {
				continue
			}
			queue = append(queue, [2]int{nx, ny})
		}
	}
	return false
}

func TestCastStepFindsWall(t *testing.T) {
	grid := padGrid([]string{
		"#####",
		"#...#",
		"#...#",
		"#####",
	})
	dist, _, tile, _ := castStep(grid, 5, 4, 2.5, 1.5, 1, 0)
	if tile != TileBrick || math.Abs(dist-1.5) > 1e-9 {
		t.Fatalf("want dist 1.5 brick, got dist=%v tile=%q", dist, tile)
	}
	if _, side, _, _ := castStep(grid, 5, 4, 2.5, 1.5, 0, 1); side != 1 {
		t.Fatalf("want horizontal side hit, got side=%d", side)
	}
}

func TestLOSClearAndBlocked(t *testing.T) {
	grid := padGrid([]string{
		"#######",
		"#.....#",
		"#..#..#",
		"#.....#",
		"#######",
	})
	if !losClear(grid, 7, 5, 1.5, 1.5, 5.5, 1.5) {
		t.Fatal("LOS should be clear along open corridor")
	}
	if losClear(grid, 7, 5, 1.5, 1.5, 5.5, 3.5) {
		t.Fatal("LOS should be blocked by center wall")
	}
}

func TestTryMoveStopsAtWalls(t *testing.T) {
	grid := padGrid([]string{
		"#####",
		"#...#",
		"#####",
	})
	x, _ := tryMove(grid, 5, 3, 2.5, 1.5, 2.0, 0, 0.2)
	if math.Abs(x-2.5) > 1e-9 {
		t.Fatalf("x should not move into the wall, got %v", x)
	}
	_, y := tryMove(grid, 5, 3, 2.5, 1.5, 0, -1.0, 0.2)
	if y < 1.5 {
		t.Fatalf("should not pass through top wall: y=%v", y)
	}
}

func TestFireHitsEnemyInFront(t *testing.T) {
	g := &Game{}
	w := loadLevel(1, 1)
	g.world = w
	p := &w.Player
	p.Dir = 0
	e := &Enemy{X: p.X + 3, Y: p.Y, HP: 40}
	w.Enemies = append(w.Enemies, e)
	p.Ammo = 10
	w.T = 1
	if !w.fire(g) {
		t.Fatal("shot should hit enemy directly ahead")
	}
	if e.HP >= 40 {
		t.Fatal("enemy should have taken damage")
	}
}
