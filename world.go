package main

import (
	"math"
	"math/rand"
)

type Player struct {
	X, Y, Dir float64
	HP        int
	Ammo      int
	LastFire  float64
	Bob       float64
}

type Enemy struct {
	X, Y    float64
	Heading float64
	HP      int
	AtkCD   float64
	Flash   float64
	WanderT float64
	Frame   float64
	Dead    bool
	DeadT   float64
}

type Pickup struct {
	X, Y float64
	Kind byte
	Bob0 float64
}

type World struct {
	Grid      [][]byte
	W, H      int
	Player    Player
	Enemies   []*Enemy
	Pickups   []Pickup
	Exits     []Pickup
	T         float64
	FlashRed  float64
	FlashGood float64
	Muzzle    float64
	Kills     int
	Cleared   bool
	rng       *rand.Rand
}

const (
	playerRadius = 0.22
	enemyRadius  = 0.30
	moveSpeed    = 3.1
	sprintMult   = 1.55
	turnSpeed    = 2.7
	mouseSens    = 0.0038
	fireRate     = 0.34
	enemySpeed   = 1.65
	enemyDamage  = 12
	enemyAtkCD   = 0.9
	enemyRange   = 0.95
	enemyAggroD  = 12.5
	pistolDamage = 26
	pistolSpread = 0.055
	maxShotDist  = 16.0
	maxHP        = 100
	maxAmmo      = 99
)

func loadLevel(level int, seed int64) *World {
	var rows []string
	switch level {
	case 2:
		rows = level2
	case 3:
		rows = level3
	default:
		rows = level1
	}

	maxLen := 0
	for _, r := range rows {
		if len(r) > maxLen {
			maxLen = len(r)
		}
	}
	grid := make([][]byte, len(rows))
	for i, r := range rows {
		row := make([]byte, maxLen)
		copy(row, []byte(r))
		for j := len(r); j < maxLen; j++ {
			row[j] = TileBrick
		}
		grid[i] = row
	}

	w := &World{
		Grid: grid,
		W:    maxLen,
		H:    len(rows),
		rng:  rand.New(rand.NewSource(seed)),
	}

	for y, row := range grid {
		for x, c := range row {
			fx, fy := float64(x)+0.5, float64(y)+0.5
			switch c {
			case TileSpawn:
				w.Player = Player{X: fx, Y: fy, Dir: 0.5, HP: maxHP, Ammo: 24}
			case TileEnemy:
				w.Enemies = append(w.Enemies, &Enemy{X: fx, Y: fy, HP: 40, Heading: w.rng.Float64() * 2 * math.Pi})
			case TileHealth:
				w.Pickups = append(w.Pickups, Pickup{X: fx, Y: fy, Kind: TileHealth, Bob0: w.rng.Float64() * 6})
			case TileAmmo:
				w.Pickups = append(w.Pickups, Pickup{X: fx, Y: fy, Kind: TileAmmo, Bob0: w.rng.Float64() * 6})
			case TileExit:
				w.Exits = append(w.Exits, Pickup{X: fx, Y: fy, Kind: TileExit})
			}
		}
	}
	return w
}

// TileAt returns the tile at grid coordinates, or a wall outside bounds.
func (w *World) TileAt(x, y int) byte {
	if x < 0 || x >= w.W || y < 0 || y >= w.H {
		return TileBrick
	}
	return w.Grid[y][x]
}

func (w *World) solidAt(x, y int) bool {
	return isSolid(w.TileAt(x, y))
}

func (w *World) update(dt float64, g *Game) {
	if w.Cleared {
		return
	}
	w.T += dt
	p := &w.Player
	p.Bob += dt * 8

	for _, e := range w.Enemies {
		w.updateEnemy(e, dt, g)
	}
	alive := w.Enemies[:0]
	for _, e := range w.Enemies {
		if !(e.Dead && e.DeadT > 0.8) {
			alive = append(alive, e)
		}
	}
	w.Enemies = alive

	w.pickupLoop(g)
	w.decayFlashes(dt)
}

