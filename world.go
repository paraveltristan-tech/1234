package main

import (
	"math"
	"math/rand"
)

// ---- weapons ----

type Weapon struct {
	Name       string
	Dmg        int
	Pellets    int
	Spread     float64
	Rate       float64
	Auto       bool
	FlashMax   float64
	Recoil     float64
	AmmoPickup int
	MaxAmmo    int
}

var weapons = [3]Weapon{
	{Name: "PISTOLET", Dmg: 26, Pellets: 1, Spread: 0.05, Rate: 0.34, FlashMax: 0.08, Recoil: 14, AmmoPickup: 12, MaxAmmo: 99},
	{Name: "FUSIL A POMPE", Dmg: 11, Pellets: 6, Spread: 0.30, Rate: 0.95, FlashMax: 0.15, Recoil: 30, AmmoPickup: 4, MaxAmmo: 40},
	{Name: "FUSIL A PLASMA", Dmg: 14, Pellets: 1, Spread: 0.03, Rate: 0.11, Auto: true, FlashMax: 0.05, Recoil: 7, AmmoPickup: 16, MaxAmmo: 99},
}

// ---- enemies ----

type enemyKind struct {
	HP      int
	Speed   float64
	Damage  int
	AtkCD   float64
	Range   float64
	Aggro   float64
	Scale   float64
	Score   int
	Ranged  bool
	ShootCD float64
	Keep    float64
}

var enemyKinds = map[byte]enemyKind{
	TileEnemy:   {HP: 40, Speed: 1.65, Damage: 12, AtkCD: 0.9, Range: 0.95, Aggro: 12.5, Scale: 0.80, Score: 100},
	TileSpectre: {HP: 22, Speed: 2.7, Damage: 8, AtkCD: 0.6, Range: 0.85, Aggro: 14, Scale: 0.55, Score: 150},
	TileGolem:   {HP: 120, Speed: 0.85, Damage: 22, AtkCD: 1.1, Range: 1.05, Aggro: 11, Scale: 1.1, Score: 250},
	TileSpitter: {HP: 35, Speed: 1.15, Damage: 10, AtkCD: 0, Range: 0.8, Aggro: 13, Scale: 0.75, Score: 200, Ranged: true, ShootCD: 1.7, Keep: 5.5},
}

type Player struct {
	X, Y, Dir float64
	HP        int
	Ammo      [3]int
	Weapon    int
	LastFire  float64
	Bob       float64
}

type Enemy struct {
	X, Y    float64
	Kind    byte
	Heading float64
	Score   int
	HP      int
	AtkCD   float64
	Flash   float64
	WanderT float64
	Frame   float64
	ShootT  float64
	Dead    bool
	DeadT   float64
}

type Pickup struct {
	X, Y float64
	Kind byte
	Bob0 float64
}

type Proj struct {
	X, Y, DX, DY float64
}

type World struct {
	Grid      [][]byte
	W, H      int
	Level     int
	Player    Player
	Enemies   []*Enemy
	Pickups   []Pickup
	Exits     []Pickup
	Projs     []*Proj
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
	projSpeed    = 6.5
	maxHP        = 100
	maxShotDist  = 16.0
)

