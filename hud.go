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

// finitions d'armes achetables dans la boutique du mode counter-strike.
var skinNames = []string{"STANDARD", "CHROME", "OR", "NEON", "INFERNO", "TOXIQUE"}
var skinPrices = []int{0, 400, 700, 900, 900, 900}

// applySkin reteint la couleur de base selon le finition equipee.
// skin 0 (standard) conserve la couleur d'origine.
func applySkin(c uint32, skin int) uint32 {
	if skin <= 0 || skin >= len(skinNames) {
		return c
	}
	r := float64((c>>16)&0xff) / 255
	gr := float64((c>>8)&0xff) / 255
	b := float64(c&0xff) / 255
	switch skin {
	case 1: // chrome : desature, reflet bleute
		a := (r + gr + b) / 3
		r, gr, b = a*0.85, a*0.93, clamp1(a*1.15)
	case 2: // or brosse
		l := 0.35*r + 0.5*gr + 0.15*b
		r, gr, b = clamp1(l*1.55+0.12), clamp1(l*1.15), l*0.30
	case 3: // neon : corps sombre, accents cyan
		r, gr, b = r*0.30, gr*0.60, clamp1(b*1.45+0.10)
	case 4: // inferno : acier brule, accents orange
		r, gr, b = clamp1(r*1.45+0.06), gr*0.42, b*0.30
	case 5: // toxique : vert acide
		r, gr, b = r*0.35, clamp1(gr*1.45+0.05), b*0.38
	}
	return uint32(r*255)<<16 | uint32(gr*255)<<8 | uint32(b*255)
}

// clamp1 borne une valeur flottante entre 0 et 1.
func clamp1(v float64) float64 {
	if v > 1 {
		return 1
	}
	return v
}

// sk applique le finition de l'arme equipee a une couleur de base.
func (g *Game) sk(c uint32) uint32 {
	p := &g.world.Player
	return applySkin(c, p.Skins[p.Weapon])
}

