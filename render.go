package main

import (
	"image"
	"math"
	"sort"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var (
	wallTex     map[byte]*ebiten.Image
	impFrames   []*ebiten.Image
	impDead     *ebiten.Image
	heartSprite *ebiten.Image
	ammoSprite  *ebiten.Image
	beamSprite  *ebiten.Image
	background  *ebiten.Image
	screenBuf   *ebiten.Image
)

const texSize = 64
const maxWallDist = 14.0

func shadeFor(dist float64, side int) float64 {
	sh := 1.0 - dist/maxWallDist
	if sh < 0.12 {
		sh = 0.12
	}
	if side == 1 {
		sh *= 0.72
	}
	return sh
}

func initRender() {
	screenBuf = ebiten.NewImage(InternalW, InternalH)
}

type billboard struct {
	x, y  float64
	img   *ebiten.Image
	scale float64
	vOff  float64
	glow  bool
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
	}
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

	for i := 0; i < InternalH/2; i += 6 {
		a := float64(i) / 4
		vector.DrawFilledRect(screenBuf, 0, float32(i), InternalW, 6, rgba(14+uint8(a/3), 18+uint8(a/2), 30+uint8(a/2), 255), false)
		fy := float32(InternalH/2 + i)
		vector.DrawFilledRect(screenBuf, 0, fy, InternalW, 6, rgba(36+uint8(a), 30+uint8(a), 28+uint8(a), 255), false)
	}

	dirX, dirY := math.Cos(p.Dir), math.Sin(p.Dir)
	fov := 0.66
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
		tex := wallTex[tile]
		if tex == nil {
			tex = wallTex[TileBrick]
		}
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
		sh := shadeFor(dist, side)
		op.ColorScale.Scale(float32(sh), float32(sh), float32(sh), 1)
		screenBuf.DrawImage(src, op)
	}

	g.drawBillboards(zbuf)
}

func (g *Game) collectBillboards() []billboard {
	var bbs []billboard
	w := g.world
	for _, e := range w.Enemies {
		if e.Dead {
			bbs = append(bbs, billboard{e.X, e.Y, impDead, 0.8, 0.2, false})
			continue
		}
		var img *ebiten.Image
		switch {
		case e.Flash > 0:
			img = impFrames[2]
		case e.AtkCD > enemyAtkCD-0.35:
			img = impFrames[2]
		default:
			img = impFrames[int(e.Frame*6)%2]
		}
		bbs = append(bbs, billboard{e.X, e.Y, img, 0.8, 0.2, false})
	}
	for _, pk := range w.Pickups {
		img := heartSprite
		if pk.Kind == TileAmmo {
			img = ammoSprite
		}
		bbs = append(bbs, billboard{pk.X, pk.Y, img, 0.30, 0.52, true})
	}
	for _, ex := range w.Exits {
		bbs = append(bbs, billboard{ex.X, ex.Y, beamSprite, 1.5, -0.25, true})
	}
	return bbs
}

func (g *Game) drawBillboards(zbuf []float64) {
	w := g.world
	p := &w.Player
	dirX, dirY := math.Cos(p.Dir), math.Sin(p.Dir)
	fov := 0.66
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
			op.ColorScale.Scale(sh, sh, sh, 1)
			if bb.glow {
				op.CompositeMode = ebiten.CompositeModeLighter
			}
			screenBuf.DrawImage(src, op)
		}
	}
}
