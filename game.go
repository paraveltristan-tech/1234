package main

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type gameMode int

const (
	modeCampaign gameMode = iota
	modeSurvive
	modeCS
)

var (
	lastMouseX int
	lastMouse  bool
)

func (g *Game) Update() error {
	g.time += 1.0 / 60.0
	if !assetsInited {
		initAssets()
		initRender()
		assetsInited = true
	}
	if g.state == statePlay && g.mode == modeSurvive && g.world != nil && g.world.Wave != g.lvlWave {
		g.lvlWave = g.world.Wave
		g.level = clampInt(g.world.Wave, 1, 5)
		curLevel = (g.world.Wave-1)%5 + 1
		rebuildFogSets()
		g.grantLoadout()
	}
	if g.state == statePlay {
		ebiten.SetCursorMode(ebiten.CursorModeCaptured)
	} else {
		ebiten.SetCursorMode(ebiten.CursorModeVisible)
	}
	switch g.state {
	case stateTitle:
		mv := 0
		if inpututil.IsKeyJustPressed(ebiten.KeyUp) || inpututil.IsKeyJustPressed(ebiten.KeyW) {
			mv--
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyDown) || inpututil.IsKeyJustPressed(ebiten.KeyS) {
			mv++
		}
		g.menuIdx = (g.menuIdx + mv + 3) % 3
		if anyConfirm() {
			switch g.menuIdx {
			case 0:
				g.mode = modeCampaign
			case 1:
				g.mode = modeSurvive
			default:
				g.mode = modeCS
			}
			g.score = 0
			g.kills = 0
			g.ShowMinimap = true
			g.ShowBigMap = false
			g.startLevel(1)
			g.state = statePlay
			ebiten.SetCursorMode(ebiten.CursorModeCaptured)
		}
	case statePlay:
		g.handleToggleKeys()
		if inpututil.IsKeyJustPressed(ebiten.KeyR) && !g.buyOpen {
			g.world.startReload()
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
			g.state = statePause
			g.zoomed = false
			return nil
		}
		if g.mode == modeCS && inpututil.IsKeyJustPressed(ebiten.KeyB) && !g.world.Cleared {
			g.buyOpen = !g.buyOpen
		}
		if g.buyOpen {
			g.handleBuyKeys()
			return nil
		}
		if g.world.Cleared {
			g.clearT += 1.0 / 60.0
			if g.clearT > 2.2 {
				g.nextLevelOrWin()
			}
			return nil
		}
		g.updatePlayer()
		g.world.update(1.0/60.0, g)
	case statePause:
		g.handleToggleKeys()
		if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
			g.state = statePlay
			lastMouse = false
		}
	case stateDead, stateWin:
		if anyConfirm() {
			g.score = 0
			g.kills = 0
			g.buyOpen = false
			g.zoomed = false
			g.startLevel(1)
			g.state = statePlay
			lastMouse = false
			ebiten.SetCursorMode(ebiten.CursorModeCaptured)
		}
	}
	return nil
}

func anyConfirm() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeySpace) ||
		inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft)
}

// handleBuyKeys : boutique du mode counter-strike (touche B pour ouvrir).
func (g *Game) handleBuyKeys() {
	p := &g.world.Player
	press := func(k ebiten.Key) bool { return inpututil.IsKeyJustPressed(k) }
	switch {
	case press(ebiten.KeyDigit1):
		g.tryBuy(300, 1, 12)
	case press(ebiten.KeyDigit2):
		g.tryBuy(700, 2, 40)
	case press(ebiten.KeyDigit3):
		g.tryBuy(900, 3, 6)
	case press(ebiten.KeyDigit4):
		g.tryBuy(1200, 4, 4)
	case press(ebiten.KeyDigit5):
		if g.world.cash >= 250 && p.HP < maxHP {
			g.world.cash -= 250
			p.HP = min(maxHP, p.HP+40)
			g.world.FlashGood = 0.4
		}
	case press(ebiten.KeyB) || press(ebiten.KeyEscape):
		g.buyOpen = false
	}
}

func (g *Game) tryBuy(price, wi, ammo int) {
	w := g.world
	p := &w.Player
	if w.cash < price || p.Ammo[wi]+p.Reserve[wi] >= weapons[wi].MaxAmmo {
		return
	}
	w.cash -= price
	give := min(ammo, weapons[wi].MaxAmmo-p.Ammo[wi]-p.Reserve[wi])
	fill := min(weapons[wi].MagSize-p.Ammo[wi], give)
	if fill < 0 {
		fill = 0
	}
	p.Ammo[wi] += fill
	p.Reserve[wi] += give - fill
	w.FlashGood = 0.3
}

func (g *Game) handleToggleKeys() {
	if inpututil.IsKeyJustPressed(ebiten.KeyM) {
		g.ShowMinimap = !g.ShowMinimap
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyTab) {
		g.ShowBigMap = !g.ShowBigMap
	}
}