// drawWeapon : vue premiere personne de l'arme equipee, dessinee pas a pas
// avec passager de details (glissiere, lunette, serpentin, bandes danger...)
// et animation de rechargement (plongeon + balancement).
func (g *Game) drawWeapon() {
	p := &g.world.Player
	if g.zoomed {
		return
	}
	if g.world.VehMounted {
		if g.world.VehKind == 'T' {
			g.drawTankView()
		}
		if g.world.VehKind == 'H' {
			g.drawHeliView()
		}
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
	reloadSway := 0.0
	if p.ReloadT > 0 {
		prog := 1 - p.ReloadT/max(0.01, p.ReloadDur)
		reloadDip = math.Sin(prog*math.Pi) * 54
		reloadSway = math.Sin(prog*math.Pi*2) * 5
	}

	gx := float64(InternalW)/2 + 18 + bobX + reloadSway
	gy := float64(InternalH) - 78 + bobY + recoil + reloadDip
	fx, fy := float32(gx), float32(gy)
	pump := float32(recoilPct * 6)

	switch p.Weapon {
	case 1: // fusil a pompe : canon nervure + guidon, culasse acier, crosse noyer veinee, pompe moletee
		vector.DrawFilledRect(screenBuf, fx-22, fy-8, 64, 13, hex(g.sk(0x2a313b)), false)
		vector.DrawFilledRect(screenBuf, fx-22, fy-8, 64, 3, hex(g.sk(0x46525f)), false)
		vector.DrawFilledRect(screenBuf, fx-22, fy+1, 64, 2, hex(g.sk(0x151a20)), false)
		for i := 0; i < 8; i++ {
			vector.DrawFilledRect(screenBuf, fx-16+float32(i*7), fy+3, 3, 3, hex(g.sk(0x10141a)), false)
		}
		vector.DrawFilledCircle(screenBuf, fx-26, fy-5, 3, hex(g.sk(0xc8ccd2)), false)
		vector.DrawFilledRect(screenBuf, fx-20, fy+5, 10, 15, hex(g.sk(0x11161b)), false)
		vector.DrawFilledRect(screenBuf, fx-18, fy+7, 6, 4, hex(g.sk(0x2c343e)), false)
		vector.DrawFilledRect(screenBuf, fx-16+pump, fy+6+pump*0.7, 26, 10, hex(g.sk(0x4a3423)), false)
		vector.DrawFilledRect(screenBuf, fx-16+pump, fy+9+pump*0.7, 26, 3, hex(g.sk(0x6d5238)), false)
		for i := 0; i < 3; i++ {
			vector.DrawFilledRect(screenBuf, fx-14+pump, fy+8+pump*0.7+float32(2+i*2), 22, 1, hex(g.sk(0x332619)), false)
		}
		vector.DrawFilledRect(screenBuf, fx+8, fy+10, 30, 36, hex(g.sk(0x4a3423)), false)
		vector.DrawFilledRect(screenBuf, fx+8, fy+10, 30, 5, hex(g.sk(0x5d4430)), false)
		vector.DrawFilledRect(screenBuf, fx+10, fy+22, 26, 3, hex(g.sk(0x3d2c1e)), false)
		vector.DrawFilledRect(screenBuf, fx+12, fy+32, 22, 3, hex(g.sk(0x5d4430)), false)
	case 2: // fusil a plasma : rail, cellules pulsees, serpentin, bandes danger, crosse gaine
		vector.DrawFilledRect(screenBuf, fx-10, fy-12, 52, 20, hex(g.sk(0x1c3a40)), false)
		vector.DrawFilledRect(screenBuf, fx-10, fy-12, 52, 3, hex(g.sk(0x2c5a63)), false)
		for i := 0; i < 6; i++ {
			vector.DrawFilledRect(screenBuf, fx-6+float32(i*8), fy-14, 4, 2, hex(g.sk(0x39434f)), false)
		}
		vector.DrawFilledRect(screenBuf, fx-8, fy-6, 46, 6, hex(g.sk(0x0d2f36)), false)
		pulse1 := 3 + int(2*math.Sin(p.Bob*2.6))
		for i := 0; i < 3; i++ {
			vector.DrawFilledCircle(screenBuf, float32(6+i*15)+fx-2, fy-3, float32(pulse1), hex(g.sk(0x39d7e0)), false)
			vector.DrawFilledCircle(screenBuf, float32(6+i*15)+fx-2, fy-3, float32(pulse1-1), hex(g.sk(0xd0faff)), false)
		}
		vector.DrawFilledRect(screenBuf, fx+2, fy+8, 34, 32, hex(g.sk(0x24363d)), false)
		vector.DrawFilledRect(screenBuf, fx+2, fy+8, 34, 4, hex(g.sk(0x39546a)), false)
		vector.DrawFilledRect(screenBuf, fx-4, fy+14, 4, 4, hex(g.sk(0xffc93a)), false)
		vector.DrawFilledRect(screenBuf, fx, fy+14, 4, 4, hex(g.sk(0x1a1c14)), false)
		vector.DrawFilledRect(screenBuf, fx-4, fy+20, 4, 4, hex(g.sk(0xffc93a)), false)
		vector.DrawFilledRect(screenBuf, fx, fy+20, 4, 4, hex(g.sk(0x1a1c14)), false)
		vector.DrawFilledRect(screenBuf, fx, fy+40, 26, 26, hex(g.sk(0x1d232b)), false)
		vector.DrawFilledRect(screenBuf, fx, fy+54, 26, 3, hex(g.sk(0x39d7e0)), false)
		vector.DrawFilledRect(screenBuf, fx+28, fy+44, 8, 12, hex(g.sk(0x14181d)), false)
	case 3: // sniper : canon long + frein de bouche, lunette sur montages, culasse a poignee, chargeur incline
		vector.DrawFilledRect(screenBuf, fx-40, fy-4, 84, 11, hex(g.sk(0x1c2630)), false)
		vector.DrawFilledRect(screenBuf, fx-40, fy-4, 84, 3, hex(g.sk(0x2f3f4c)), false)
		vector.DrawFilledRect(screenBuf, fx-40, fy+4, 84, 2, hex(g.sk(0x39d7e0)), false)
		for i := 0; i < 3; i++ {
			vector.DrawFilledRect(screenBuf, fx-46+float32(i*4), fy-2, 2, 7, hex(g.sk(0x0d1116)), false)
		}
		vector.DrawFilledRect(screenBuf, fx-30, fy-14, 34, 8, hex(g.sk(0x0e1116)), false)
		vector.DrawFilledRect(screenBuf, fx-26, fy-12, 26, 4, hex(g.sk(0x151c24)), false)
		vector.DrawFilledRect(screenBuf, fx-14, fy-6, 6, 4, hex(g.sk(0x0e1116)), false)
		vector.DrawFilledRect(screenBuf, fx+12, fy-6, 6, 4, hex(g.sk(0x0e1116)), false)
		vector.DrawFilledCircle(screenBuf, fx-30, fy-10, 4, hex(g.sk(0x59c8ff)), false)
		vector.DrawFilledRect(screenBuf, fx+6, fy+1, 4, 8, hex(g.sk(0x3a4149)), false)
		vector.DrawFilledCircle(screenBuf, fx+8, fy+8, 5, hex(g.sk(0x3a4149)), false)
		vector.DrawFilledRect(screenBuf, fx+14, fy+7, 12, 16, hex(g.sk(0x1a2027)), false)
		vector.DrawFilledRect(screenBuf, fx+16, fy+23, 8, 6, hex(g.sk(0x12171d)), false)
		vector.DrawFilledRect(screenBuf, fx+4, fy+18, 24, 26, hex(g.sk(0x22272e)), false)
		vector.DrawFilledRect(screenBuf, fx+6, fy+20, 20, 5, hex(g.sk(0x333c45)), false)
		vector.DrawFilledRect(screenBuf, fx-34, fy+8, 74, 1, hex(g.sk(0x59c8ff)), false)
	case 4: // canon d'ions : reservoir a jauge de charge, serpentin double, orbe, bandes danger
		vector.DrawFilledRect(screenBuf, fx-14, fy-8, 64, 18, hex(g.sk(0x1a2620)), false)
		vector.DrawFilledRect(screenBuf, fx-14, fy-8, 64, 3, hex(g.sk(0x2c4438)), false)
		vector.DrawFilledRect(screenBuf, fx-14, fy+4, 64, 4, hex(g.sk(0x0c1210)), false)
		for i := 0; i < 4; i++ {
			vector.DrawFilledCircle(screenBuf, float32(4+i*13)+fx-10, fy-1, 5, hex(g.sk(0x0f3326)), false)
			vector.DrawFilledCircle(screenBuf, float32(4+i*13)+fx-10, fy-1, 2, hex(g.sk(0x2c8a5e)), false)
		}
		vector.DrawFilledRect(screenBuf, fx-8, fy+14, 4, 4, hex(g.sk(0xffc93a)), false)
		vector.DrawFilledRect(screenBuf, fx-4, fy+14, 4, 4, hex(g.sk(0x1a1c14)), false)
		vector.DrawFilledRect(screenBuf, fx-8, fy+20, 4, 4, hex(g.sk(0xffc93a)), false)
		vector.DrawFilledRect(screenBuf, fx-4, fy+20, 4, 4, hex(g.sk(0x1a1c14)), false)
		pulse := 6 + int(3*math.Sin(p.Bob*1.4))
		vector.DrawFilledCircle(screenBuf, fx+44, fy+2, float32(pulse+6), hex(g.sk(0x0f5e46)), false)
		vector.DrawFilledCircle(screenBuf, fx+44, fy+2, float32(pulse+3), hex(g.sk(0x1c8a5e)), false)
		vector.DrawFilledCircle(screenBuf, fx+44, fy+2, float32(pulse), hex(g.sk(0x5fe8a0)), false)
		vector.DrawFilledRect(screenBuf, fx, fy+10, 30, 34, hex(g.sk(0x232b28)), false)
		vector.DrawFilledRect(screenBuf, fx+2, fy+12, 26, 4, hex(g.sk(0x35433d)), false)
		for i := 0; i < 4; i++ { // jauge de charge du reservoir = munitions restantes
			pip := hex(g.sk(0x1a2a24))
			if i < p.Ammo[4] {
				pip = hex(g.sk(0x5fe8a0))
			}
			vector.DrawFilledRect(screenBuf, fx+6, fy+38-float32(i*5), 12, 3, pip, false)
		}
	default: // pistolet : marteau, glissiere 3 tons, rayures, port d'ejection, guidon, poignee quadrillee
		vector.DrawFilledRect(screenBuf, fx-4, fy-16, 12, 6, hex(g.sk(0x14181d)), false)
		vector.DrawFilledRect(screenBuf, fx-2, fy-12, 38, 16, hex(g.sk(0x272e38)), false)
		vector.DrawFilledRect(screenBuf, fx-2, fy-12, 38, 4, hex(g.sk(0x46525f)), false)
		vector.DrawFilledRect(screenBuf, fx-2, fy-2, 38, 4, hex(g.sk(0x1a1f26)), false)
		for i := 0; i < 5; i++ {
			vector.DrawFilledRect(screenBuf, fx+float32(3+i*6), fy-9, 2, 5, hex(g.sk(0x10141a)), false)
		}
		vector.DrawFilledRect(screenBuf, fx+20, fy-7, 8, 6, hex(g.sk(0x12161c)), false)
		vector.DrawFilledRect(screenBuf, fx-6, fy-11, 4, 4, hex(g.sk(0xd0f6ff)), false)
		vector.DrawFilledRect(screenBuf, fx-2, fy+4, 30, 8, hex(g.sk(0x1d232b)), false)
		vector.DrawFilledRect(screenBuf, fx-2, fy+12, 24, 34, hex(g.sk(0x2b313a)), false)
		vector.DrawFilledRect(screenBuf, fx-2, fy+12, 24, 3, hex(g.sk(0x3d4650)), false)
		for i := 0; i < 4; i++ {
			vector.DrawFilledRect(screenBuf, fx-1, fy+18+float32(i*7), 21, 2, hex(g.sk(0x171c22)), false)
		}
	}

	if p.ReloadT > 0 {
		prog := 1 - p.ReloadT/max(0.01, p.ReloadDur)
		drop := math.Sin(prog*math.Pi) * 26
		mx, my := fx+6, fy+16+float32(drop)
		vector.DrawFilledRect(screenBuf, mx, my, 10, 16, hex(g.sk(0x22272e)), false)
		vector.DrawFilledRect(screenBuf, mx, my, 10, 3, hex(g.sk(0x39434f)), false)
		vector.DrawFilledRect(screenBuf, mx+2, my+6, 6, 2, hex(g.sk(0xffc93a)), false)
	}

	if g.world.Muzzle > 0 {
		switch p.Weapon {
		case 1:
			vector.DrawFilledCircle(screenBuf, fx-30, fy-2, 10, hex(0xffe08a), false)
			vector.DrawFilledCircle(screenBuf, fx-30, fy-2, 5, hex(0xfff6cf), false)
			vector.DrawFilledRect(screenBuf, fx-42, fy-3, 12, 2, hex(0xffe08a), false)
		case 2:
			vector.DrawFilledCircle(screenBuf, fx+18, fy-20, 11, hex(0x5fe8ff), false)
			vector.DrawFilledCircle(screenBuf, fx+18, fy-20, 6, hex(0xd0faff), false)
		case 3:
			vector.DrawFilledCircle(screenBuf, fx-36, fy+1, 10, hex(0xfffbe8), false)
			vector.DrawFilledCircle(screenBuf, fx-36, fy+1, 5, hex(0xd0f6ff), false)
			vector.DrawFilledRect(screenBuf, fx-50, fy, 14, 2, hex(0xfffbe8), false)
		case 4:
			vector.DrawFilledCircle(screenBuf, fx+44, fy-12, 15, hex(0x5fe8a0), false)
			vector.DrawFilledCircle(screenBuf, fx+44, fy-12, 8, hex(0xd0ffe8), false)
		default:
			vector.DrawFilledCircle(screenBuf, fx-6, fy-14, 9, hex(0xffe08a), false)
			vector.DrawFilledCircle(screenBuf, fx-6, fy-14, 4, hex(0xfff6cf), false)
		}
	}
}

// drawScope : overlay de lunette quand le sniper est armé et zoomé.
// drawTankView dessine la vue pilote du tank : canon central qui recule,
// frein de bouche et capot de tourelle en bas d'ecran.
func (g *Game) drawTankView() {
	w := float32(InternalW)
	h := float32(InternalH)
	cx := w / 2
	rec := float32(g.world.TankCd/0.9) * 12
	vector.DrawFilledRect(screenBuf, cx-10, -30+rec, 20, h*0.62, hex(0x4a4a30), false)
	vector.DrawFilledRect(screenBuf, cx-10, -30+rec, 5, h*0.62, hex(0x60603c), false)
	vector.DrawFilledRect(screenBuf, cx-10, -30+rec, 20, 5, hex(0x60603c), false)
	vector.DrawFilledRect(screenBuf, cx-15, h*0.12+rec, 30, 9, hex(0x2e3338), false)
	vector.DrawFilledRect(screenBuf, cx-60, h-58, 120, 58, hex(0x5c5a34), false)
	vector.DrawFilledRect(screenBuf, cx-60, h-58, 120, 7, hex(0x7a7846), false)
	vector.DrawFilledRect(screenBuf, cx-20, h-48, 40, 12, hex(0x34341f), false)
}

// drawHeliView dessine la verriere de l'helicoptere : cadre de bulle,
// plafond sombre et plancher de cabine.
func (g *Game) drawHeliView() {
	w := float32(InternalW)
	h := float32(InternalH)
	vector.DrawFilledRect(screenBuf, 0, 0, w, 24, hex(0x0c1410), false)
	vector.DrawFilledRect(screenBuf, 0, 24, w, 3, hex(0x37543f), false)
	vector.DrawFilledRect(screenBuf, 0, h-38, w, 38, hex(0x18231c), false)
	vector.DrawFilledRect(screenBuf, 0, h-40, w, 3, hex(0x37543f), false)
}

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
	if p.BurnT > 0 {
		drawText(screen, "EN FEU !", faceSmall, float64(w)/2, float64(h)*0.44, hex(0xff7a2a), text.AlignCenter)
	}
	hint := "N SON · R RECHARGER · M MINIMAP · TAB CARTE · ÉCHAP PAUSE"
	if g.world.Vehicle.X > 0 {
		hint = "E VTT · " + hint
	}
	if g.world.Tank.X > 0 {
		hint = "E TANK · " + hint
	}
	if g.world.Heli.X > 0 {
		hint = "E HELI · " + hint
	}
	if g.world.VehMounted && g.world.VehKind == 'H' {
		hint = "CLIC ROQUETTES · E DESCENDRE"
	}
	if g.world.VehMounted && g.world.VehKind == 'T' {
		hint = "CLIC TIR CANON · E DESCENDRE"
	}
	if g.mode == modeCS {
		hint = "B BOUTIQUE · N SON · R RECHARGER · M MINIMAP · TAB CARTE · ÉCHAP PAUSE"
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
	px, py, pw, ph := w/2-250, h/2-170, float64(500), float64(340)
	vector.DrawFilledRect(screen, float32(px), float32(py), float32(pw), float32(ph), rgbaSemi(8, 14, 11, 240), false)
	vector.StrokeRect(screen, float32(px), float32(py), float32(pw), float32(ph), 2, hex(0x39d7e0), false)
	title := "BOUTIQUE - MUNITIONS"
	if g.buyPage == 1 {
		title = "BOUTIQUE - SKINS D'ARMES"
	}
	drawText(screen, title, faceBig, w/2, py+16, hex(0x39d7e0), text.AlignCenter)
	drawText(screen, fmt.Sprintf("$ %d", g.world.cash), faceMed, w/2, py+56, hex(0xffe066), text.AlignCenter)

	type line struct {
		txt  string
		lock bool
	}
	lines := []line{}
	p := &g.world.Player
	if g.buyPage == 0 {
		for _, it := range []struct {
			txt   string
			price int
		}{
			{"1  FUSIL A POMPE", 300},
			{"2  FUSIL A PLASMA", 700},
			{"3  SNIPER", 900},
			{"4  CANON D'IONS", 1200},
			{"5  SOINS +40", 250},
		} {
			lines = append(lines, line{txt: fmt.Sprintf("%s — $%d", it.txt, it.price), lock: g.world.cash < it.price})
		}
		lines = append(lines, line{txt: "TAB  PAGE SKINS D'ARMES", lock: false})
	} else {
		for wi := range weapons {
			cur := skinNames[p.Skins[wi]]
			if p.Skins[wi] >= len(skinNames)-1 {
				lines = append(lines, line{txt: fmt.Sprintf("%d  %s : %s (MAX)", wi+1, weapons[wi].Name, cur), lock: false})
				continue
			}
			price := skinPrices[p.Skins[wi]+1]
			next := skinNames[p.Skins[wi]+1]
			txt := fmt.Sprintf("%d  %s : %s vers %s ($ %d)", wi+1, weapons[wi].Name, cur, next, price)
			lines = append(lines, line{txt: txt, lock: g.world.cash < price})
		}
		lines = append(lines, line{txt: "TAB  PAGE MUNITIONS", lock: false})
	}
	lines = append(lines, line{txt: "B  FERMER LA BOUTIQUE", lock: false})
	for i, l := range lines {
		col := hex(0xc8d2dd)
		if l.lock {
			col = hex(0x6a4a44)
		}
		drawText(screen, l.txt, faceSmall, w/2, py+104+float64(i)*26, col, text.AlignCenter)
	}
}

func (g *Game) drawWinScreen(screen *ebiten.Image) {
	w := float64(InternalW * 2)
	h := float64(InternalH * 2)
	drawText(screen, "VICTOIRE !", faceTitle, w/2, h*0.30, hex(0x5fe8c0), text.AlignCenter)
	drawText(screen, fmt.Sprintf("SCORE FINAL : %d", g.score), faceBig, w/2, h*0.50, hex(0xf2f2f4), text.AlignCenter)
	drawText(screen, "ENTRÉE POUR REJOUER", faceSmall, w/2, h*0.66, hex(0xffe066), text.AlignCenter)
}
