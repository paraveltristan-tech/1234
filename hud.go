package main

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

func drawText(dst *ebiten.Image, s string, face text.Face, x, y float64, col color.Color, align text.Align) {
	op := &text.DrawOptions{}
	op.GeoM.Translate(x, y)
	op.PrimaryAlign = align
	op.ColorScale.ScaleWithColor(col)
	text.Draw(dst, s, face, op)
}

func itoa(n int) string { return fmt.Sprintf("%d", n) }

// ---- weapon ----

func (g *Game) drawWeapon() {
	p := &g.world.Player
	bobX := math.Sin(p.Bob*0.6) * 3
	bobY := math.Abs(math.Cos(p.Bob*0.6)) * 2
	wp := weapons[p.Weapon]
	recoilPct := 0.0
	if wp.FlashMax > 0 {
		recoilPct = g.world.Muzzle / wp.FlashMax
	}
	recoil := recoilPct * wp.Recoil

	gx := float64(InternalW)/2 + 18 + bobX
	gy := float64(InternalH) - 78 + bobY + recoil
	fx, fy := float32(gx), float32(gy)

	switch p.Weapon {
	case 1: // fusil à pompe
		vector.DrawFilledRect(screenBuf, fx-14, fy-4, 58, 12, hex(0x232a33), false)
		vector.DrawFilledRect(screenBuf, fx-14, fy-1, 58, 3, hex(0x0f1216), false)
		vector.DrawFilledRect(screenBuf, fx+6, fy+8, 28, 34, hex(0x4a3423), false)
		vector.DrawFilledRect(screenBuf, fx+10, fy+10, 20, 8, hex(0x5d4430), false)
		vector.DrawFilledRect(screenBuf, fx-2, fy+30, 24, 30, hex(0x33261b), false)
	case 2: // plasma
		vector.DrawFilledRect(screenBuf, fx-4, fy-10, 44, 16, hex(0x1c3a40), false)
		vector.DrawFilledRect(screenBuf, fx-2, fy-8, 40, 5, hex(0x39d7e0), false)
		vector.DrawFilledRect(screenBuf, fx, fy+4, 34, 36, hex(0x24363d), false)
		vector.DrawFilledCircle(screenBuf, fx+30, fy+22, 7, hex(0x0f5e66), false)
		vector.DrawFilledCircle(screenBuf, fx+30, fy+22, 4, hex(0x5fe8ff), false)
		vector.DrawFilledRect(screenBuf, fx-2, fy+38, 22, 26, hex(0x1d232b), false)
	case 3: // canon d'ions
		vector.DrawFilledRect(screenBuf, fx-10, fy-6, 60, 18, hex(0x1a2620), false)
		vector.DrawFilledRect(screenBuf, fx-10, fy-1, 60, 4, hex(0x0c1210), false)
		pulse := 6 + int(3*math.Sin(p.Bob*1.4))
		vector.DrawFilledCircle(screenBuf, fx+44, fy+3, float32(pulse+6), hex(0x0f5e46), false)
		vector.DrawFilledCircle(screenBuf, fx+44, fy+3, float32(pulse+3), hex(0x1c8a5e), false)
		vector.DrawFilledCircle(screenBuf, fx+44, fy+3, float32(pulse), hex(0x5fe8a0), false)
		vector.DrawFilledRect(screenBuf, fx, fy+10, 30, 34, hex(0x232b28), false)
		vector.DrawFilledRect(screenBuf, fx+2, fy+12, 26, 4, hex(0x35433d), false)
	default: // pistolet
		vector.DrawFilledRect(screenBuf, fx-2, fy-6, 34, 14, hex(0x232a33), false)
		vector.DrawFilledRect(screenBuf, fx, fy, 26, 40, hex(0x39434f), false)
		vector.DrawFilledRect(screenBuf, fx+2, fy+2, 22, 8, hex(0x4a565f), false)
		vector.DrawFilledRect(screenBuf, fx-1, fy+38, 20, 28, hex(0x1d232b), false)
		vector.DrawFilledRect(screenBuf, fx+6, fy+16, 6, 10, hex(0x14181d), false)
	}

	if g.world.Muzzle > 0 {
		if p.Weapon == 3 {
			vector.DrawFilledCircle(screenBuf, fx+30, fy-10, 15, hex(0x5fe8a0), false)
			vector.DrawFilledCircle(screenBuf, fx+30, fy-10, 8, hex(0xd0ffe8), false)
		} else if p.Weapon == 2 {
			vector.DrawFilledCircle(screenBuf, fx+18, fy-16, 11, hex(0x5fe8ff), false)
			vector.DrawFilledCircle(screenBuf, fx+18, fy-16, 6, hex(0xd0faff), false)
		} else {
			r := float32(9)
			if p.Weapon == 1 {
				r = 13
			}
			vector.DrawFilledCircle(screenBuf, fx+15, fy-12, r, hex(0xffe08a), false)
			vector.DrawFilledCircle(screenBuf, fx+15, fy-12, r/2, hex(0xfff6cf), false)
		}
	}
}

