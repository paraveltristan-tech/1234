package main

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
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
	if g.state == statePlay {
		ebiten.SetCursorMode(ebiten.CursorModeCaptured)
	} else {
		ebiten.SetCursorMode(ebiten.CursorModeVisible)
	}
	switch g.state {
	case stateTitle:
		if anyConfirm() {
			g.score = 0
			g.kills = 0
			g.startLevel(1)
			g.state = statePlay
			ebiten.SetCursorMode(ebiten.CursorModeCaptured)
		}
	case statePlay:
		if g.world.Cleared {
			g.clearT += 1.0 / 60.0
			if g.clearT > 2.2 {
				g.nextLevelOrWin()
			}
			return nil
		}
		g.updatePlayer()
		g.world.update(1.0/60.0, g)
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

func (g *Game) startLevel(n int) {
	g.level = n
	g.world = loadLevel(n, 20240101+int64(n))
	g.clearT = 0
	lastMouse = false
}

func (g *Game) nextLevelOrWin() {
	if g.level >= 3 {
		g.state = stateWin
		return
	}
	g.startLevel(g.level + 1)
}

func (g *Game) updatePlayer() {
	w := g.world
	p := &w.Player

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

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		w.fire(g)
	}
	if g.world.Player.HP <= 0 {
		g.state = stateDead
		ebiten.SetCursorMode(ebiten.CursorModeVisible)
	}
}
