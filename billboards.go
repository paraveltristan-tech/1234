package main

import "github.com/hajimehoshi/ebiten/v2"

var (
	enemySprites   map[byte][]*ebiten.Image
	enemyDead      map[byte]*ebiten.Image
	projSprite     *ebiten.Image
	particleSprite *ebiten.Image
)

type billboard struct {
	x, y  float64
	img   *ebiten.Image
	scale float64
	vOff  float64
	glow  bool
	col   [3]float32
}

func (g *Game) collectBillboards() []billboard {
	var bbs []billboard
	w := g.world
	for _, e := range w.Enemies {
		k, ok := enemyKinds[e.Kind]
		if !ok {
			k = enemyKinds[TileEnemy]
		}
		frames := enemySprites[e.Kind]
		if len(frames) == 0 {
			frames = enemySprites[TileEnemy]
		}
		deadImg := enemyDead[e.Kind]
		if deadImg == nil {
			deadImg = enemyDead[TileEnemy]
		}
		if e.Dead {
			bbs = append(bbs, billboard{x: e.X, y: e.Y, img: deadImg, scale: k.Scale, vOff: 0.2})
			continue
		}
		var img *ebiten.Image
		last := frames[len(frames)-1]
		switch {
		case e.Flash > 0:
			img = last
		case (k.Ranged || k.Boss) && e.ShootT > k.ShootCD-0.3:
			img = last
		case e.AtkCD > k.AtkCD-0.35:
			img = last
		default:
			img = frames[int(e.Frame*6)%2]
		}
		bbs = append(bbs, billboard{x: e.X, y: e.Y, img: img, scale: k.Scale, vOff: 0.2})
	}
	for _, pk := range w.Pickups {
		img := heartSprite
		if pk.Kind == TileAmmo {
			img = ammoSprite
		}
		bbs = append(bbs, billboard{x: pk.X, y: pk.Y, img: img, scale: 0.30, vOff: 0.52, glow: true})
	}
	for _, pr := range w.Projs {
		bbs = append(bbs, billboard{x: pr.X, y: pr.Y, img: projSprite, scale: 0.16, vOff: 0.12, glow: true})
	}
	for _, ex := range w.Exits {
		bbs = append(bbs, billboard{x: ex.X, y: ex.Y, img: beamSprite, scale: 1.5, vOff: -0.25, glow: true})
	}
	if particleSprite != nil {
		for _, pa := range w.Parts {
			fade := 1 - pa.T/pa.Life
			bbs = append(bbs, billboard{
				x: pa.X, y: pa.Y, img: particleSprite,
				scale: pa.Size * (0.4 + 0.6*fade),
				vOff:  0.5 - pa.Z,
				glow:  pa.Glow, col: pa.Col,
			})
		}
	}
	return bbs
}