func loadLevel(level int, seed int64) *World {
	var rows []string
	switch level {
	case 2:
		rows = level2
	case 3:
		rows = level3
	case 4:
		rows = level4
	case 5:
		rows = level5
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
		Grid:  grid,
		W:     maxLen,
		H:     len(rows),
		Level: level,
		rng:   rand.New(rand.NewSource(seed)),
	}

	for y, row := range grid {
		for x, c := range row {
			fx, fy := float64(x)+0.5, float64(y)+0.5
			if isEnemy(c) {
				k, ok := enemyKinds[c]
				if !ok {
					k = enemyKinds[TileEnemy]
				}
				w.Enemies = append(w.Enemies, &Enemy{
					X: fx, Y: fy, Kind: c, Score: k.Score, HP: k.HP,
					Heading: 0, // set below by wander init
				})
				w.Enemies[len(w.Enemies)-1].Heading = w.rng.Float64() * 2 * math.Pi
			}
			switch c {
			case TileSpawn:
				w.Player = Player{X: fx, Y: fy, Dir: 0.5, HP: maxHP}
				w.Player.Ammo[0] = 24
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
	w.Player.Bob += dt * 8

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

	w.updateProjectiles(dt, g)
	w.pickupLoop(g)
	w.decayFlashes(dt)
}

func (w *World) updateEnemy(e *Enemy, dt float64, g *Game) {
	if e.Dead {
		e.DeadT += dt
		return
	}
	k, ok := enemyKinds[e.Kind]
	if !ok {
		k = enemyKinds[TileEnemy]
	}
	e.Flash = math.Max(0, e.Flash-dt)
	e.AtkCD = math.Max(0, e.AtkCD-dt)
	p := &w.Player
	dx, dy := p.X-e.X, p.Y-e.Y
	dist := math.Hypot(dx, dy)
	seen := dist < k.Aggro && losClear(w.Grid, w.W, w.H, e.X, e.Y, p.X, p.Y)

	switch {
	case k.Ranged && seen:
		switch {
		case dist > k.Keep+0.8:
			e.moveToward(dx, dy, dist, k.Speed, dt, w)
		case dist < k.Keep-1.6:
			e.moveToward(-dx, -dy, dist, k.Speed*0.8, dt, w)
		}
		e.Heading = math.Atan2(dy, dx)
		e.ShootT -= dt
		if e.ShootT <= 0 {
			e.ShootT = k.ShootCD
			w.Projs = append(w.Projs, &Proj{
				X: e.X + dx/dist*0.45, Y: e.Y + dy/dist*0.45,
				DX: dx / dist, DY: dy / dist,
			})
		}
	case seen:
		if dist > k.Range*0.7 {
			e.moveToward(dx, dy, dist, k.Speed, dt, w)
		}
		e.Heading = math.Atan2(dy, dx)
		if dist < k.Range && e.AtkCD <= 0 {
			e.AtkCD = k.AtkCD
			p.HP -= k.Damage
			w.FlashRed = 0.5
			if p.HP <= 0 {
				g.state = stateDead
			}
		}
	default:
		e.WanderT -= dt
		if e.WanderT <= 0 {
			e.WanderT = w.rng.Float64()*2 + 1
			e.Heading = w.rng.Float64() * 2 * math.Pi
		}
		nx, ny := tryMove(w.Grid, w.W, w.H, e.X, e.Y, math.Cos(e.Heading)*k.Speed*0.4*dt, math.Sin(e.Heading)*k.Speed*0.4*dt, enemyRadius)
		if nx == e.X && ny == e.Y {
			e.Heading = w.rng.Float64() * 2 * math.Pi
		}
		e.X, e.Y = nx, ny
	}
	e.Frame += dt
}

func (e *Enemy) moveToward(dx, dy, dist, speed, dt float64, w *World) {
	e.X, e.Y = tryMove(w.Grid, w.W, w.H, e.X, e.Y, dx/dist*speed*dt, dy/dist*speed*dt, enemyRadius)
}

func (w *World) updateProjectiles(dt float64, g *Game) {
	kept := w.Projs[:0]
	p := &w.Player
	for _, pr := range w.Projs {
		pr.X += pr.DX * projSpeed * dt
		pr.Y += pr.DY * projSpeed * dt
		if w.solidAt(int(pr.X), int(pr.Y)) {
			continue
		}
		if math.Hypot(pr.X-p.X, pr.Y-p.Y) < 0.45 {
			p.HP -= 10
			w.FlashRed = 0.5
			if p.HP <= 0 {
				g.state = stateDead
			}
			continue
		}
		kept = append(kept, pr)
	}
	w.Projs = kept
}

func (w *World) pickupLoop(g *Game) {
	p := &w.Player
	for i := 0; i < len(w.Pickups); i++ {
		pk := w.Pickups[i]
		if math.Hypot(pk.X-p.X, pk.Y-p.Y) < 0.55 {
			switch pk.Kind {
			case TileHealth:
				if p.HP < maxHP {
					p.HP = min(maxHP, p.HP+25)
					w.FlashGood = 0.4
					g.score += 10
				}
			case TileAmmo:
				picked := false
				for wi := range weapons {
					want := weapons[wi].AmmoPickup
					if p.Ammo[wi] < weapons[wi].MaxAmmo {
						p.Ammo[wi] = min(weapons[wi].MaxAmmo, p.Ammo[wi]+want)
						picked = true
					}
				}
				if picked {
					w.FlashGood = 0.4
					g.score += 10
				}
			}
			w.Pickups = append(w.Pickups[:i], w.Pickups[i+1:]...)
			i--
		}
	}
	for _, ex := range w.Exits {
		if math.Hypot(ex.X-p.X, ex.Y-p.Y) < 0.6 {
			w.Cleared = true
			g.score += 250
		}
	}
}

func (w *World) decayFlashes(dt float64) {
	w.FlashRed = math.Max(0, w.FlashRed-dt)
	w.FlashGood = math.Max(0, w.FlashGood-dt)
	w.Muzzle = math.Max(0, w.Muzzle-dt)
}

// fire discharges the equipped weapon and returns true if anything was hit.
func (w *World) fire(g *Game) bool {
	p := &w.Player
	wp := weapons[p.Weapon]
	if p.Ammo[p.Weapon] <= 0 || w.T-p.LastFire < wp.Rate || w.Cleared {
		return false
	}
	p.Ammo[p.Weapon]--
	p.LastFire = w.T
	w.Muzzle = wp.FlashMax

	anyHit := false
	for i := 0; i < wp.Pellets; i++ {
		ang := p.Dir + (w.rng.Float64()-0.5)*wp.Spread
		if w.hitscan(ang, wp.Dmg, g) {
			anyHit = true
		}
	}
	return anyHit
}

func (w *World) hitscan(ang float64, dmg int, g *Game) bool {
	p := &w.Player
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
		if math.Abs(angDiff(ang, math.Atan2(dy, dx))) > 0.06+0.3/dist {
			continue
		}
		if !losClear(w.Grid, w.W, w.H, p.X, p.Y, e.X, e.Y) {
			continue
		}
		hit = e
		hitDist = dist
	}
	if hit == nil {
		return false
	}
	hit.HP -= dmg + w.rng.Intn(6)
	hit.Flash = 0.14
	if hit.HP <= 0 {
		hit.Dead = true
		g.kills++
		g.score += hit.Score
	}
	return true
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
