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
	return [][]string{level1, level2, level3, level4, level5, levelDesertMap}
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
	p.Weapon = 4
	p.Ammo[4] = 5
	w.T = 2

	e1 := &Enemy{X: 5.5, Y: 1.5, Kind: TileEnemy, HP: 100, Score: 100}
	e2 := &Enemy{X: 5.5, Y: 2.3, Kind: TileEnemy, HP: 100, Score: 100}
	far := &Enemy{X: 2.5, Y: 2.5, Kind: TileEnemy, HP: 100, Score: 100}
	w.Enemies = append(w.Enemies, e1, e2, far)

	if !w.fire(g) {
		t.Fatal("ion shot should hit something")
	}
	if p.Ammo[4] != 4 {
		t.Fatalf("ion shot should cost one shell, got %d", p.Ammo[4])
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

func TestBossGatesExit(t *testing.T) {
	g := &Game{}
	g.startLevel(5)
	g.state = statePlay
	w := g.world
	bosses := 0
	for _, e := range w.Enemies {
		if e.Kind == TileBoss {
			bosses++
		}
	}
	if bosses != 1 {
		t.Fatalf("level 5 should host exactly 1 guardian, got %d", bosses)
	}
	ex := findKind(w, TileExit)
	w.Player.X, w.Player.Y = ex.X, ex.Y
	w.update(1.0/60.0, g)
	if w.Cleared {
		t.Fatal("exit must stay locked while the guardian lives")
	}
	w.damageEnemy(w.boss(), 1000, g)
	w.update(1.0/60.0, g)
	if !w.Cleared {
		t.Fatal("exit should unlock once the guardian is dead")
	}
}

func TestBossVolleyFiresTripleProjectiles(t *testing.T) {
	w := &World{Grid: padGrid([]string{
		"#########",
		"#.......#",
		"#.......#",
		"#.......#",
		"#########",
	}), W: 9, H: 5, rng: rand.New(rand.NewSource(3))}
	g := &Game{}
	g.world = w
	w.Player.X, w.Player.Y = 2.5, 1.5
	boss := &Enemy{X: 6.5, Y: 2.5, Kind: TileBoss, HP: 900, Score: 2000, ShootT: 0.01}
	w.Enemies = append(w.Enemies, boss)
	w.update(1.0/60.0, g)
	if len(w.Projs) != 3 {
		t.Fatalf("guardian volley should fire 3 projectiles, got %d", len(w.Projs))
	}
	for _, pr := range w.Projs {
		if pr.Dmg != 14 {
			t.Fatalf("guardian projectiles should deal 14 damage, got %d", pr.Dmg)
		}
	}
}

func TestSniperReachesBeyondStandardRange(t *testing.T) {
	rows := []string{
		"##############################",
		"#............................#",
		"##############################",
	}
	w := &World{Grid: padGrid(rows), W: 30, H: 3, rng: rand.New(rand.NewSource(1))}
	g := &Game{}
	g.world = w
	p := &w.Player
	p.X, p.Y, p.Dir = 2.5, 1.5, 0
	e := &Enemy{X: 26.5, Y: 1.5, Kind: TileEnemy, HP: 400, Score: 100}
	w.Enemies = append(w.Enemies, e)

	p.Weapon = 0
	p.Ammo[0] = 10
	w.T = 1
	if w.fire(g) {
		t.Fatal("pistol should not reach a target at 24 tiles")
	}

	p.Weapon = 3
	p.Ammo[3] = 5
	w.T = 3
	if !w.fire(g) {
		t.Fatal("sniper should hit a target at 24 tiles")
	}
	if e.HP >= 400 {
		t.Fatal("sniper shot should damage distant enemy")
	}
}

func TestCSBuyAndCash(t *testing.T) {
	g := &Game{}
	g.mode = modeCS
	g.startLevel(1)
	w := g.world
	if !w.Survive || !w.CS {
		t.Fatal("CS level should set Survive and CS flags")
	}
	if w.cash != 800 {
		t.Fatalf("CS should start with $800, got %d", w.cash)
	}
	p := &w.Player
	before := p.Ammo[1] + p.Reserve[1]
	g.tryBuy(300, 1, 12)
	if w.cash != 500 {
		t.Fatalf("cash should drop to 500, got %d", w.cash)
	}
	if p.Ammo[1]+p.Reserve[1] != before+12 {
		t.Fatalf("pump ammo should increase, got %d", p.Ammo[1]+p.Reserve[1])
	}
	g.tryBuy(300, 1, 12)
	if w.cash != 200 {
		t.Fatalf("second buy should also cost 300, got %d", w.cash)
	}
	p.Reserve[2] = weapons[2].MaxAmmo - p.Ammo[2]
	ammoFull := p.Ammo[2] + p.Reserve[2]
	g.tryBuy(700, 2, 40)
	if w.cash != 200 {
		t.Fatalf("full-ammo buy should be refused, got %d", w.cash)
	}
	if p.Ammo[2]+p.Reserve[2] != ammoFull {
		t.Fatal("full-ammo buy should not change ammo")
	}
}

func TestCSKillGivesCash(t *testing.T) {
	g := &Game{}
	g.mode = modeCS
	g.startLevel(1)
	w := g.world
	w.Enemies = nil
	e := &Enemy{X: 3.5, Y: 3.5, Kind: TileEnemy, HP: 40, Score: 100}
	w.Enemies = append(w.Enemies, e)
	cashBefore := w.cash
	w.damageEnemy(e, 1000, g)
	if !e.Dead {
		t.Fatal("enemy should die")
	}
	if w.cash != cashBefore+60 {
		t.Fatalf("imp kill should pay $60, got %d", w.cash-cashBefore)
	}
	b := &Enemy{X: 4.5, Y: 3.5, Kind: TileBoss, HP: 900, Score: 2000}
	w.Enemies = append(w.Enemies, b)
	cashBefore = w.cash
	w.damageEnemy(b, 1000, g)
	if w.cash != cashBefore+300 {
		t.Fatalf("boss kill should pay $300, got %d", w.cash-cashBefore)
	}
}

func TestReloadPullsFromReserve(t *testing.T) {
	rows := []string{
		"#####",
		"#...#",
		"#####",
	}
	w := &World{Grid: padGrid(rows), W: 5, H: 3, rng: rand.New(rand.NewSource(1))}
	g := &Game{}
	g.world = w
	p := &w.Player
	p.X, p.Y = 2.5, 1.5
	p.Ammo[0], p.Reserve[0] = 2, 12

	w.startReload()
	if p.ReloadT <= 0 {
		t.Fatal("reload should start when magazine is not full")
	}
	for f := 0; f < 120 && p.ReloadT > 0; f++ {
		w.update(1.0/60.0, g)
	}
	if p.ReloadT != 0 {
		t.Fatal("reload should finish within 2 seconds")
	}
	if p.Ammo[0] != weapons[0].MagSize {
		t.Fatalf("magazine should be full, got %d", p.Ammo[0])
	}
	if p.Reserve[0] != 12-(weapons[0].MagSize-2) {
		t.Fatalf("reserve should pay for the reload, got %d", p.Reserve[0])
	}

	p.Ammo[0] = 0
	p.Reserve[0] = 0
	w.T += 1
	w.startReload()
	if p.ReloadT != 0 {
		t.Fatal("reload must not start without reserve")
	}
	if w.fire(g) {
		t.Fatal("empty weapon without reserve should not fire")
	}
	p.Reserve[0] = 3
	w.T += 1
	if w.fire(g) {
		t.Fatal("empty weapon should trigger reload instead of firing")
	}
	if p.ReloadT <= 0 {
		t.Fatal("firing empty should start a reload")
	}
}

func TestColossusSecondBoss(t *testing.T) {
	w := loadLevel(levelSurvive, 1)
	w.startWave(10)
	if !w.Boss2Wave {
		t.Fatal("survival wave 10 should call the colossus")
	}
	w3 := loadLevel(levelSurvive, 2)
	w3.startWave(5)
	if w3.Boss2Wave || !w3.BossWave {
		t.Fatal("wave 5 should only bring the gardien")
	}
	w2 := loadLevel(levelCS, 1)
	w2.startWave(6)
	if !w2.Boss2Wave {
		t.Fatal("CS wave 6 should call the colossus")
	}
	w2.Enemies = nil
	g := &Game{}
	c := &Enemy{X: 3.5, Y: 3.5, Kind: TileColossus, HP: 1400, Score: 5000}
	w2.Enemies = append(w2.Enemies, c)
	cashBefore := w2.cash
	w2.damageEnemy(c, 5000, g)
	if !c.Dead {
		t.Fatal("colossus should die")
	}
	if w2.cash != cashBefore+500 {
		t.Fatalf("colossus kill should pay $500, got %d", w2.cash-cashBefore)
	}
	if kind := enemyKinds[TileColossus]; kind.Volley != 5 {
		t.Fatalf("colossus should fire a 5-projectile volley, got %d", kind.Volley)
	}
}

func TestBuySkinCostsAndStacks(t *testing.T) {
	g := &Game{}
	g.mode = modeCS
	g.startLevel(1)
	w := g.world
	p := &w.Player

	g.buySkin(0)
	if p.Skins[0] != 1 || w.cash != 800-400 {
		t.Fatalf("first skin should cost 400, got skin=%d cash=%d", p.Skins[0], w.cash)
	}
	g.buySkin(0)
	if p.Skins[0] != 1 {
		t.Fatal("skin purchase must be refused without cash")
	}
	w.cash = 5000
	g.buySkin(0)
	if p.Skins[0] != 2 || w.cash != 4300 {
		t.Fatalf("second skin should cost 700, got skin=%d cash=%d", p.Skins[0], w.cash)
	}
	g.buySkin(9)
	if p.Skins[0] != 2 {
		t.Fatal("out-of-range weapon must be ignored")
	}
	for i := 0; i < 4; i++ {
		g.buySkin(0)
	}
	if p.Skins[0] != 5 {
		t.Fatalf("skins should stack to the last finish, got %d", p.Skins[0])
	}
	g.buySkin(0)
	if p.Skins[0] != 5 {
		t.Fatal("last finish should be capped")
	}
}

func TestApplySkinRemap(t *testing.T) {
	if applySkin(0x808080, 0) != 0x808080 {
		t.Fatal("standard skin should keep colors")
	}
	ch := applySkin(0x808080, 1)
	if (ch & 0xff) <= (ch>>16)&0xff {
		t.Fatal("chrome should push a blue tint")
	}
	gold := applySkin(0x808080, 2)
	if (gold>>16)&0xff <= gold&0xff {
		t.Fatal("gold should push red above blue")
	}
	neon := applySkin(0x808080, 3)
	if neon&0xff <= (neon>>16)&0xff {
		t.Fatal("neon should push blue above red")
	}
	if applySkin(0x123456, 42) != 0x123456 {
		t.Fatal("unknown skin should be identity")
	}
}

func TestChestOpensAndRewards(t *testing.T) {
	rows := []string{
		"#####",
		"#...#",
		"#####",
	}
	w := &World{Grid: padGrid(rows), W: 5, H: 3, rng: rand.New(rand.NewSource(3))}
	g := &Game{}
	g.world = w
	p := &w.Player
	p.X, p.Y, p.HP = 2.5, 1.5, 60
	w.Pickups = append(w.Pickups, Pickup{X: 2.5, Y: 1.5, Kind: TileChest})
	w.Pickups = append(w.Pickups, Pickup{X: 4.2, Y: 1.5, Kind: TileChest})

	beforeHP, beforeScore, beforeCash := p.HP, g.score, w.cash
	beforeAmmo := p.Ammo[0] + p.Reserve[0]
	w.update(1.0/60.0, g)

	if !w.Pickups[0].Opened {
		t.Fatal("chest should open on contact")
	}
	if w.Pickups[1].Opened {
		t.Fatal("distant chest must stay closed")
	}
	looted := p.HP != beforeHP || g.score != beforeScore || w.cash != beforeCash ||
		p.Ammo[0]+p.Reserve[0] != beforeAmmo
	if !looted {
		t.Fatal("opened chest should give something")
	}
	hp2, sc2 := p.HP, g.score
	w.update(1.0/60.0, g)
	if p.HP != hp2 || g.score != sc2 {
		t.Fatal("opened chest must not loot twice")
	}
}

func TestWavesSpawnChests(t *testing.T) {
	w := loadLevel(levelSurvive, 5)
	w.startWave(3)
	found := false
	for _, pk := range w.Pickups {
		if pk.Kind == TileChest && !pk.Opened {
			found = true
		}
	}
	if !found {
		t.Fatal("each wave should drop a chest")
	}
}

func TestVehicleMountRampage(t *testing.T) {
	rows := []string{
		"#########",
		"#.......#",
		"#########",
	}
	w := &World{Grid: padGrid(rows), W: 9, H: 3, rng: rand.New(rand.NewSource(1))}
	g := &Game{}
	g.world = w
	p := &w.Player
	p.X, p.Y, p.Dir = 3.5, 1.5, 0
	w.Vehicle = Vehicle{X: 7.0, Y: 1.5}

	w.toggleVehicle()
	if w.VehMounted {
		t.Fatal("mount must fail when the vehicle is far")
	}
	w.Vehicle = Vehicle{X: 4.0, Y: 1.5}
	w.toggleVehicle()
	if !w.VehMounted {
		t.Fatal("mount should work when close")
	}
	if w.fire(g) {
		t.Fatal("cannot fire while driving")
	}

	e := &Enemy{X: 3.6, Y: 1.5, Kind: TileEnemy, HP: 100, Score: 100}
	w.Enemies = append(w.Enemies, e)
	w.vehRampage(0.1, g)
	if e.HP >= 100 {
		t.Fatal("vehicle should crush the enemy")
	}
	w.vehRampage(0.1, g)
	if e.HP > 0 && e.Dead {
		t.Fatal("dead enemy must not be crushed twice")
	}

	w.toggleVehicle()
	if w.VehMounted {
		t.Fatal("dismount should work")
	}
	if w.Vehicle.X < 2.0 || w.Vehicle.X > 3.5 {
		t.Fatalf("vehicle should be parked behind the player, got %v", w.Vehicle.X)
	}
}

func TestDesertLevelHugeWithVehicleAndChests(t *testing.T) {
	g := &Game{}
	g.startLevel(7)
	w := g.world
	if w.W < 40 || w.H < 20 {
		t.Fatalf("desert map should be huge, got %dx%d", w.W, w.H)
	}
	if w.Vehicle.X == 0 && w.Vehicle.Y == 0 {
		t.Fatal("desert must contain the VTT")
	}
	chests := 0
	for _, pk := range w.Pickups {
		if pk.Kind == TileChest {
			chests++
		}
	}
	if chests < 2 {
		t.Fatalf("desert should hide several chests, got %d", chests)
	}
	kinds := map[byte]bool{}
	for _, e := range w.Enemies {
		kinds[e.Kind] = true
	}
	if !kinds[TileColossus] || !kinds[TileBoss] {
		t.Fatal("desert finale must host both bosses")
	}
	if len(w.Exits) != 1 {
		t.Fatalf("desert must keep exactly one exit, got %d", len(w.Exits))
	}
}

func TestTankMountFiresCannon(t *testing.T) {
	g := &Game{}
	g.startLevel(7)
	g.state = statePlay
	w := g.world
	p := &w.Player
	w.Tank = Vehicle{X: p.X + 1.0, Y: p.Y}
	w.toggleVehicle()
	if !w.VehMounted || w.VehKind != 'T' {
		t.Fatal("should mount the tank when close")
	}
	p.Dir = 0
	e1 := &Enemy{X: p.X + 3.5, Y: p.Y, Kind: TileEnemy, HP: 40, Score: 100}
	e2 := &Enemy{X: p.X + 4.4, Y: p.Y, Kind: TileEnemy, HP: 40, Score: 100}
	w.Enemies = append(w.Enemies, e1, e2)
	w.tankFire(g)
	if len(w.Projs) != 1 {
		t.Fatal("the tank should fire one shell")
	}
	if w.TankCd <= 0 {
		t.Fatal("the cannon needs a cooldown")
	}
	w.tankFire(g)
	if len(w.Projs) != 1 {
		t.Fatal("the cannon has a cooldown between shots")
	}
	hp1, hp2 := e1.HP, e2.HP
	for i := 0; i < 80 && e1.HP == hp1 && !e1.Dead; i++ {
		w.update(1.0/60.0, g)
	}
	if e1.HP == hp1 && !e1.Dead {
		t.Fatal("the shell should reach and hurt the enemy")
	}
	if e2.HP == hp2 && !e2.Dead {
		t.Fatal("the blast should splash nearby enemies")
	}
	p.X += 5
	w.toggleVehicle()
	if w.VehMounted {
		t.Fatal("dismount should work from the tank")
	}
	if math.Hypot(w.Tank.X-p.X, w.Tank.Y-p.Y) < 1.05 {
		t.Fatal("the tank should be dropped behind the player")
	}
}

func TestDragonFireVolleyOnDesert(t *testing.T) {
	g := &Game{}
	g.startLevel(7)
	g.state = statePlay
	w := g.world
	var d *Enemy
	for _, e := range w.Enemies {
		if e.Kind == TileDragon {
			d = e
		}
	}
	if d == nil {
		t.Fatal("the desert finale must host the dragon")
	}
	p := &w.Player
	p.X, p.Y = d.X-3, d.Y
	w.Enemies = []*Enemy{d}
	d.ShootT = 0.01
	w.update(1.0/60.0, g)
	hot := 0
	for _, pr := range w.Projs {
		if pr.Hot {
			hot++
		}
	}
	if hot == 0 {
		t.Fatal("the dragon should spit a volley of fire projectiles")
	}
}

func TestJungleLevelWithHeli(t *testing.T) {
	g := &Game{}
	g.startLevel(6)
	g.state = statePlay
	w := g.world
	if w.W < 30 || w.H < 20 {
		t.Fatalf("jungle should be a big map, got %dx%d", w.W, w.H)
	}
	if w.Heli.X == 0 && w.Heli.Y == 0 {
		t.Fatal("jungle must host the heli")
	}
	chests := 0
	for _, pk := range w.Pickups {
		if pk.Kind == TileChest {
			chests++
		}
	}
	if chests < 2 {
		t.Fatalf("jungle should hide chests, got %d", chests)
	}
	if findKind(w, TileExit) == nil {
		t.Fatal("jungle needs an exit")
	}
}

func TestHeliFreeMovementAndRockets(t *testing.T) {
	g := &Game{}
	g.startLevel(7)
	g.state = statePlay
	w := g.world
	p := &w.Player
	w.Heli = Vehicle{X: p.X + 1.0, Y: p.Y}
	w.toggleVehicle()
	if !w.VehMounted || w.VehKind != 'H' {
		t.Fatal("should mount the heli when close")
	}
	p.X, p.Y = 2.5, 9.5
	w.flyMove(1.2, 0)
	if !w.solidAt(int(p.X), int(p.Y)) {
		t.Fatal("the heli should fly over walls")
	}
	p.X, p.Y = 2.5, 1.5
	w.flyMove(100, 0)
	if p.X != float64(w.W)-1.0 {
		t.Fatal("the heli must stay inside the map")
	}
	p.Dir = 0
	p.X, p.Y = 10.5, 1.5
	e := &Enemy{X: p.X + 3.5, Y: p.Y, Kind: TileEnemy, HP: 40, Score: 100}
	w.Enemies = append(w.Enemies, e)
	w.heliFire(g)
	if len(w.Projs) == 0 {
		t.Fatal("the heli should fire a rocket")
	}
	hp := e.HP
	for i := 0; i < 120 && e.HP == hp && !e.Dead; i++ {
		w.update(1.0/60.0, g)
	}
	if e.HP == hp && !e.Dead {
		t.Fatal("the rocket should reach and hurt the enemy")
	}
}

func TestBurnFromDragonFire(t *testing.T) {
	g := &Game{}
	g.startLevel(7)
	g.state = statePlay
	w := g.world
	var d *Enemy
	for _, e := range w.Enemies {
		if e.Kind == TileDragon {
			d = e
		}
	}
	if d == nil {
		t.Fatal("the desert finale must host the dragon")
	}
	p := &w.Player
	p.X, p.Y = d.X-3, d.Y
	w.Enemies = []*Enemy{d}
	d.ShootT = 0.01
	w.update(1.0/60.0, g)
	ignites := false
	for _, pr := range w.Projs {
		if pr.Burns && pr.Hot {
			ignites = true
		}
	}
	if !ignites {
		t.Fatal("dragon fire should ignite the player")
	}
	w.Projs = append(w.Projs, &Proj{X: p.X - 0.5, Y: p.Y, DX: 1, DY: 0, Dmg: 5, Hot: true, Burns: true})
	hp0 := p.HP
	w.update(1.0/60.0, g)
	if p.BurnT <= 0 {
		t.Fatal("the player should catch fire")
	}
	for i := 0; i < 120; i++ {
		w.update(1.0/60.0, g)
	}
	if p.HP >= hp0 {
		t.Fatal("burn damage should tick over time")
	}
}

func TestSoundBuffersGenerated(t *testing.T) {
	bufs := buildSounds()
	if len(bufs) < 12 {
		t.Fatalf("expected many sounds, got %d", len(bufs))
	}
	for name, b := range bufs {
		if len(b) == 0 || len(b)%4 != 0 {
			t.Fatalf("sound %s must be a non-empty PCM int16 stereo buffer", name)
		}
	}
}
