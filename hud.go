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
	if g.zoomed {
		return
	}
	bobX := math.Sin(p.Bob*0.6) * 3
	bobY := math.Abs(math.Cos(p.Bob*0.6)) * 2
	wp := weapons[p.Weapon]
	recoilPct := 0.0
	if wp.FlashMax > 0 {
		recoilPct = g.world.Muzzle / wp.FlashMax
	}
	recoil := recoilPct * wp.Recoil
	reloadDip := 0.0
	if p.ReloadT > 0 {
		reloadDip = math.Sin((1-p.ReloadT/math.Max(0.01, p.ReloadDur))*math.Pi) * 54
	}

	gx := float64(InternalW)/2 + 18 + bobX
	gy := float64(InternalH) - 78 + bobY + recoil + reloadDip
	fx, fy := float32(gx), float32(gy)
	pump := float32(recoilPct * 6)

	switch p.Weapon {
	case 1: // fusil a pompe : crosse veinée, culasse acier, canons doubles, pompe animée
		vector.DrawFilledRect(screenBuf, fx-14, fy-4, 58, 12, hex(0x232a33), false)
		vector.DrawFilledRect(screenBuf, fx-14, fy-4, 58, 3, hex(0x39434f), false)
		vector.DrawFilledRect(screenBuf, fx-14, fy-1, 58, 3, hex(0x0f1216), false)
		vector.DrawFilledRect(screenBuf, fx-16, fy-6, 4, 16, hex(0x11161b), false)
		vector.DrawFilledRect(screenBuf, fx-16, fy+6, 4, 5, hex(0x11161b), false)
		vector.DrawFilledRect(screenBuf, fx+6, fy+8, 28, 34, hex(0x4a3423), false)
		vector.DrawFilledRect(screenBuf, fx+8, fy+10, 24, 5, hex(0x5d4430), false)
		vector.DrawFilledRect(screenBuf, fx+8, fy+20, 24, 3, hex(0x3d2c1e), false)
		vector.DrawFilledRect(screenBuf, fx+10, fy+30, 20, 3, hex(0x5d4430), false)
		vector.DrawFilledRect(screenBuf, fx-2, fy+30, 24, 30, hex(0x33261b), false)
		vector.DrawFilledRect(screenBuf, fx-6+pump, fy+8+pump*0.7, 22, 9, hex(0x4a3423), false)
		vector.DrawFilledRect(screenBuf, fx-6+pump, fy+11+pump*0.7, 22, 3, hex(0x6d5238), false)
		vector.DrawFilledRect(screenBuf, fx+2, fy-8, 8, 6, hex(0xffc93a), false)
	case 2: // fusil a plasma : cellule cyan pulseante, serpentin, évents
		vector.DrawFilledRect(screenBuf, fx-6, fy-12, 46, 18, hex(0x1c3a40), false)
		vector.DrawFilledRect(screenBuf, fx-6, fy-12, 46, 3, hex(0x2c5a63), false)
		vector.DrawFilledRect(screenBuf, fx-4, fy-8, 42, 5, hex(0x39d7e0), false)
		for i := 0; i < 4; i++ {
			vector.DrawFilledCircle(screenBuf, float32(2+i*11)+fx-2, fy-10, 3, hex(0x5fe8ff), false)
		}
		vector.DrawFilledRect(screenBuf, fx, fy+4, 34, 36, hex(0x24363d), false)
		vector.DrawFilledRect(screenBuf, fx+2, fy+6, 30, 4, hex(0x39546a), false)
		vector.DrawFilledRect(screenBuf, fx+2, fy+14, 30, 2, hex(0x1a5a63), false)
		vector.DrawFilledCircle(screenBuf, fx+30, fy+22, 8, hex(0x0f5e66), false)
		pulse := 4 + int(2*math.Sin(p.Bob*2.2))
		vector.DrawFilledCircle(screenBuf, fx+30, fy+22, float32(pulse), hex(0x5fe8ff), false)
		vector.DrawFilledCircle(screenBuf, fx+30, fy+22, float32(pulse-2), hex(0xd0faff), false)
		vector.DrawFilledRect(screenBuf, fx-2, fy+38, 22, 26, hex(0x1d232b), false)
		vector.DrawFilledRect(screenBuf, fx-2, fy+44, 22, 3, hex(0x39d7e0), false)
	case 3: // sniper : canon long, lunette à lentille bleue, frein de bouche, liseré néon
		vector.DrawFilledRect(screenBuf, fx-28, fy-2, 74, 12, hex(0x1c2630), false)
		vector.DrawFilledRect(screenBuf, fx-28, fy-2, 74, 3, hex(0x2f3f4c), false)
		vector.DrawFilledRect(screenBuf, fx-28, fy+7, 74, 2, hex(0x39d7e0), false)
		vector.DrawFilledRect(screenBuf, fx-34, fy-4, 10, 10, hex(0x11161b), false)
		vector.DrawFilledRect(screenBuf, fx-34, fy-1, 10, 2, hex(0x333c45), false)
		vector.DrawFilledRect(screenBuf, fx+2, fy-12, 18, 6, hex(0x141a20), false)
		vector.DrawFilledRect(screenBuf, fx+8, fy-8, 3, 6, hex(0x141a20), false)
		vector.DrawFilledCircle(screenBuf, fx+20, fy-11, 5, hex(0x0e1116), false)
		vector.DrawFilledCircle(screenBuf, fx+20, fy-11, 3, hex(0x59c8ff), false)
		vector.DrawFilledCircle(screenBuf, fx+1, fy-11, 4, hex(0x0e1116), false)
		vector.DrawFilledRect(screenBuf, fx+2, fy+10, 26, 30, hex(0x22272e), false)
		vector.DrawFilledRect(screenBuf, fx+4, fy+12, 22, 5, hex(0x333c45), false)
		vector.DrawFilledRect(screenBuf, fx+4, fy+22, 22, 3, hex(0x22272e), false)
		vector.DrawFilledRect(screenBuf, fx+6, fy+14, 4, 8, hex(0x4fd06a), false)
		// culasse et levier
		vector.DrawFilledRect(screenBuf, fx+24, fy-6, 8, 4, hex(0x4a565f), false)
		vector.DrawFilledRect(screenBuf, fx+30, fy-6, 4, 8, hex(0x4a565f), false)
	case 4: // canon d'ions : double serpentin, orbre géante pulsante, bandes danger
		vector.DrawFilledRect(screenBuf, fx-10, fy-6, 60, 18, hex(0x1a2620), false)
		vector.DrawFilledRect(screenBuf, fx-10, fy-1, 60, 4, hex(0x0c1210), false)
		vector.DrawFilledRect(screenBuf, fx-10, fy-6, 60, 2, hex(0x2c4438), false)
		for i := 0; i < 3; i++ {
			vector.DrawFilledCircle(screenBuf, float32(6+i*14)+fx-6, fy+3, 5, hex(0x0f3326), false)
			vector.DrawFilledCircle(screenBuf, float32(6+i*14)+fx-6, fy+3, 2, hex(0x2c8a5e), false)
		}
		pulse := 6 + int(3*math.Sin(p.Bob*1.4))
		vector.DrawFilledCircle(screenBuf, fx+44, fy+3, float32(pulse+6), hex(0x0f5e46), false)
		vector.DrawFilledCircle(screenBuf, fx+44, fy+3, float32(pulse+3), hex(0x1c8a5e), false)
		vector.DrawFilledCircle(screenBuf, fx+44, fy+3, float32(pulse), hex(0x5fe8a0), false)
		vector.DrawFilledRect(screenBuf, fx, fy+10, 30, 34, hex(0x232b28), false)
		vector.DrawFilledRect(screenBuf, fx+2, fy+12, 26, 4, hex(0x35433d), false)
		vector.DrawFilledRect(screenBuf, fx-8, fy+16, 4, 4, hex(0xffc93a), false)
		vector.DrawFilledRect(screenBuf, fx-4, fy+16, 4, 4, hex(0x1a1c14), false)
		vector.DrawFilledRect(screenBuf, fx-8, fy+22, 4, 4, hex(0xffc93a), false)
		vector.DrawFilledRect(screenBuf, fx-4, fy+22, 4, 4, hex(0x1a1c14), false)
	default: // pistolet : glissière bicolore, rayures de culasse, poignée texturée
		vector.DrawFilledRect(screenBuf, fx-2, fy-6, 34, 14, hex(0x232a33), false)
		vector.DrawFilledRect(screenBuf, fx-2, fy-6, 34, 3, hex(0x39434f), false)
		for i := 0; i < 4; i++ {
			vector.DrawFilledRect(screenBuf, fx+float32(2+i*7), fy-4, 2, 9, hex(0x14181d), false)
		}
		vector.DrawFilledRect(screenBuf, fx, fy, 26, 40, hex(0x39434f), false)
		vector.DrawFilledRect(screenBuf, fx+2, fy+2, 22, 8, hex(0x4a565f), false)
		vector.DrawFilledRect(screenBuf, fx-1, fy+38, 20, 28, hex(0x1d232b), false)
		vector.DrawFilledRect(screenBuf, fx-1, fy+42, 20, 3, hex(0x333c45), false)
		vector.DrawFilledRect(screenBuf, fx-1, fy+50, 20, 3, hex(0x333c45), false)
		vector.DrawFilledRect(screenBuf, fx+6, fy+16, 6, 10, hex(0x14181d), false)
	}

	if g.world.Muzzle > 0 {
		switch p.Weapon {
		case 2:
			vector.DrawFilledCircle(screenBuf, fx+18, fy-16, 11, hex(0x5fe8ff), false)
			vector.DrawFilledCircle(screenBuf, fx+18, fy-16, 6, hex(0xd0faff), false)
		case 3:
			vector.DrawFilledCircle(screenBuf, fx-32, fy+1, 10, hex(0xfffbe8), false)
			vector.DrawFilledCircle(screenBuf, fx-32, fy+1, 5, hex(0xd0f6ff), false)
		case 4:
			vector.DrawFilledCircle(screenBuf, fx+30, fy-10, 15, hex(0x5fe8a0), false)
			vector.DrawFilledCircle(screenBuf, fx+30, fy-10, 8, hex(0xd0ffe8), false)
		default:
			r := float32(9)
			vector.DrawFilledCircle(screenBuf, fx+15, fy-12, r, hex(0xffe08a), false)
			vector.DrawFilledCircle(screenBuf, fx+15, fy-12, r/2, hex(0xfff6cf), false)
		}
	}
}

