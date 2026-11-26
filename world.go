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
	Splash     float64
	Range      float64
	MagSize    int
	ReloadTime float64
}

var weapons = [5]Weapon{
	{Name: "PISTOLET", Dmg: 26, Pellets: 1, Spread: 0.05, Rate: 0.34, FlashMax: 0.08, Recoil: 14, AmmoPickup: 12, MaxAmmo: 99, MagSize: 12, ReloadTime: 0.9},
	{Name: "FUSIL A POMPE", Dmg: 11, Pellets: 6, Spread: 0.30, Rate: 0.95, FlashMax: 0.15, Recoil: 30, AmmoPickup: 4, MaxAmmo: 40, MagSize: 8, ReloadTime: 1.4},
	{Name: "FUSIL A PLASMA", Dmg: 14, Pellets: 1, Spread: 0.03, Rate: 0.11, Auto: true, FlashMax: 0.05, Recoil: 7, AmmoPickup: 16, MaxAmmo: 99, MagSize: 40, ReloadTime: 1.2},
	{Name: "SNIPER", Dmg: 95, Pellets: 1, Spread: 0.0, Rate: 1.5, FlashMax: 0.18, Recoil: 34, AmmoPickup: 6, MaxAmmo: 30, Range: 30, MagSize: 5, ReloadTime: 1.8},
	{Name: "CANON D'IONS", Dmg: 70, Pellets: 1, Spread: 0.0, Rate: 1.35, FlashMax: 0.22, Recoil: 40, AmmoPickup: 3, MaxAmmo: 15, Splash: 2.6, MagSize: 4, ReloadTime: 2.2},
}

func wpRange(wp Weapon) float64 {
	if wp.Range > 0 {
		return wp.Range
	}
	return maxShotDist
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
	Boss    bool
	Volley  int
	Reward  int
}

var enemyKinds = map[byte]enemyKind{
	TileEnemy:    {HP: 40, Speed: 1.65, Damage: 12, AtkCD: 0.9, Range: 0.95, Aggro: 12.5, Scale: 0.80, Score: 100, Reward: 60},
	TileSpectre:  {HP: 22, Speed: 2.7, Damage: 8, AtkCD: 0.6, Range: 0.85, Aggro: 14, Scale: 0.55, Score: 150, Reward: 80},
	TileGolem:    {HP: 120, Speed: 0.85, Damage: 22, AtkCD: 1.1, Range: 1.05, Aggro: 11, Scale: 1.1, Score: 250, Reward: 150},
	TileSpitter:  {HP: 35, Speed: 1.15, Damage: 10, AtkCD: 0, Range: 0.8, Aggro: 13, Scale: 0.75, Score: 200, Ranged: true, ShootCD: 1.7, Keep: 5.5, Reward: 100},
	TileBoss:     {HP: 900, Speed: 1.05, Damage: 26, AtkCD: 1.2, Range: 1.35, Aggro: 30, Scale: 1.9, Score: 2000, ShootCD: 2.0, Boss: true, Volley: 3, Reward: 300},
	TileColossus: {HP: 1400, Speed: 1.35, Damage: 30, AtkCD: 1.0, Range: 1.45, Aggro: 30, Scale: 2.1, Score: 5000, ShootCD: 1.7, Boss: true, Volley: 5, Reward: 500},
	TileDragon:   {HP: 1100, Speed: 1.6, Damage: 32, AtkCD: 0.9, Range: 1.4, Aggro: 30, Scale: 2.3, Score: 7000, ShootCD: 1.5, Boss: true, Volley: 4, Reward: 700},
}

type Player struct {
	X, Y, Dir float64
	HP        int
	Ammo      [5]int
	Reserve   [5]int
	Skins     [5]int
	Weapon    int
	LastFire  float64
	Bob       float64
	ReloadT   float64
	ReloadDur float64
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
	X, Y   float64
	Kind   byte
	Bob0   float64
	Opened bool
	OpenT  float64
}

