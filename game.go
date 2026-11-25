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
	if g.state == statePlay && g.world != nil && g.world.Survive && g.world.Wave != g.lvlWave {
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
		g.menuIdx = (g.menuIdx + mv + 2) % 2
		if anyConfirm() {
			if g.menuIdx == 0 {
				g.mode = modeCampaign
			} else {
				g.mode = modeSurvive
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
		if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
			g.state = statePause
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
	} else {
		g.world = loadLevel(n, 20240101+int64(n))
		g.level = n
		g.lvlWave = 0
		curLevel = n
	}
	rebuildFogSets()
	g.grantLoadout()
}

// unlockedCount : armes disponibles selon la progression ; le canon
// d'ions s'ajoute au niveau 5.
func (g *Game) unlockedCount() int {
	n := maxInt(1, min(3, g.level))
	if g.level >= 5 {
		n = 4
	}
	return n
}

func (g *Game) grantLoadout() {
	p := &g.world.Player
	p.Weapon = 0
	if p.Ammo[0] < 24 {
		p.Ammo[0] = 24
	}
	if g.unlockedCount() >= 2 && p.Ammo[1] < 6 {
		p.Ammo[1] = 6
	}
	if g.unlockedCount() >= 3 && p.Ammo[2] < 40 {
		p.Ammo[2] = 40
	}
	if g.unlockedCount() >= 4 && p.Ammo[3] < 5 {
		p.Ammo[3] = 5
	}
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

	mx, _ := ebiten.CursorPosition()
	if !lastMouse {
		lastMouse = true
	} else {
		p.Dir += float64(mx-lastMouseX) * mouseSens
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
	switchKeys := []ebiten.Key{ebiten.KeyDigit1, ebiten.KeyDigit2, ebiten.KeyDigit3, ebiten.KeyDigit4}
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