// drawScope : overlay de lunette quand le sniper est armé et zoomé.
func (g *Game) drawScope(screen *ebiten.Image) {
	w := float32(InternalW * 2)
	h := float32(InternalH * 2)
	vector.DrawFilledRect(screen, 0, 0, w, h, rgbaSemi(2, 4, 8, 110), false)
	cx, cy := w/2, h/2
	r := h * 0.44
	vector.DrawFilledCircle(screen, cx, cy, r, rgbaSemi(8, 14, 24, 170), false)
	vector.StrokeCircle(screen, cx, cy, r, 3, hex(0x39d7e0), false)
	vector.StrokeCircle(screen, cx, cy, r+6, 8, hex(0x060a10), false)
	vector.DrawFilledRect(screen, cx-1, cy-r, 2, 2*r, hex(0x9fd8e8), false)
	vector.DrawFilledRect(screen, cx-r, cy-1, 2*r, 2, hex(0x9fd8e8), false)
	for i := 1; i <= 4; i++ {
		vector.DrawFilledRect(screen, cx-6, cy+float32(i)*22, 12, 2, hex(0x9fd8e8), false)
	}
	vector.StrokeCircle(screen, cx, cy, 8, 1, hex(0xd0f6ff), false)
}

// ---- hud ----
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

	vector.DrawFilledRect(screen, w-352, h-46, 340, 34, rgbaSemi(30, 34, 40, 210), false)
	for wi := 0; wi < len(weapons); wi++ {
		sx := w - 344 + float32(wi)*27
		boxCol := hex(0x1c2128)
		if wi == p.Weapon {
			boxCol = hex(0x39d7e0)
		} else if wi >= g.unlockedCount() {
			boxCol = hex(0x12141a)
		}
		vector.DrawFilledRect(screen, sx, h-40, 22, 22, boxCol, false)
		drawText(screen, fmt.Sprintf("%d", wi+1), faceTiny, float64(sx)+11, float64(h)-37, hex(0xf2f2f4), text.AlignCenter)
	}
	if p.ReloadT > 0 {
		drawText(screen, "RECHARGEMENT...", faceTiny, float64(w)-14, float64(h)-40, hex(0x5fe8ff), text.AlignEnd)
	} else {
		drawText(screen, fmt.Sprintf("%d+%d · %s", p.Ammo[p.Weapon], p.Reserve[p.Weapon], wp.Name), faceSmall, float64(w)-14, float64(h)-40, hex(0xffc93a), text.AlignEnd)
	}

	if g.mode == modeCS {
		drawText(screen, fmt.Sprintf("MANCHE %d", g.world.Wave), faceSmall, 14, 12, hex(0xd8c88a), text.AlignStart)
		drawText(screen, fmt.Sprintf("$ %d", g.world.cash), faceMed, float64(w)-14, 48, hex(0xffe066), text.AlignEnd)
	} else if g.world.Survive {
		drawText(screen, fmt.Sprintf("VAGUE %d", g.world.Wave), faceSmall, 14, 12, hex(0x9fd8e8), text.AlignStart)
	} else {
		drawText(screen, fmt.Sprintf("NIVEAU %d", g.level), faceSmall, 14, 12, hex(0x9fd8e8), text.AlignStart)
	}
	drawText(screen, fmt.Sprintf("SCORE %d", g.score), faceSmall, float64(w)-14, 12, hex(0x9fd8e8), text.AlignEnd)
	drawText(screen, fmt.Sprintf("ENNEMIS %d", g.aliveEnemies()), faceTiny, float64(w)-14, 34, hex(0x8890a0), text.AlignEnd)
	drawText(screen, fmt.Sprintf("ARME %d/%d", p.Weapon+1, g.unlockedCount()), faceTiny, 14, 34, hex(0x8890a0), text.AlignStart)

	barY := float64(10)
	for _, bk := range [2]byte{TileColossus, TileBoss} {
		var be *Enemy
		for _, e := range g.world.Enemies {
			if e.Kind == bk && !e.Dead {
				be = e
				break
			}
		}
		if be == nil {
			continue
		}
		frac := float64(be.HP) / float64(enemyKinds[bk].HP)
		if frac < 0 {
			frac = 0
		}
		if frac > 1 {
			frac = 1
		}
		name, txtCol, barCol := "GARDIEN", hex(0xff7070), hex(0xd02a34)
		if bk == TileColossus {
			name, txtCol, barCol = "COLOSSE", hex(0xffb066), hex(0xff8c2a)
		}
		drawText(screen, name, faceTiny, float64(w)/2, barY, txtCol, text.AlignCenter)
		vector.DrawFilledRect(screen, w/2-112, float32(barY)+18, 224, 12, rgbaSemi(24, 6, 8, 210), false)
		vector.DrawFilledRect(screen, w/2-110, float32(barY)+20, float32(220*frac), 8, barCol, false)
		barY += 40
	}
	cx, cy := w/2, h/2
	vector.DrawFilledRect(screen, cx-1, cy-7, 2, 4, hex(0xd0f6ff), false)
	vector.DrawFilledRect(screen, cx-1, cy+3, 2, 4, hex(0xd0f6ff), false)
	vector.DrawFilledRect(screen, cx-7, cy-1, 4, 2, hex(0xd0f6ff), false)
	vector.DrawFilledRect(screen, cx+3, cy-1, 4, 2, hex(0xd0f6ff), false)
	if p.ReloadT > 0 {
		prog := 1 - p.ReloadT/math.Max(0.01, p.ReloadDur)
		drawText(screen, "RECHARGE", faceTiny, float64(cx), float64(cy)+14, hex(0x5fe8ff), text.AlignCenter)
		vector.DrawFilledRect(screen, cx-40, cy+24, 80, 6, rgbaSemi(10, 14, 20, 200), false)
		vector.DrawFilledRect(screen, cx-40, cy+24, float32(80*prog), 6, hex(0x39d7e0), false)
	} else if p.Ammo[p.Weapon] == 0 && p.Reserve[p.Weapon] > 0 {
		drawText(screen, "R POUR RECHARGER", faceTiny, float64(cx), float64(cy)+16, hex(0xff8c2a), text.AlignCenter)
	}

	if g.world.Cleared {
		drawText(screen, "NIVEAU TERMINÉ !  +BONUS", faceBig, float64(w)/2, float64(h)/2-60, hex(0x5fe8c0), text.AlignCenter)
	} else if g.world.bossAlive() {
		for _, ex := range g.world.Exits {
			if math.Hypot(ex.X-p.X, ex.Y-p.Y) < 3 {
				drawText(screen, "LE GARDIEN BLOQUE LE PORTAIL !", faceSmall, float64(w)/2, float64(h)*0.30, hex(0xff7070), text.AlignCenter)
				break
			}
		}
	}
	hint := "R RECHARGER · M MINIMAP · TAB CARTE · ÉCHAP PAUSE"
	if g.mode == modeCS {
		hint = "B BOUTIQUE · R RECHARGER · M MINIMAP · TAB CARTE · ÉCHAP PAUSE"
	} else if g.world.Player.Weapon == 3 {
		hint = "CLIC DROIT ZOOM · R RECHARGER · M MINIMAP · TAB CARTE · ÉCHAP PAUSE"
	}
	drawText(screen, hint, faceTiny, float64(w)/2, float64(h)-16, hex(0x6a7686), text.AlignCenter)
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
		label := "VAGUE ATTEINTE : %d  ·  SCORE : %d"
		if g.mode == modeCS {
			label = "MANCHE ATTEINTE : %d  ·  SCORE : %d"
		}
		drawText(screen, fmt.Sprintf(label, g.world.Wave, g.score), faceMed, w/2, h*0.52, hex(0xf2f2f4), text.AlignCenter)
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
	if g.mode == modeCS {
		drawText(screen, fmt.Sprintf("MANCHE %d  ·  SCORE %d  ·  $ %d", g.world.Wave, g.score, g.world.cash), faceMed, w/2, h*0.50, hex(0xf2f2f4), text.AlignCenter)
	} else if g.world.Survive {
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
	drawText(screen, "FPS RÉTRO · 5 ARMES · BOSS · MODE CS", faceMed, w/2, h*0.22, hex(0x8890a0), text.AlignCenter)

	opts := [3]string{"CAMPAGNE — 5 NIVEAUX", "SURVIE — VAGUES INFINIES", "COUNTER-STRIKE — BOUTIQUE"}
	for i, o := range opts {
		col := hex(0x8890a0)
		label := o
		if i == g.menuIdx {
			label = "» " + o + " «"
			col = hex(0xffe066)
		}
		drawText(screen, label, faceMed, w/2, h*0.32+float64(i)*36, col, text.AlignCenter)
	}

	lines := []string{
		"ZQSD/WASD : DÉPLACEMENT   ·   SOURIS+CLIC : VISER/TIRER   ·   MAJ : SPRINT",
		"1-5 + MOLETTE : ARMES   ·   CLIC DROIT : ZOOM SNIPER   ·   B : BOUTIQUE (MODE CS)",
		"M : MINIMAP   ·   TAB : GRANDE CARTE   ·   ÉCHAP : PAUSE",
		"FLÈCHES + ENTRÉE : CHOISIR LE MODE",
	}
	for i, l := range lines {
		drawText(screen, l, faceSmall, w/2, h*0.54+float64(i)*28, hex(0xc8d2dd), text.AlignCenter)
	}
	if int(math.Mod(g.time, 1.2)/0.6)%2 == 0 {
		drawText(screen, "ENTRÉE POUR JOUER", faceMed, w/2, h*0.82, hex(0xffe066), text.AlignCenter)
	}
}