// ---- hud ----

func (g *Game) drawHUD(screen *ebiten.Image) {
	w := float32(InternalW * 2)
	h := float32(InternalH * 2)
	p := &g.world.Player
	wp := weapons[p.Weapon]

	vector.DrawFilledRect(screen, 12, h-46, 20, 34, rgbaSemi(40, 30, 34, 210), false)
	verticalBar(screen, 16, int(h)-12, 30, p.HP)
	drawText(screen, "SANTÉ", faceTiny, 40, float64(h)-44, hex(0xf2f2f4), text.AlignStart)
	drawText(screen, fmt.Sprintf("%d", clampInt(p.HP, 0, 100)), faceMed, 40, float64(h)-26, hex(0xf2f2f4), text.AlignStart)

	vector.DrawFilledRect(screen, w-308, h-46, 296, 34, rgbaSemi(30, 34, 40, 210), false)
	for wi := 0; wi < len(weapons); wi++ {
		sx := w - 300 + float32(wi)*27
		boxCol := hex(0x1c2128)
		if wi == p.Weapon {
			boxCol = hex(0x39d7e0)
		} else if wi >= g.unlockedCount() {
			boxCol = hex(0x12141a)
		}
		vector.DrawFilledRect(screen, sx, h-40, 22, 22, boxCol, false)
		drawText(screen, fmt.Sprintf("%d", wi+1), faceTiny, float64(sx)+11, float64(h)-37, hex(0xf2f2f4), text.AlignCenter)
	}
	drawText(screen, fmt.Sprintf("%d · %s", p.Ammo[p.Weapon], wp.Name), faceSmall, float64(w)-14, float64(h)-40, hex(0xffc93a), text.AlignEnd)

	if g.world.Survive {
		drawText(screen, fmt.Sprintf("VAGUE %d", g.world.Wave), faceSmall, 14, 12, hex(0x9fd8e8), text.AlignStart)
	} else {
		drawText(screen, fmt.Sprintf("NIVEAU %d", g.level), faceSmall, 14, 12, hex(0x9fd8e8), text.AlignStart)
	}
	drawText(screen, fmt.Sprintf("SCORE %d", g.score), faceSmall, float64(w)-14, 12, hex(0x9fd8e8), text.AlignEnd)
	drawText(screen, fmt.Sprintf("ENNEMIS %d", g.aliveEnemies()), faceTiny, float64(w)-14, 34, hex(0x8890a0), text.AlignEnd)
	drawText(screen, fmt.Sprintf("ARME %d/%d", p.Weapon+1, g.unlockedCount()), faceTiny, 14, 34, hex(0x8890a0), text.AlignStart)

	cx, cy := w/2, h/2
	vector.DrawFilledRect(screen, cx-1, cy-7, 2, 4, hex(0xd0f6ff), false)
	vector.DrawFilledRect(screen, cx-1, cy+3, 2, 4, hex(0xd0f6ff), false)
	vector.DrawFilledRect(screen, cx-7, cy-1, 4, 2, hex(0xd0f6ff), false)
	vector.DrawFilledRect(screen, cx+3, cy-1, 4, 2, hex(0xd0f6ff), false)

	if g.world.Cleared {
		drawText(screen, "NIVEAU TERMINÉ !  +BONUS", faceBig, float64(w)/2, float64(h)/2-60, hex(0x5fe8c0), text.AlignCenter)
	}
	drawText(screen, "M MINIMAP · TAB CARTE · ÉCHAP PAUSE", faceTiny, float64(w)/2, float64(h)-16, hex(0x6a7686), text.AlignCenter)
}

func (g *Game) aliveEnemies() int {
	n := 0
	for _, e := range g.world.Enemies {
		if !e.Dead {
			n++
		}
	}
	return n
}

func verticalBar(screen *ebiten.Image, x, yBottom, heightFull int, hp int) {
	frac := float64(clampInt(hp, 0, 100)) / 100.0
	backH := float32(heightFull)
	vector.DrawFilledRect(screen, float32(x), float32(yBottom)-backH, 10, backH, rgbaSemi(60, 18, 24, 230), false)
	fh := float32(float64(heightFull) * frac)
	col := hex(0x4fd06a)
	if hp <= 30 {
		col = hex(0xe0483f)
	} else if hp <= 60 {
		col = hex(0xe8a83c)
	}
	vector.DrawFilledRect(screen, float32(x), float32(yBottom)-fh, 10, fh, col, false)
}

