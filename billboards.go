package main

import "github.com/hajimehoshi/ebiten/v2"

var (
	enemySprites map[byte][]*ebiten.Image
	enemyDead    map[byte]*ebiten.Image
	projSprite   *ebiten.Image
)

type billboard struct {
	x, y  float64
	img   *ebiten.Image
	scale float64
	vOff  float64
	glow  bool
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
			bbs = append(bbs, billboard{e.X, e.Y, deadImg, k.Scale, 0.2, false})
			continue
		}
		var img *ebiten.Image
		last := frames[len(frames)-1]
		switch {
		case e.Flash > 0:
			img = last
		case k.Ranged && e.ShootT > k.ShootCD-0.3:
			img = last
		case e.AtkCD > k.AtkCD-0.35:
			img = last
		default:
			img = frames[int(e.Frame*6)%2]
		}
		bbs = append(bbs, billboard{e.X, e.Y, img, k.Scale, 0.2, false})
	}
	for _, pk := range w.Pickups {
		img := heartSprite
		if pk.Kind == TileAmmo {
			img = ammoSprite
		}
		bbs = append(bbs, billboard{pk.X, pk.Y, img, 0.30, 0.52, true})
	}
	for _, pr := range w.Projs {
		bbs = append(bbs, billboard{pr.X, pr.Y, projSprite, 0.16, 0.12, true})
	}
	for _, ex := range w.Exits {
		bbs = append(bbs, billboard{ex.X, ex.Y, beamSprite, 1.5, -0.25, true})
	}
	return bbs
}
