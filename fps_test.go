package main

import (
	"math"
	"math/rand"
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
	return [][]string{level1, level2, level3, level4, level5}
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
				case TileHealth, TileAmmo:
					marks++
				}
				if isEnemy(c) {
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
	e := &Enemy{X: p.X + 3, Y: p.Y, Kind: TileEnemy, HP: 40, Score: 100}
	w.Enemies = append(w.Enemies, e)
	p.Ammo[0] = 10
	w.T = 1
	if !w.fire(g) {
		t.Fatal("shot should hit enemy directly ahead")
	}
	if e.HP >= 40 {
		t.Fatal("enemy should have taken damage")
	}
}

func TestIonSplashDamagesGroup(t *testing.T) {
	w := &World{Grid: padGrid([]string{
		"########",
		"#......#",
		"#......#",
		"########",
	}), W: 8, H: 4, rng: rand.New(rand.NewSource(1))}
	g := &Game{}
	g.world = w
	p := &w.Player
	p.X, p.Y, p.Dir = 2.5, 1.5, 0
	p.Weapon = 3
	p.Ammo[3] = 5
	w.T = 2

	e1 := &Enemy{X: 5.5, Y: 1.5, Kind: TileEnemy, HP: 100, Score: 100}
	e2 := &Enemy{X: 5.5, Y: 2.3, Kind: TileEnemy, HP: 100, Score: 100}
	far := &Enemy{X: 2.5, Y: 2.5, Kind: TileEnemy, HP: 100, Score: 100}
	w.Enemies = append(w.Enemies, e1, e2, far)

	if !w.fire(g) {
		t.Fatal("ion shot should hit something")
	}
	if p.Ammo[3] != 4 {
		t.Fatalf("ion shot should cost one shell, got %d", p.Ammo[3])
	}
	if !e1.Dead {
		t.Fatalf("direct target should die: hp=%d", e1.HP)
	}
	if e2.Dead {
		t.Fatal("splash target should be wounded, not killed")
	}
	if e2.HP >= 100 {
		t.Fatalf("splash should wound nearby enemy, got hp=%d", e2.HP)
	}
	if far.HP != 100 {
		t.Fatalf("enemy outside splash radius must not be hit, got hp=%d", far.HP)
	}
	if g.kills != 1 || g.score != 100 {
		t.Fatalf("want 1 kill scored, got kills=%d score=%d", g.kills, g.score)
	}
}

func TestParticlesDecayAndCap(t *testing.T) {
	w := loadLevel(1, 1)
	w.spawnBlood(w.Player.X, w.Player.Y, TileEnemy, 20)
	if len(w.Parts) != 20 {
		t.Fatalf("want 20 particles, got %d", len(w.Parts))
	}
	for i := 0; i < 300; i++ {
		w.update(1.0/60.0, &Game{})
	}
	if len(w.Parts) != 0 {
		t.Fatalf("particles should vanish after their lifetime, left %d", len(w.Parts))
	}
	for i := 0; i < 100; i++ {
		w.spawnBlood(1.5, 1.5, TileEnemy, 20)
	}
	if len(w.Parts) > maxParts {
		t.Fatalf("particle cap exceeded: %d", len(w.Parts))
	}
}

func TestSurvivalArenaFullyConnected(t *testing.T) {
	grid := padGrid(levelSurv)
	w := &World{Grid: grid, W: len(grid[0]), H: len(grid)}
	start := findKind(w, TileSpawn)
	if start == nil {
		t.Fatal("arena has no spawn")
	}
	if findKind(w, TileExit) != nil {
		t.Fatal("arena must not contain an exit")
	}
	total := 0
	for y := 0; y < w.H; y++ {
		for x := 0; x < w.W; x++ {
			if !w.solidAt(x, y) {
				total++
			}
		}
	}
	reached := 0
	for y := 0; y < w.H; y++ {
		for x := 0; x < w.W; x++ {
			if w.solidAt(x, y) {
				continue
			}
			if floodReach(w, int(start.X), int(start.Y), x, y) {
				reached++
			}
		}
	}
	if reached != total {
		t.Fatalf("arena has unreachable floor cells: %d/%d reachable", reached, total)
	}
}

func TestSurvivalWavesProgress(t *testing.T) {
	g := &Game{}
	g.mode = modeSurvive
	g.startLevel(1)
	g.state = statePlay
	w := g.world
	if !w.Survive || w.Wave != 1 || w.ToSpawn != 6 {
		t.Fatalf("survival should start on wave 1 with 6 enemies to spawn (wave=%d toSpawn=%d)", w.Wave, w.ToSpawn)
	}
	for f := 0; f < 60*8; f++ {
		w.update(1.0/60.0, g)
	}
	if w.ToSpawn != 0 || w.aliveCount() < 5 {
		t.Fatalf("wave 1 should be fully spawned (toSpawn=%d alive=%d)", w.ToSpawn, w.aliveCount())
	}
	for _, e := range w.Enemies {
		e.Dead = true
		e.DeadT = 2
	}
	w.update(1.0/60.0, g)
	for f := 0; f < 60*4; f++ {
		w.update(1.0/60.0, g)
	}
	if w.Wave != 2 {
		t.Fatalf("clearing the wave should trigger wave 2, got %d", w.Wave)
	}
}