// drawBuyMenu : boutique du mode counter-strike.
func (g *Game) drawBuyMenu(screen *ebiten.Image) {
	w := float64(InternalW * 2)
	h := float64(InternalH * 2)
	vector.DrawFilledRect(screen, 0, 0, float32(w), float32(h), rgbaSemi(4, 8, 6, 150), false)
	px, py, pw, ph := w/2-240, h/2-160, float64(480), float64(320)
	vector.DrawFilledRect(screen, float32(px), float32(py), float32(pw), float32(ph), rgbaSemi(8, 14, 11, 240), false)
	vector.StrokeRect(screen, float32(px), float32(py), float32(pw), float32(ph), 2, hex(0x39d7e0), false)
	drawText(screen, "BOUTIQUE", faceBig, w/2, py+18, hex(0x39d7e0), text.AlignCenter)
	drawText(screen, fmt.Sprintf("$ %d", g.world.cash), faceMed, w/2, py+62, hex(0xffe066), text.AlignCenter)

	items := []struct {
		txt   string
		price int
	}{}
	suffix := ""
	items = append(items,
		joined("1  FUSIL A POMPE", 300),
		joined("2  FUSIL A PLASMA", 700),
		joined("3  SNIPER", 900),
		joined("4  CANON D'IONS", 1200),
		joined("5  SOINS +40", 250),
	)
	items = append(items, struct {
		txt   string
		price int
	}{txt: "B  FERMER LA BOUTIQUE", price: 0})
	for i, it := range items {
		col := hex(0xc8d2dd)
		if it.price > 0 && g.world.cash < it.price {
			col = hex(0x6a4a44)
		}
		drawText(screen, it.txt+suffix, faceSmall, w/2, py+110+float64(i)*26, col, text.AlignCenter)
	}
}

func joined(text string, price int) struct {
	txt   string
	price int
} {
	return struct {
		txt   string
		price int
	}{fmt.Sprintf("%s — $%d", text, price), price}
}

func (g *Game) drawWinScreen(screen *ebiten.Image) {
	w := float64(InternalW * 2)
	h := float64(InternalH * 2)
	drawText(screen, "VICTOIRE !", faceTitle, w/2, h*0.30, hex(0x5fe8c0), text.AlignCenter)
	drawText(screen, fmt.Sprintf("SCORE FINAL : %d", g.score), faceBig, w/2, h*0.50, hex(0xf2f2f4), text.AlignCenter)
	drawText(screen, "ENTRÉE POUR REJOUER", faceSmall, w/2, h*0.66, hex(0xffe066), text.AlignCenter)
}
