package main

import (
	"bytes"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font/gofont/goregular"
)

func hex(c uint32) color.Color {
	return color.RGBA{R: uint8(c >> 16), G: uint8(c >> 8), B: uint8(c), A: 0xff}
}

func rgba(r, g, b, a uint8) color.Color {
	return color.RGBA{R: r, G: g, B: b, A: a}
}

func newCanvas(w, h int) *ebiten.Image {
	return ebiten.NewImage(w, h)
}

func initFaces() {
	src, err := text.NewGoTextFaceSource(bytes.NewReader(goregular.TTF))
	if err != nil {
		panic(err)
	}
	faceTiny = &text.GoTextFace{Source: src, Size: 11}
	faceSmall = &text.GoTextFace{Source: src, Size: 14}
	faceMed = &text.GoTextFace{Source: src, Size: 20}
	faceBig = &text.GoTextFace{Source: src, Size: 32}
	faceTitle = &text.GoTextFace{Source: src, Size: 52}
}

// ---- wall textures (64x64) ----

func makeBrickTexture() *ebiten.Image {
	t := newCanvas(64, 64)
	vector.DrawFilledRect(t, 0, 0, 64, 64, hex(0x6e3328), false)
	mortar := hex(0x2b1a16)
	for row := 0; row < 4; row++ {
		y := row * 16
		vector.DrawFilledRect(t, 0, float32(y), 64, 2, mortar, false)
		off := float32(0)
		if row%2 == 1 {
			off = 16
		}
		for bx := -32; bx <= 64; bx += 32 {
			vector.DrawFilledRect(t, float32(bx)+off+30, float32(y), 2, 16, mortar, false)
		}
	}
	return t
}

func makeStoneTexture() *ebiten.Image {
	t := newCanvas(64, 64)
	vector.DrawFilledRect(t, 0, 0, 64, 64, hex(0x4a5560), false)
	edge := hex(0x333c46)
	light := hex(0x5d6b78)
	for by := 0; by < 2; by++ {
		for bx := 0; bx < 2; bx++ {
			x, y := bx*32, by*32
			vector.DrawFilledRect(t, float32(x), float32(y), 30, 30, hex(0x525d68), false)
			vector.DrawFilledRect(t, float32(x), float32(y), 30, 2, light, false)
			vector.DrawFilledRect(t, float32(x), float32(y), 2, 30, light, false)
			vector.DrawFilledRect(t, float32(x), float32(y+28), 30, 2, edge, false)
			vector.DrawFilledRect(t, float32(x+28), float32(y), 2, 30, edge, false)
		}
	}
	return t
}

func makeTechTexture() *ebiten.Image {
	t := newCanvas(64, 64)
	vector.DrawFilledRect(t, 0, 0, 64, 64, hex(0x1c2a30), false)
	hex2 := hex(0x24363d)
	for y := 0; y < 64; y += 16 {
		vector.DrawFilledRect(t, 0, float32(y), 64, 1, hex2, false)
	}
	vector.DrawFilledRect(t, 0, 30, 64, 4, hex(0x0f5e66), false)
	vector.DrawFilledRect(t, 0, 31, 64, 2, hex(0x39d7e0), false)
	for x := 8; x < 64; x += 24 {
		vector.DrawFilledRect(t, float32(x), 8, 6, 6, hex(0x0f5e66), false)
	}
	return t
}

// ---- sprite assets ----