type Vehicle struct {
	X, Y float64
}

type Proj struct {
	X, Y, DX, DY float64
	Dmg          int
	Hot          bool
	Tank         bool
}

type World struct {
	Grid       [][]byte
	W, H       int
	Level      int
	Player     Player
	Enemies    []*Enemy
	Pickups    []Pickup
	Exits      []Pickup
	Projs      []*Proj
	Parts      []*Particle
	T          float64
	FlashRed   float64
	FlashGood  float64
	Muzzle     float64
	Kills      int
	Cleared    bool
	Survive    bool
	CS         bool
	cash       int
	Wave       int
	ToSpawn    int
	SpawnT     float64
	RestT      float64
	BossWave   bool
	Boss2Wave  bool
	Boss3Wave  bool
	Vehicle    Vehicle
	Tank       Vehicle
	VehMounted bool
	VehKind    byte
	VehCd      float64
	TankCd     float64
	rng        *rand.Rand
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
	case levelDesert:
		rows = levelDesertMap
	case levelSurvive:
		rows = levelSurv
	case levelCS:
		rows = levelCSArena
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
				w.Player.Ammo[0] = weapons[0].MagSize
				w.Player.Reserve[0] = 24 - weapons[0].MagSize
			case TileHealth:
				w.Pickups = append(w.Pickups, Pickup{X: fx, Y: fy, Kind: TileHealth, Bob0: w.rng.Float64() * 6})
			case TileAmmo:
				w.Pickups = append(w.Pickups, Pickup{X: fx, Y: fy, Kind: TileAmmo, Bob0: w.rng.Float64() * 6})
			case TileChest:
				w.Pickups = append(w.Pickups, Pickup{X: fx, Y: fy, Kind: TileChest, Bob0: w.rng.Float64() * 6})
			case TileVehicle:
				w.Vehicle = Vehicle{X: fx, Y: fy}
			case TileTank:
				w.Tank = Vehicle{X: fx, Y: fy}
			case TileExit:
				w.Exits = append(w.Exits, Pickup{X: fx, Y: fy, Kind: TileExit})
			}
		}
	}
	if level == levelSurvive || level == levelCS {
		w.Survive = true
		if level == levelCS {
			w.CS = true
			w.cash = 800
		}
		w.startWave(1)
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
	w.tickReload(dt)

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

	w.vehRampage(dt, g)
	w.updateProjectiles(dt, g)
	w.updateParticles(dt)
	w.pickupLoop(g)
	if w.Survive {
		w.updateSurvive(dt, g)
	}
	w.decayFlashes(dt)
}