func (g *Game) startLevel(n int) {
	g.clearT = 0
	lastMouse = false
	if g.mode == modeSurvive {
		g.world = loadLevel(levelSurvive, 20240101)
		g.level = 1
		g.lvlWave = 1
		curLevel = 1
	} else if g.mode == modeCS {
		g.world = loadLevel(levelCS, 20240101)
		g.level = 1
		g.lvlWave = 1
		curLevel = 6
	} else {
		g.world = loadLevel(n, 20240101+int64(n))
		g.level = n
		g.lvlWave = 0
		curLevel = n
	}
	rebuildFogSets()
	g.grantLoadout()
}

// unlockedCount : armes disponibles selon la progression ; le sniper
// s'ajoute au niveau 4, le canon d'ions au niveau 5. En mode counter-strike,
// tout s'achete a la boutique.
func (g *Game) unlockedCount() int {
	if g.mode == modeCS {
		return len(weapons)
	}
	n := maxInt(1, min(3, g.level))
	if g.level >= 4 {
		n = 4
	}
	if g.level >= 5 {
		n = 5
	}
	return n
}

func (g *Game) grantLoadout() {
	p := &g.world.Player
	p.Weapon = 0
	g.topUp(0, 24)
	if g.mode == modeCS {
		for wi := 1; wi < len(weapons); wi++ {
			if p.Ammo[wi]+p.Reserve[wi] == 0 {
				g.topUp(wi, 1)
			}
		}
		return
	}
	if g.unlockedCount() >= 2 {
		g.topUp(1, 6)
	}
	if g.unlockedCount() >= 3 {
		g.topUp(2, 40)
	}
	if g.unlockedCount() >= 4 {
		g.topUp(3, 5)
	}
	if g.unlockedCount() >= 5 {
		g.topUp(4, 5)
	}
}

// topUp complete chargeur puis reserve jusqu'a la quantite voulue.
func (g *Game) topUp(wi, want int) {
	p := &g.world.Player
	total := p.Ammo[wi] + p.Reserve[wi]
	if total >= want || total >= weapons[wi].MaxAmmo {
		return
	}
	give := min(want, weapons[wi].MaxAmmo) - total
	fill := min(weapons[wi].MagSize-p.Ammo[wi], give)
	if fill < 0 {
		fill = 0
	}
	p.Ammo[wi] += fill
	p.Reserve[wi] += give - fill
}

func (g *Game) nextLevelOrWin() {
	if g.level >= 5 {
		g.state = stateWin
		return
	}
	g.startLevel(g.level + 1)
}

func (g *Game) updatePlayer() {
	w := g.world
	p := &w.Player

	g.switchWeapon()

	sens := mouseSens
	if g.zoomed {
		sens *= 0.35
	}
	mx, _ := ebiten.CursorPosition()
	if !lastMouse {
		lastMouse = true
	} else {
		p.Dir += float64(mx-lastMouseX) * sens
	}
	lastMouseX = mx

	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		p.Dir -= turnSpeed / 60.0
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		p.Dir += turnSpeed / 60.0
	}

	speed := moveSpeed / 60.0
	if ebiten.IsKeyPressed(ebiten.KeyShift) {
		speed *= sprintMult
	}
	c, s := math.Cos(p.Dir), math.Sin(p.Dir)
	var fw, sw float64
	if ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyUp) {
		fw += 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) || ebiten.IsKeyPressed(ebiten.KeyDown) {
		fw -= 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		sw -= 0.75
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		sw += 0.75
	}
	dx := (c*fw - s*sw) * speed
	dy := (s*fw + c*sw) * speed
	p.X, p.Y = tryMove(w.Grid, w.W, w.H, p.X, p.Y, dx, dy, playerRadius)

	g.zoomed = p.Weapon == 3 && ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight)

	wp := weapons[p.Weapon]
	if (wp.Auto && ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)) ||
		(!wp.Auto && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft)) {
		w.fire(g)
	}
	if g.world.Player.HP <= 0 {
		g.state = stateDead
		ebiten.SetCursorMode(ebiten.CursorModeVisible)
	}
}

func (g *Game) switchWeapon() {
	p := &g.world.Player
	u := g.unlockedCount()
	switchKeys := []ebiten.Key{ebiten.KeyDigit1, ebiten.KeyDigit2, ebiten.KeyDigit3, ebiten.KeyDigit4, ebiten.KeyDigit5}
	for i, k := range switchKeys {
		if inpututil.IsKeyJustPressed(k) && i < u {
			p.Weapon = i
		}
	}
	_, wy := ebiten.Wheel()
	g.wheelAcc += wy
	if g.wheelAcc >= 1 {
		p.Weapon = min(u-1, p.Weapon+1)
		g.wheelAcc = 0
	}
	if g.wheelAcc <= -1 {
		p.Weapon = maxInt(0, p.Weapon-1)
		g.wheelAcc = 0
	}
}
