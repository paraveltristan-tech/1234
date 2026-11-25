package main

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type mapRenderer struct {
	ox, oy, cell float32
	big          bool
}

func (g *Game) mapRect(screenW, screenH int) mapRenderer {
	w := g.world
	cell := 9.0
	if w.W > 24 || w.H > 18 {
		cell = 8.0
	}
	ox := float32(screenW) - float32(w.W)*float32(cell) - 16
	oy := 64
	if ox < 150 {
		ox = 150
	}
	return mapRenderer{ox: ox, oy: float32(oy), cell: float32(cell), big: false}
}

func (g *Game) bigMapRect(screenW, screenH int) mapRenderer {
	w := g.world
	cell := float64(screenH-200) / float64(maxInt(w.W, w.H))
	if cell > 26 {
		cell = 26
	}
	mw := float32(w.W) * float32(cell)
	mh := float32(w.H) * float32(cell)
	return mapRenderer{
		ox:   (float32(screenW) - mw) / 2,
		oy:   (float32(screenH)-mh)/2 + 20,
		cell: float32(cell),
		big:  true,
	}
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (r mapRenderer) drawTiles(g *Game, screen *ebiten.Image) {
	w := g.world
	bg := rgbaSemi(6, 10, 16, 200)
	if r.big {
		bg = rgbaSemi(4, 8, 14, 235)
	}
	vector.DrawFilledRect(screen, r.ox-4, r.oy-4, float32(w.W)*r.cell+8, float32(w.H)*r.cell+8, bg, false)
	for y := 0; y < w.H; y++ {
		for x := 0; x < w.W; x++ {
			c := w.Grid[y][x]
			var col color.Color = rgba(26, 34, 46, 150)
			if isSolid(c) {
				switch c {
				case TileBrick:
					col = hex(0x8a4634)
				case TileStone:
					col = hex(0x5a6774)
				case TileTech:
					col = hex(0x1f7c8c)
				}
			}
			vector.DrawFilledRect(screen, r.ox+float32(x)*r.cell, r.oy+float32(y)*r.cell, r.cell, r.cell, col, false)
		}
	}
}

func (r mapRenderer) drawEntities(g *Game, screen *ebiten.Image) {
	w := g.world
	pulse := uint8(160 + 80*math.Sin(g.time*4))
	for _, pk := range w.Pickups {
		col := rgba(220, 60, 70, 220)
		if pk.Kind == TileAmmo {
			col = rgba(255, 200, 60, 220)
		}
		r.dot(screen, pk.X, pk.Y, col)
	}
	if !w.bossAlive() {
		for _, ex := range w.Exits {
			r.dot(screen, ex.X, ex.Y, rgba(60, 235, 245, pulse))
		}
	}
	for _, e := range w.Enemies {
		if !e.Dead {
			if e.Kind == TileBoss {
				vector.DrawFilledCircle(screen,
					r.ox+float32(e.X)*r.cell, r.oy+float32(e.Y)*r.cell,
					r.cell*0.55, rgba(255, 40, 40, 240), false)
			} else {
				r.dot(screen, e.X, e.Y, rgba(255, 80, 80, 230))
			}
		}
	}
	p := &w.Player
	r.dot(screen, p.X, p.Y, rgbWhite())
	vector.StrokeLine(
		screen,
		r.ox+float32(p.X)*r.cell,
		r.oy+float32(p.Y)*r.cell,
		r.ox+float32(p.X+math.Cos(p.Dir)*1.6)*r.cell,
		r.oy+float32(p.Y+math.Sin(p.Dir)*1.6)*r.cell,
		2, rgbWhite(), false,
	)
	if r.big {
		vector.StrokeRect(screen, r.ox-4, r.oy-4, float32(w.W)*r.cell+8, float32(w.H)*r.cell+8, 2, hex(0x39d7e0), false)
	}
}

func (r mapRenderer) dot(screen *ebiten.Image, gx, gy float64, col color.Color) {
	vector.DrawFilledCircle(screen,
		r.ox+(float32(gx))*r.cell,
		r.oy+(float32(gy))*r.cell,
		r.cell*0.32, col, false)
}

func rgbWhite() color.Color {
	return color.RGBA{R: 240, G: 245, B: 250, A: 255}
}