func (w *World) updateEnemy(e *Enemy, dt float64, g *Game) {
	if e.Dead {
		e.DeadT += dt
		return
	}
	e.Flash = maxF(0, e.Flash-dt)
	e.AtkCD = maxF(0, e.AtkCD-dt)
	p := &w.Player
	dx, dy := p.X-e.X, p.Y-e.Y
	dist := math.Hypot(dx, dy)

	if dist < enemyAggroD && losClear(w.Grid, w.W, w.H, e.X, e.Y, p.X, p.Y) {
		if dist > enemyRange*0.7 {
			nx, ny := tryMove(w.Grid, w.W, w.H, e.X, e.Y, dx/dist*enemySpeed*dt, dy/dist*enemySpeed*dt, enemyRadius)
			e.X, e.Y = nx, ny
		}
		e.Heading = math.Atan2(dy, dx)
		if dist < enemyRange && e.AtkCD <= 0 {
			e.AtkCD = enemyAtkCD
			p.HP -= enemyDamage
			w.FlashRed = 0.5
			if p.HP <= 0 {
				g.state = stateDead
			}
		}
	} else {
		e.WanderT -= dt
		if e.WanderT <= 0 {
			e.WanderT = w.rng.Float64()*2 + 1
			e.Heading = w.rng.Float64() * 2 * math.Pi
		}
		nx, ny := tryMove(w.Grid, w.W, w.H, e.X, e.Y, math.Cos(e.Heading)*enemySpeed*0.4*dt, math.Sin(e.Heading)*enemySpeed*0.4*dt, enemyRadius)
		if nx == e.X && ny == e.Y {
			e.Heading = w.rng.Float64() * 2 * math.Pi
		}
		e.X, e.Y = nx, ny
	}
	e.Frame += dt
}

func (w *World) pickupLoop(g *Game) {
	p := &w.Player
	for i := 0; i < len(w.Pickups); i++ {
		pk := w.Pickups[i]
		if math.Hypot(pk.X-p.X, pk.Y-p.Y) < 0.55 {
			switch pk.Kind {
			case TileHealth:
				p.HP = min(maxHP, p.HP+25)
				w.FlashGood = 0.4
				g.score += 10
			case TileAmmo:
				p.Ammo = min(maxAmmo, p.Ammo+12)
				w.FlashGood = 0.4
				g.score += 10
			}
			w.Pickups = append(w.Pickups[:i], w.Pickups[i+1:]...)
			i--
		}
	}
	for _, ex := range w.Exits {
		if math.Hypot(ex.X-p.X, ex.Y-p.Y) < 0.6 {
			w.Cleared = true
			g.score += 250 + g.kills*100
		}
	}
}

func (w *World) decayFlashes(dt float64) {
	w.FlashRed = maxF(0, w.FlashRed-dt)
	w.FlashGood = maxF(0, w.FlashGood-dt)
	w.Muzzle = maxF(0, w.Muzzle-dt)
}

// fire performs a hitscan shot toward where the player looks.
func (w *World) fire(g *Game) bool {
	p := &w.Player
	if p.Ammo <= 0 || w.T-p.LastFire < fireRate || w.Cleared {
		return false
	}
	p.Ammo--
	p.LastFire = w.T
	w.Muzzle = 0.08

	var hit *Enemy
	hitDist := math.MaxFloat64
	for _, e := range w.Enemies {
		if e.Dead {
			continue
		}
		dx, dy := e.X-p.X, e.Y-p.Y
		dist := math.Hypot(dx, dy)
		if dist > maxShotDist || dist > hitDist {
			continue
		}
		spreadTol := pistolSpread + 0.35/dist
		if math.Abs(angDiff(p.Dir, math.Atan2(dy, dx))) > spreadTol {
			continue
		}
		if !losClear(w.Grid, w.W, w.H, p.X, p.Y, e.X, e.Y) {
			continue
		}
		hit = e
		hitDist = dist
	}
	if hit != nil {
		hit.HP -= pistolDamage + w.rng.Intn(11)
		hit.Flash = 0.14
		if hit.HP <= 0 {
			hit.Dead = true
			g.kills++
			g.score += 100
		}
		return true
	}
	return false
}

func maxF(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
