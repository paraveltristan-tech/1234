package main

import (
	"image"
	"math"
	"sort"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var (
	wallTexSets map[byte][]*ebiten.Image
	brickTex    *image.RGBA
	stoneTex    *image.RGBA
	techTex     *image.RGBA
	heartSprite *ebiten.Image
	ammoSprite  *ebiten.Image
	beamSprite  *ebiten.Image
	background  *ebiten.Image
	screenBuf   *ebiten.Image
	curLevel    int
)

// rebuildFogSets regenerates the wall texture shade/fog ramp for the level.
// Skipped when no texture assets exist yet (headless tests).
func rebuildFogSets() {
	if brickTex == nil || stoneTex == nil || techTex == nil {
		return
	}
	pal := paletteFor(curLevel)
	wallTexSets = map[byte][]*ebiten.Image{
		TileBrick: makeFogVariants(brickTex, pal.Fog),
		TileStone: makeFogVariants(stoneTex, pal.Fog),
		TileTech:  makeFogVariants(techTex, pal.Fog),
	}
}

const texSize = 64
const maxWallDist = 14.0

func shadeFor(dist float64, side int) float64 {
	sh := 1.0 - dist/maxWallDist
	if sh < 0.12 {
		sh = 0.12
	}
	if side == 1 {
		sh *= 0.78
	}
	return sh
}

func lerp3(a, b [3]uint8, t float64) [3]uint8 {
	return [3]uint8{
		clamp8(int(float64(a[0]) + (float64(b[0])-float64(a[0]))*t)),
		clamp8(int(float64(a[1]) + (float64(b[1])-float64(a[1]))*t)),
		clamp8(int(float64(a[2]) + (float64(b[2])-float64(a[2]))*t)),
	}
}

func initRender() {
	screenBuf = ebiten.NewImage(InternalW, InternalH)
}

func (g *Game) Draw(screen *ebiten.Image) {
	switch g.state {
	case stateTitle:
		screen.DrawImage(background, nil)
		g.drawTitle(screen)
	case stateDead:
		g.drawWorldFrame()
		screen.DrawImage(screenBuf, scaled())
		g.drawOverlays(screen)
		g.drawDeathScreen(screen)
	case stateWin:
		screen.DrawImage(background, nil)
		g.drawWinScreen(screen)
	default:
		g.drawWorldFrame()
		g.drawWeapon()
		screen.DrawImage(screenBuf, scaled())
		g.drawOverlays(screen)
		g.drawHUD(screen)
		if g.zoomed {
			g.drawScope(screen)
		}
		g.drawMaps(screen)
		if g.buyOpen {
			g.drawBuyMenu(screen)
		}
		if g.state == statePause {
			g.drawPauseScreen(screen)
		}
	}
}

func (g *Game) drawMaps(screen *ebiten.Image) {
	if g.world == nil {
		return
	}
	sw, sh := int(InternalW*ScaleF), int(InternalH*ScaleF)
	if g.ShowMinimap && !g.ShowBigMap {
		if r := g.mapRect(sw, sh); r.cell > 0 {
			r.drawTiles(g, screen)
			r.drawEntities(g, screen)
		}
	}
	if g.ShowBigMap {
		r := g.bigMapRect(sw, sh)
		r.drawTiles(g, screen)
		r.drawEntities(g, screen)
		if g.world.Survive {
			label := "SURVIE · VAGUE " + itoa(g.world.Wave)
			if g.mode == modeCS {
				label = "COUNTER-STRIKE · MANCHE " + itoa(g.world.Wave)
			}
			drawText(screen, "CARTE — "+label, faceMed, float64(sw)/2, 40, hex(0x5fe8ff), text.AlignCenter)
		} else {
			drawText(screen, "CARTE — NIVEAU "+itoa(g.level), faceMed, float64(sw)/2, 40, hex(0x5fe8ff), text.AlignCenter)
		}
	}
}

// fovNow : le zoom du sniper reduit la moitie du champ de vision.
func (g *Game) fovNow() float64 {
	if g.zoomed {
		return 0.24
	}
	return 0.66
}

func scaled() *ebiten.DrawImageOptions {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(ScaleF, ScaleF)
	return op
}

func (g *Game) drawWorldFrame() {
	w := g.world
	if w == nil {
		return
	}
	p := &w.Player
	pal := paletteFor(curLevel)

	bands := InternalH / 12
	for i := 0; i < bands; i++ {
		t := float64(i) / float64(bands-1)
		cr := lerp3(pal.CeilTop, pal.CeilBottom, t)
		fr := lerp3(pal.FloorNear, pal.FloorFar, t)
		vector.DrawFilledRect(screenBuf, 0, float32(i*6), InternalW, 6, rgba(cr[0], cr[1], cr[2], 255), false)
		vector.DrawFilledRect(screenBuf, 0, float32(InternalH-(i+1)*6), InternalW, 6, rgba(fr[0], fr[1], fr[2], 255), false)
	}

	dirX, dirY := math.Cos(p.Dir), math.Sin(p.Dir)
	fov := g.fovNow()
	planeX, planeY := -dirY*fov, dirX*fov
	zbuf := make([]float64, InternalW)

	for x := 0; x < InternalW; x++ {
		cameraX := 2*float64(x)/InternalW - 1
		rdx := dirX + planeX*cameraX
		rdy := dirY + planeY*cameraX
		dist, side, tile, wallX := castStep(w.Grid, w.W, w.H, p.X, p.Y, rdx, rdy)
		zbuf[x] = dist

		lineH := int(float64(InternalH) / dist)
		if lineH > InternalH*6 {
			lineH = InternalH * 6
		}
		drawStart := -lineH/2 + InternalH/2
		set := wallTexSets[tile]
		if set == nil {
			set = wallTexSets[TileBrick]
		}
		shFrac := 1.0 - dist/maxWallDist
		if shFrac < 0.02 {
			shFrac = 0.02
		}
		bucket := int(shFrac * 7.999)
		tex := set[bucket]
		texX := int(wallX * texSize)
		if (side == 0 && rdx > 0) || (side == 1 && rdy < 0) {
			texX = texSize - texX - 1
		}
		if texX < 0 {
			texX = 0
		}
		src := tex.SubImage(image.Rect(texX, 0, texX+1, texSize)).(*ebiten.Image)
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(1, float64(lineH)/texSize)
		op.GeoM.Translate(float64(x), float64(drawStart))
		if side == 1 {
			op.ColorScale.Scale(0.72, 0.72, 0.78, 1)
		}
		screenBuf.DrawImage(src, op)
	}

	g.drawBillboards(zbuf)
}

func (g *Game) drawBillboards(zbuf []float64) {
	w := g.world
	p := &w.Player
	dirX, dirY := math.Cos(p.Dir), math.Sin(p.Dir)
	fov := g.fovNow()
	planeX, planeY := -dirY*fov, dirX*fov

	bbs := g.collectBillboards()
	sort.SliceStable(bbs, func(i, j int) bool {
		di := (bbs[i].x-p.X)*(bbs[i].x-p.X) + (bbs[i].y-p.Y)*(bbs[i].y-p.Y)
		dj := (bbs[j].x-p.X)*(bbs[j].x-p.X) + (bbs[j].y-p.Y)*(bbs[j].y-p.Y)
		return di > dj
	})

	for _, bb := range bbs {
		sx, sy := bb.x-p.X, bb.y-p.Y
		invDet := 1.0 / (planeX*dirY - dirX*planeY)
		tX := invDet * (dirY*sx - dirX*sy)
		tY := invDet * (-planeY*sx + planeX*sy)
		if tY < 0.08 {
			continue
		}
		imgW := bb.img.Bounds().Dx()
		imgH := bb.img.Bounds().Dy()
		spriteH := int(math.Abs(float64(InternalH)/tY) * bb.scale)
		spriteW := spriteH * imgW / imgH
		if spriteW <= 0 || spriteH <= 0 {
			continue
		}
		screenX := int(float64(InternalW) / 2 * (1 + tX/tY))
		vMoveScreen := int(bb.vOff * float64(InternalH) / tY)
		drawStartY := -spriteH/2 + InternalH/2 + vMoveScreen

		dist := math.Hypot(sx, sy)
		sh := float32(shadeFor(dist, 0))
		if bb.glow {
			sh = 1
		}

		for stripe := -spriteW / 2; stripe < spriteW-spriteW/2; stripe++ {
			x := screenX + stripe
			if x < 0 || x >= InternalW {
				continue
			}
			if tY >= zbuf[x] {
				continue
			}
			texX := (stripe + spriteW/2) * imgW / spriteW
			if texX < 0 {
				texX = 0
			}
			if texX >= imgW {
				texX = imgW - 1
			}
			src := bb.img.SubImage(image.Rect(texX, 0, texX+1, imgH)).(*ebiten.Image)
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Scale(1, float64(spriteH)/float64(imgH))
			op.GeoM.Translate(float64(x), float64(drawStartY))
			tint := bb.col
			if tint == ([3]float32{}) {
				tint = [3]float32{1, 1, 1}
			}
			op.ColorScale.Scale(tint[0]*sh, tint[1]*sh, tint[2]*sh, 1)
			if bb.glow {
				op.CompositeMode = ebiten.CompositeModeLighter
			}
			screenBuf.DrawImage(src, op)
		}
	}
}