func drawImp(t *ebiten.Image, leg bool, attack bool) {
	vector.DrawFilledRect(t, 0, 0, 64, 64, color.RGBA{0, 0, 0, 0}, false)
	skin := hex(0x7a2f52)
	dark := hex(0x4e1c38)
	// legs
	l1, l2 := float32(0), float32(0)
	if leg {
		l1, l2 = -3, 3
	}
	vector.DrawFilledRect(t, 24+l1, 46, 6, 16, dark, false)
	vector.DrawFilledRect(t, 34+l2, 46, 6, 16, dark, false)
	// body
	vector.DrawFilledCircle(t, 32, 40, 14, skin, false)
	vector.DrawFilledCircle(t, 32, 42, 9, dark, false)
	// arms
	if attack {
		vector.DrawFilledRect(t, 12, 30, 10, 6, skin, false)
		vector.DrawFilledRect(t, 42, 30, 10, 6, skin, false)
		vector.DrawFilledCircle(t, 10, 33, 3, hex(0xffcf5c), false)
		vector.DrawFilledCircle(t, 54, 33, 3, hex(0xffcf5c), false)
	} else {
		vector.DrawFilledRect(t, 16, 34, 6, 16, skin, false)
		vector.DrawFilledRect(t, 42, 34, 6, 16, skin, false)
	}
	// head
	vector.DrawFilledCircle(t, 32, 20, 9, skin, false)
	vector.DrawFilledRect(t, 22, 8, 4, 8, skin, false)
	vector.DrawFilledRect(t, 38, 8, 4, 8, skin, false)
	vector.DrawFilledCircle(t, 28, 19, 2, hex(0xffe066), false)
	vector.DrawFilledCircle(t, 36, 19, 2, hex(0xffe066), false)
	vector.DrawFilledRect(t, 28, 25, 8, 2, dark, false)
}

func makeImpFrames() []*ebiten.Image {
	a := newCanvas(64, 64)
	drawImp(a, false, false)
	b := newCanvas(64, 64)
	drawImp(b, true, false)
	atk := newCanvas(64, 64)
	drawImp(atk, false, true)
	return []*ebiten.Image{a, b, atk}
}

func makeImpDead() *ebiten.Image {
	t := newCanvas(64, 64)
	vector.DrawFilledCircle(t, 32, 56, 13, hex(0x541d33), false)
	vector.DrawFilledCircle(t, 26, 54, 6, hex(0x6d2742), false)
	return t
}

func makeHeartSprite() *ebiten.Image {
	t := newCanvas(32, 32)
	vector.DrawFilledRect(t, 4, 8, 24, 18, hex(0xe8e8ea), false)
	vector.DrawFilledRect(t, 6, 10, 20, 14, hex(0xf4f4f6), false)
	vector.DrawFilledRect(t, 14, 12, 4, 10, hex(0xd63a4e), false)
	vector.DrawFilledRect(t, 11, 15, 10, 4, hex(0xd63a4e), false)
	return t
}

func makeAmmoSprite() *ebiten.Image {
	t := newCanvas(32, 32)
	vector.DrawFilledRect(t, 4, 10, 24, 16, hex(0x5a636e), false)
	vector.DrawFilledRect(t, 4, 10, 24, 3, hex(0x6d7885), false)
	for i := 0; i < 3; i++ {
		vector.DrawFilledRect(t, float32(8+i*6), 14, 3, 7, hex(0xffc93a), false)
	}
	return t
}

func makeBeamSprite() *ebiten.Image {
	t := newCanvas(48, 96)
	for i := 0; i < 96; i += 3 {
		a := uint8(160 - i)
		if i > 64 {
			a = uint8(96 - (i-64)*3)
		}
		vector.DrawFilledRect(t, 6, float32(i), 36, 3, rgba(60, 230, 240, a), false)
	}
	return t
}

func makeBackground() *ebiten.Image {
	t := newCanvas(InternalW, InternalH)
	vector.DrawFilledRect(t, 0, 0, InternalW, InternalH, hex(0x070a12), false)
	for i := 0; i < InternalH; i += 4 {
		a := uint8(18 + i/8)
		vector.DrawFilledRect(t, 0, float32(i), InternalW, 2, rgba(30, 60, 90, a), false)
	}
	for x := -InternalW; x < InternalW*2; x += 24 {
		vector.StrokeLine(t, float32(InternalW/2+x/3), 120, float32(x), InternalH, 1, rgba(40, 120, 160, 60), false)
	}
	vector.StrokeLine(t, 0, float32(InternalH-40), InternalW, float32(InternalH-40), 1, rgba(60, 200, 230, 90), false)
	return t
}

func initAssets() {
	wallTex = map[byte]*ebiten.Image{
		TileBrick: makeBrickTexture(),
		TileStone: makeStoneTexture(),
		TileTech:  makeTechTexture(),
	}
	impFrames = makeImpFrames()
	impDead = makeImpDead()
	heartSprite = makeHeartSprite()
	ammoSprite = makeAmmoSprite()
	beamSprite = makeBeamSprite()
	background = makeBackground()
}