func clampInt(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

func rgbaSemi(r, gg, b, a uint8) color.Color {
	return color.RGBA{R: r, G: gg, B: b, A: a}
}

// ---- screens & overlays ----

func (g *Game) drawOverlays(screen *ebiten.Image) {
	w := float32(InternalW * 2)
	h := float32(InternalH * 2)
	if r := g.world.FlashRed; r > 0 {
		vector.DrawFilledRect(screen, 0, 0, w, h, rgbaSemi(200, 20, 30, uint8(r*130)), false)
	}
	if gr := g.world.FlashGood; gr > 0 {
		vector.DrawFilledRect(screen, 0, 0, w, h, rgbaSemi(40, 200, 90, uint8(gr*70)), false)
	}
	if m := g.world.Muzzle; m > 0 {
		vector.DrawFilledRect(screen, 0, 0, w, h, rgbaSemi(255, 240, 200, uint8(m*900)), false)
	}
}

func (g *Game) drawDeathScreen(screen *ebiten.Image) {
	w := float64(InternalW * 2)
	h := float64(InternalH * 2)
	vector.DrawFilledRect(screen, 0, 0, float32(w), float32(h), rgbaSemi(60, 0, 4, 150), false)
	drawText(screen, "VOUS ÊTES MORT", faceTitle, w/2, h*0.30, hex(0xff5a5a), text.AlignCenter)
	if g.world.Survive {
		drawText(screen, fmt.Sprintf("VAGUE ATTEINTE : %d  ·  SCORE : %d", g.world.Wave, g.score), faceMed, w/2, h*0.52, hex(0xf2f2f4), text.AlignCenter)
	} else {
		drawText(screen, fmt.Sprintf("SCORE : %d  ·  ÉLIMINATIONS : %d", g.score, g.kills), faceMed, w/2, h*0.52, hex(0xf2f2f4), text.AlignCenter)
	}
	if int(math.Mod(g.time, 1.2)/0.6)%2 == 0 {
		drawText(screen, "ENTRÉE POUR RECOMMENCER", faceSmall, w/2, h*0.68, hex(0xffe066), text.AlignCenter)
	}
}

func (g *Game) drawPauseScreen(screen *ebiten.Image) {
	w := float64(InternalW * 2)
	h := float64(InternalH * 2)
	vector.DrawFilledRect(screen, 0, 0, float32(w), float32(h), rgbaSemi(6, 10, 18, 190), false)
	drawText(screen, "PAUSE", faceTitle, w/2, h*0.30, hex(0x5fe8ff), text.AlignCenter)
	if g.world.Survive {
		drawText(screen, fmt.Sprintf("VAGUE %d  ·  SCORE %d", g.world.Wave, g.score), faceMed, w/2, h*0.50, hex(0xf2f2f4), text.AlignCenter)
	} else {
		drawText(screen, fmt.Sprintf("NIVEAU %d  ·  SCORE %d", g.level, g.score), faceMed, w/2, h*0.50, hex(0xf2f2f4), text.AlignCenter)
	}
	if int(math.Mod(g.time, 1.2)/0.6)%2 == 0 {
		drawText(screen, "ÉCHAP POUR REPRENDRE", faceSmall, w/2, h*0.64, hex(0xffe066), text.AlignCenter)
	}
}

func (g *Game) drawTitle(screen *ebiten.Image) {
	w := float64(InternalW * 2)
	h := float64(InternalH * 2)
	drawText(screen, "PROTOCOLE NÉON", faceTitle, w/2, h*0.10, hex(0x5fe8ff), text.AlignCenter)
	drawText(screen, "FPS RÉTRO · PARTICULES · CANON D'IONS", faceMed, w/2, h*0.24, hex(0x8890a0), text.AlignCenter)

	opts := [2]string{"CAMPAGNE — 5 NIVEAUX", "SURVIE — VAGUES INFINIES"}
	for i, o := range opts {
		col := hex(0x8890a0)
		label := o
		if i == g.menuIdx {
			label = "» " + o + " «"
			col = hex(0xffe066)
		}
		drawText(screen, label, faceMed, w/2, h*0.34+float64(i)*38, col, text.AlignCenter)
	}

	lines := []string{
		"ZQSD/WASD .... SE DÉPLACER   ·   SOURIS+CLIC ... VISER/TIRER",
		"MAJ .......... SPRINT   ·   1/2/3/4 + MOLETTE ... ARMES",
		"M ............ MINIMAP   ·   TAB ... GRANDE CARTE",
		"ÉCHAP ........ PAUSE   ·   FLÈCHES + ENTRÉE ... CHOISIR LE MODE",
	}
	for i, l := range lines {
		drawText(screen, l, faceSmall, w/2, h*0.52+float64(i)*28, hex(0xc8d2dd), text.AlignCenter)
	}
	if int(math.Mod(g.time, 1.2)/0.6)%2 == 0 {
		drawText(screen, "ENTRÉE POUR JOUER", faceMed, w/2, h*0.80, hex(0xffe066), text.AlignCenter)
	}
}

func (g *Game) drawWinScreen(screen *ebiten.Image) {
	w := float64(InternalW * 2)
	h := float64(InternalH * 2)
	drawText(screen, "VICTOIRE !", faceTitle, w/2, h*0.30, hex(0x5fe8c0), text.AlignCenter)
	drawText(screen, fmt.Sprintf("SCORE FINAL : %d", g.score), faceBig, w/2, h*0.50, hex(0xf2f2f4), text.AlignCenter)
	drawText(screen, "ENTRÉE POUR REJOUER", faceSmall, w/2, h*0.66, hex(0xffe066), text.AlignCenter)
}