func (w *World) aliveCount() int {
	n := 0
	for _, e := range w.Enemies {
		if !e.Dead {
			n++
		}
	}
	return n
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
				DX: dx / dist, DY: dy / dist, Dmg: 10,
			})
		}
	case k.Boss && seen:
		e.Heading = math.Atan2(dy, dx)
		if dist < k.Range {
			if dist > k.Range*0.7 {
				e.moveToward(dx, dy, dist, k.Speed, dt, w)
			}
			if e.AtkCD <= 0 {
				e.AtkCD = k.AtkCD
				w.meleePlayer(k, g)
				w.spawnBurst(p.X, p.Y, 0.6, 1.8, 6, bossCol, false, 0.4, 0.05)
			}
		} else {
			if dist > 2.0 {
				e.moveToward(dx, dy, dist, k.Speed, dt, w)
			}
			e.ShootT -= dt
			if e.ShootT <= 0 && dist < 11 {
				e.ShootT = k.ShootCD
				volley := k.Volley
				if volley <= 0 {
					volley = 3
				}
				dmg := 14
				switch e.Kind {
				case TileColossus:
					dmg = 16
				case TileDragon:
					dmg = 18
				}
				for i := 0; i < volley; i++ {
					ang := e.Heading + (float64(i)-float64(volley-1)/2)*0.35
					w.Projs = append(w.Projs, &Proj{
						X: e.X + math.Cos(ang)*0.55, Y: e.Y + math.Sin(ang)*0.55,
						DX: math.Cos(ang), DY: math.Sin(ang), Dmg: dmg,
						Hot: e.Kind == TileColossus || e.Kind == TileDragon,
					})
				}
			}
		}
	case seen:
		if dist > k.Range*0.7 {
			e.moveToward(dx, dy, dist, k.Speed, dt, w)
		}
		e.Heading = math.Atan2(dy, dx)
		if dist < k.Range && e.AtkCD <= 0 {
			e.AtkCD = k.AtkCD
			w.meleePlayer(k, g)
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

func (w *World) meleePlayer(k enemyKind, g *Game) {
	p := &w.Player
	p.HP -= k.Damage
	w.FlashRed = 0.5
	if p.HP <= 0 {
		g.state = stateDead
	}
}

func (w *World) boss() *Enemy {
	for _, e := range w.Enemies {
		if e.Kind == TileBoss && !e.Dead {
			return e
		}
	}
	return nil
}

func (w *World) bossAlive() bool { return w.boss() != nil }

func (w *World) colossus() *Enemy {
	for _, e := range w.Enemies {
		if e.Kind == TileColossus && !e.Dead {
			return e
		}
	}
	return nil
}

// startReload demarre la recharge de l'arme equipee : chargeur non plein,
// reserve disponible et aucun rechargement en cours.
func (w *World) startReload() {
	p := &w.Player
	wp := weapons[p.Weapon]
	if p.ReloadT > 0 || p.Ammo[p.Weapon] >= wp.MagSize || p.Reserve[p.Weapon] <= 0 || w.Cleared {
		return
	}
	p.ReloadDur = wp.ReloadTime
	p.ReloadT = wp.ReloadTime
}

// tickReload fait avancer la recharge puis transfere les balles du chargeur.
func (w *World) tickReload(dt float64) {
	p := &w.Player
	if p.ReloadT <= 0 {
		return
	}
	p.ReloadT -= dt
	if p.ReloadT > 0 {
		return
	}
	p.ReloadT = 0
	wp := weapons[p.Weapon]
	take := min(wp.MagSize-p.Ammo[p.Weapon], p.Reserve[p.Weapon])
	p.Ammo[p.Weapon] += take
	p.Reserve[p.Weapon] -= take
	w.spawnBurst(p.X, p.Y, 0.7, 1.2, 5, cashCol, true, 0.5, 0.04)
}

func (w *World) updateProjectiles(dt float64, g *Game) {
	kept := w.Projs[:0]
	p := &w.Player
	for _, pr := range w.Projs {
		bcol := acidCol
		if pr.Hot {
			bcol = ionCol
		}
		pr.X += pr.DX * projSpeed * dt
		pr.Y += pr.DY * projSpeed * dt
		if pr.Tank {
			expl := w.solidAt(int(pr.X), int(pr.Y))
			if !expl {
				for _, e := range w.Enemies {
					if !e.Dead && math.Hypot(e.X-pr.X, e.Y-pr.Y) < 0.65 {
						expl = true
						break
					}
				}
			}
			if expl {
				w.tankBlast(g, pr.X, pr.Y)
			} else {
				kept = append(kept, pr)
			}
			continue
		}
		if w.solidAt(int(pr.X), int(pr.Y)) {
			w.spawnBurst(pr.X-pr.DX*0.2, pr.Y-pr.DY*0.2, 0.6, 1.4, 7, bcol, true, 0.4, 0.045)
			continue
		}
		if math.Hypot(pr.X-p.X, pr.Y-p.Y) < 0.45 {
			p.HP -= pr.Dmg
			w.FlashRed = 0.5
			w.spawnBurst(pr.X, pr.Y, 0.6, 1.4, 8, bcol, true, 0.4, 0.045)
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
	for i := range w.Pickups {
		pk := &w.Pickups[i]
		if pk.Kind == TileChest && !pk.Opened && math.Hypot(pk.X-p.X, pk.Y-p.Y) < 0.8 {
			pk.Opened = true
			pk.OpenT = w.T
			w.chestLoot(g)
		}
	}
	for i := 0; i < len(w.Pickups); i++ {
		pk := w.Pickups[i]
		if pk.Kind == TileChest {
			continue
		}
		if math.Hypot(pk.X-p.X, pk.Y-p.Y) < 0.55 {
			switch pk.Kind {
			case TileHealth:
				if p.HP < maxHP {
					p.HP = min(maxHP, p.HP+25)
					w.FlashGood = 0.4
					g.score += 10
					w.spawnPickupGlow(pk.X, pk.Y, healCol)
				}
			case TileAmmo:
				picked := false
				for wi := range weapons {
					want := weapons[wi].AmmoPickup
					room := weapons[wi].MaxAmmo - p.Ammo[wi] - p.Reserve[wi]
					if room > 0 {
						give := min(want, room)
						fill := min(weapons[wi].MagSize-p.Ammo[wi], give)
						p.Ammo[wi] += fill
						p.Reserve[wi] += give - fill
						picked = true
					}
				}
				if picked {
					w.FlashGood = 0.4
					g.score += 10
					w.spawnPickupGlow(pk.X, pk.Y, ammoCol)
				}
			}
			w.Pickups = append(w.Pickups[:i], w.Pickups[i+1:]...)
			i--
		}
	}
	if w.bossAlive() {
		return
	}
	for _, ex := range w.Exits {
		if math.Hypot(ex.X-p.X, ex.Y-p.Y) < 0.6 {
			w.Cleared = true
			g.score += 250
		}
	}
}

// chestLoot distribue une recompense aleatoire quand un coffre s'ouvre.
func (w *World) chestLoot(g *Game) {
	p := &w.Player
	switch w.rng.Intn(4) {
	case 0:
		p.HP = min(maxHP, p.HP+35)
		w.spawnPickupGlow(p.X, p.Y, healCol)
	case 1:
		for wi := range weapons {
			room := weapons[wi].MaxAmmo - p.Ammo[wi] - p.Reserve[wi]
			if room > 0 {
				give := min(weapons[wi].AmmoPickup, room)
				fill := min(weapons[wi].MagSize-p.Ammo[wi], give)
				p.Ammo[wi] += fill
				p.Reserve[wi] += give - fill
			}
		}
		w.spawnPickupGlow(p.X, p.Y, ammoCol)
	case 2:
		if w.CS {
			w.cash += 150 + w.rng.Intn(150)
		} else {
			g.score += 200
		}
		w.spawnPickupGlow(p.X, p.Y, cashCol)
	default:
		g.score += 100
		p.HP = min(maxHP, p.HP+15)
		w.spawnPickupGlow(p.X, p.Y, healCol)
	}
	w.FlashGood = 0.5
}

// toggleVehicle monte ou descend du vehicule pilote (touche E) : le plus proche.
func (w *World) toggleVehicle() {
	p := &w.Player
	if w.VehMounted {
		nx := w.Player.X - math.Cos(w.Player.Dir)*1.1
		ny := w.Player.Y - math.Sin(w.Player.Dir)*1.1
		if !w.solidAt(int(nx), int(ny)) {
			if w.VehKind == 'T' {
				w.Tank.X, w.Tank.Y = nx, ny
			} else {
				w.Vehicle.X, w.Vehicle.Y = nx, ny
			}
		}
		w.VehMounted = false
		w.VehKind = 0
		w.FlashGood = 0.3
		return
	}
	dvtt := math.Hypot(w.Vehicle.X-p.X, w.Vehicle.Y-p.Y)
	dtank := math.Hypot(w.Tank.X-p.X, w.Tank.Y-p.Y)
	vttOK := w.Vehicle.X > 0 && dvtt < 1.3
	tankOK := w.Tank.X > 0 && dtank < 1.3
	if tankOK && (!vttOK || dtank <= dvtt) {
		w.VehMounted, w.VehKind = true, 'T'
		w.FlashGood = 0.4
	} else if vttOK {
		w.VehMounted, w.VehKind = true, 'V'
		w.FlashGood = 0.4
	}
}

// vehRampage ecrase les ennemis percutes par le vehicule pilote.
func (w *World) vehRampage(dt float64, g *Game) {
	w.VehCd = math.Max(0, w.VehCd-dt)
	if !w.VehMounted {
		return
	}
	dmg, rad, cd := 75, 0.95, 0.35
	if w.VehKind == 'T' {
		dmg, rad, cd = 140, 1.15, 0.55
	}
	p := &w.Player
	for _, e := range w.Enemies {
		if !e.Dead && w.VehCd <= 0 && math.Hypot(e.X-p.X, e.Y-p.Y) < rad {
			w.VehCd = cd
			w.damageEnemy(e, dmg, g)
		}
	}
}

// tankFire tire un obus explosif du canon du tank (clic gauche a bord).
func (w *World) tankFire(g *Game) {
	if w.TankCd > 0 || !w.VehMounted || w.VehKind != 'T' {
		return
	}
	w.TankCd = 0.9
	w.Muzzle = 0.5
	p := &w.Player
	ang := p.Dir
	w.Projs = append(w.Projs, &Proj{
		X:    p.X + math.Cos(ang)*0.7,
		Y:    p.Y + math.Sin(ang)*0.7,
		DX:   math.Cos(ang),
		DY:   math.Sin(ang),
		Tank: true,
	})
	w.spawnBurst(p.X+math.Cos(ang)*1.1, p.Y+math.Sin(ang)*1.1, 0.5, 1.6, 8, fireCol, true, 0.35, 0.05)
}

// tankBlast fait exploser l'obus du tank : degats de zone feu autour de l'impact.
func (w *World) tankBlast(g *Game, cx, cy float64) {
	w.spawnBurst(cx, cy, 0.6, 3.5, 26, fireCol, true, 0.6, 0.075)
	w.spawnBurst(cx, cy, 0.5, 2.2, 10, ionCol, true, 0.5, 0.060)
	w.FlashGood = 0.25
	for _, e := range w.Enemies {
		if e.Dead {
			continue
		}
		d := math.Hypot(e.X-cx, e.Y-cy)
		if d > 2.4 {
			continue
		}
		dmg := int(120*(1-0.6*d/2.4)) + 1
		w.damageEnemy(e, dmg, g)
	}
}

func (w *World) dragon() *Enemy {
	for _, e := range w.Enemies {
		if e.Kind == TileDragon && !e.Dead {
			return e
		}
	}
	return nil
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
	if w.T-p.LastFire < wp.Rate || w.Cleared {
		return false
	}
	if p.ReloadT > 0 || w.VehMounted {
		return false
	}
	if p.Ammo[p.Weapon] <= 0 {
		w.startReload()
		return false
	}
	p.Ammo[p.Weapon]--
	p.LastFire = w.T
	w.Muzzle = wp.FlashMax

	if wp.Splash > 0 {
		ang := p.Dir
		hit, hx, hy := w.traceShot(w.Player.X, w.Player.Y, ang, wpRange(wp), g)
		w.ionBlast(g, hit, hx, hy, wp)
		return hit != nil
	}

	anyHit := false
	for i := 0; i < wp.Pellets; i++ {
		ang := p.Dir + (w.rng.Float64()-0.5)*wp.Spread
		hit, hx, hy := w.traceShot(w.Player.X, w.Player.Y, ang, wpRange(wp), g)
		if hit != nil {
			anyHit = true
			w.damageEnemy(hit, wp.Dmg, g)
		} else {
			w.spawnSparks(hx, hy, 4)
		}
		if wp.Range > 0 {
			w.spawnTracer(w.Player.X, w.Player.Y, p.Dir, hx, hy)
		}
	}
	return anyHit
}

// spawnTracer dessine un traceur lumineux le long du tir de sniper.
func (w *World) spawnTracer(sx, sy, ang, ex, ey float64) {
	dx, dy := ex-sx, ey-sy
	d := math.Hypot(dx, dy)
	steps := int(d / 0.7)
	for i := 1; i < steps; i++ {
		t := float64(i) / float64(steps)
		w.spawnBurst(sx+dx*t, sy+dy*t, 0.5, 0.02, 1, plasmaCol, true, 0.22, 0.030)
	}
}

// traceShot suit un rayon depuis le joueur et renvoie l'ennemi touche
// (son centre sert de point d'impact) ou le mur avec son point d'impact.
func (w *World) traceShot(sx, sy, ang, maxDist float64, g *Game) (*Enemy, float64, float64) {
	ca, sa := math.Cos(ang), math.Sin(ang)
	var hit *Enemy
	hitDist := maxDist
	for _, e := range w.Enemies {
		if e.Dead {
			continue
		}
		dx, dy := e.X-sx, e.Y-sy
		dist := math.Hypot(dx, dy)
		if dist > maxDist || dist > hitDist {
			continue
		}
		if math.Abs(angDiff(ang, math.Atan2(dy, dx))) > 0.06+0.3/dist {
			continue
		}
		if !losClear(w.Grid, w.W, w.H, sx, sy, e.X, e.Y) {
			continue
		}
		hit = e
		hitDist = dist
	}
	if hit != nil {
		return hit, hit.X, hit.Y
	}
	dist, _, _, _ := castStep(w.Grid, w.W, w.H, sx, sy, ca, sa)
	if dist > 1e8 {
		dist = maxDist
	}
	d := math.Max(dist-0.08, 0.2)
	return nil, sx + ca*d, sy + sa*d
}

func (w *World) damageEnemy(e *Enemy, baseDmg int, g *Game) {
	if e.Dead {
		return
	}
	e.HP -= baseDmg + w.rng.Intn(6)
	e.Flash = 0.14
	w.spawnBlood(e.X, e.Y, e.Kind, 5)
	if e.HP <= 0 {
		e.Dead = true
		w.spawnBlood(e.X, e.Y, e.Kind, 18)
		g.kills++
		g.score += e.Score
		if w.CS {
			reward, ok := enemyKinds[e.Kind]
			if !ok {
				reward = enemyKinds[TileEnemy]
			}
			w.cash += reward.Reward
			w.spawnBurst(e.X, e.Y, 0.7, 1.2, 6, cashCol, true, 0.5, 0.05)
		}
	}
}

// ionBlast applique les degats de zone du canon d'ions autour du point
// d'impact, avec decroissance jusqu'a 40% du degats de base en bord de zone.
func (w *World) ionBlast(g *Game, hit *Enemy, cx, cy float64, wp Weapon) {
	w.spawnBurst(cx, cy, 0.6, 3.2, 22, ionCol, true, 0.6, 0.070)
	w.FlashGood = 0.25
	for _, e := range w.Enemies {
		if e.Dead {
			continue
		}
		d := math.Hypot(e.X-cx, e.Y-cy)
		if d > wp.Splash {
			continue
		}
		dmg := int(float64(wp.Dmg)*(1-0.6*d/wp.Splash)) + 1
		if e == hit {
			dmg += wp.Dmg + w.rng.Intn(6)
			w.spawnBlood(e.X, e.Y, e.Kind, 8)
		} else {
			w.spawnBlood(e.X, e.Y, e.Kind, 5)
		}
		e.HP -= dmg
		e.Flash = 0.14
		if e.HP <= 0 && !e.Dead {
			e.Dead = true
			w.spawnBlood(e.X, e.Y, e.Kind, 18)
			g.kills++
			g.score += e.Score
		}
	}
}

// ---- mode survie ----

func (w *World) startWave(n int) {
	w.Wave = n
	w.ToSpawn = 4 + 2*n
	w.SpawnT = 1.0
	w.RestT = 3.0
	w.BossWave = !w.CS && n%5 == 0
	w.Boss2Wave = n%10 == 0 || (w.CS && n%6 == 0)
	w.Boss3Wave = n%12 == 0 || (w.CS && n%8 == 0)
	if n > 1 {
		w.spawnWavePickups()
	}
}

func (w *World) updateSurvive(dt float64, g *Game) {
	if w.ToSpawn > 0 || w.BossWave || w.Boss2Wave || w.Boss3Wave {
		w.SpawnT -= dt
		if w.SpawnT <= 0 {
			w.SpawnT = 0.7
			if w.BossWave && w.boss() == nil && w.spawnKind(TileBoss) {
				w.BossWave = false
			}
			if w.Boss2Wave && w.colossus() == nil && w.spawnKind(TileColossus) {
				w.Boss2Wave = false
			}
			if w.Boss3Wave && w.dragon() == nil && w.spawnKind(TileDragon) {
				w.Boss3Wave = false
			}
			if w.ToSpawn > 0 && w.spawnEnemyForWave() {
				w.ToSpawn--
			}
		}
		return
	}
	if w.aliveCount() == 0 {
		w.RestT -= dt
		if w.RestT <= 0 {
			g.score += 100 * w.Wave
			if w.CS {
				w.cash += 300
			}
			w.Player.HP = min(maxHP, w.Player.HP+15)
			w.startWave(w.Wave + 1)
		}
	}
}

func (w *World) spawnEnemyForWave() bool {
	pool := []byte{TileEnemy}
	if w.Wave >= 2 {
		pool = append(pool, TileSpectre)
	}
	if w.Wave >= 3 {
		pool = append(pool, TileSpitter)
	}
	if w.Wave >= 4 {
		pool = append(pool, TileGolem, TileGolem)
	}
	return w.spawnKind(pool[w.rng.Intn(len(pool))])
}

func (w *World) spawnKind(kind byte) bool {
	k, ok := enemyKinds[kind]
	if !ok {
		k = enemyKinds[TileEnemy]
	}
	hpMul := 1 + 0.05*float64(w.Wave-1)
	if hpMul > 2 {
		hpMul = 2
	}
	p := &w.Player
	minDist := 5.0
	if enemyKinds[kind].Boss {
		minDist = 7
	}
	for tries := 0; tries < 300; tries++ {
		x := w.rng.Intn(w.W-2) + 1
		y := w.rng.Intn(w.H-2) + 1
		if w.solidAt(x, y) {
			continue
		}
		fx, fy := float64(x)+0.5, float64(y)+0.5
		d := math.Hypot(fx-p.X, fy-p.Y)
		if d < minDist || d > 14 {
			continue
		}
		clash := false
		for _, e := range w.Enemies {
			if !e.Dead && math.Hypot(e.X-fx, e.Y-fy) < 1.2 {
				clash = true
				break
			}
		}
		if clash {
			continue
		}
		w.Enemies = append(w.Enemies, &Enemy{
			X: fx, Y: fy, Kind: kind, Score: k.Score,
			HP:      int(float64(k.HP)*hpMul + 0.5),
			Heading: w.rng.Float64() * 2 * math.Pi,
		})
		return true
	}
	return false
}

func (w *World) spawnWavePickups() {
	for _, kind := range []byte{TileHealth, TileAmmo, TileChest} {
		for tries := 0; tries < 200; tries++ {
			x := w.rng.Intn(w.W-2) + 1
			y := w.rng.Intn(w.H-2) + 1
			if w.solidAt(x, y) {
				continue
			}
			fx, fy := float64(x)+0.5, float64(y)+0.5
			p := &w.Player
			if math.Hypot(fx-p.X, fy-p.Y) < 3 {
				continue
			}
			w.Pickups = append(w.Pickups, Pickup{X: fx, Y: fy, Kind: kind, Bob0: w.rng.Float64() * 6})
			break
		}
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
