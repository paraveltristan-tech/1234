package main

import (
	"bytes"
	"image"
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

// deterministic hash noise in [0,1)
func nz(x, y int, salt uint32) float64 {
	h := uint32(x)*374761393 + uint32(y)*668265263 + salt*1274126177
	h ^= h >> 13
	h *= 1103515245
	h ^= h >> 16
	return float64(h&0xffff) / 65536.0
}

func clamp8(v int) uint8 {
	if v < 0 {
		return 0
	}
	if v > 255 {
		return 255
	}
	return uint8(v)
}

// ---- wall textures (64x64, per-pixel) ----

func makeBrickRGBA() *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, 64, 64))
	for y := 0; y < 64; y++ {
		row := y / 16
		off := 0
		if row%2 == 1 {
			off = 16
		}
		for x := 0; x < 64; x++ {
			mortar := y%16 < 2 || (x+off)%32 < 2
			var r, g, b int
			if mortar {
				r, g, b = 38, 26, 24
			} else {
				brickID := (x+off)/32 + row*17
				tnz := nz(x, y, uint32(100+brickID))
				r = 122 + (brickID%5)*8 - 14
				g = 56 + (brickID%3)*7 - 10
				b = 44 + (brickID%4)*5 - 12
				if y%16 == 2 || y%16 == 3 {
					r += 26
					g += 18
					b += 12
				}
				if y%16 == 14 || y%16 == 15 {
					r -= 20
					g -= 14
					b -= 10
				}
				n := int((tnz - 0.5) * 22)
				r, g, b = r+n, g+n, b+n
			}
			img.Set(x, y, color.RGBA{R: clamp8(r), G: clamp8(g), B: clamp8(b), A: 255})
		}
	}
	return img
}

func makeStoneRGBA() *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, 64, 64))
	for y := 0; y < 64; y++ {
		for x := 0; x < 64; x++ {
			bx, by := x%32, y%32
			var r, g, b int
			if bx < 2 || by < 2 || bx > 29 || by > 29 {
				r, g, b = 40, 46, 54
			} else {
				tnz := nz(x, y, 777)
				r, g, b = 84, 94, 104
				if bx < 5 || by < 5 {
					r += 16
					g += 18
					b += 16
				}
				if bx > 26 || by > 26 {
					r -= 18
					g -= 20
					b -= 18
				}
				n := int((tnz - 0.5) * 18)
				r, g, b = r+n, g+n, b+n
				if nz(x, y, 9012) > 0.995 {
					r, g, b = 30, 34, 40
				}
				if nz(x/3, y, 4242) > 0.99 && bx > 8 && by > 8 {
					r -= 30
					g -= 30
					b -= 30
				}
			}
			img.Set(x, y, color.RGBA{R: clamp8(r), G: clamp8(g), B: clamp8(b), A: 255})
		}
	}
	return img
}

func makeTechRGBA() *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, 64, 64))
	for y := 0; y < 64; y++ {
		for x := 0; x < 64; x++ {
			var r, g, b int
			switch {
			case y >= 29 && y <= 34:
				r, g, b = 20, 130, 140
				if y >= 31 && y <= 32 {
					r, g, b = 70, 226, 235
				}
			case x%16 == 0 || y%16 == 0:
				r, g, b = 22, 32, 38
			case x%16 < 3 || y%16 < 3:
				r, g, b = 40, 54, 62
			default:
				tnz := nz(x, y, 313)
				n := int((tnz - 0.5) * 10)
				r, g, b = 34+n, 48+n, 56+n
			}
			if (x%16 == 12 || x%16 == 13) && (y%16 == 12 || y%16 == 13) {
				if x%16 == 12 && y%16 == 12 {
					r, g, b = 150, 190, 200
				} else {
					r, g, b = 110, 150, 160
				}
			}
			img.Set(x, y, color.RGBA{R: clamp8(r), G: clamp8(g), B: clamp8(b), A: 255})
		}
	}
	return img
}

// makeFogVariants pre-renders 8 shaded copies of a texture blending toward
// the level fog color, so distance fog needs no per-pixel work at draw time.
func makeFogVariants(src *image.RGBA, fog [3]uint8) []*ebiten.Image {
	var out []*ebiten.Image
	for k := 0; k < 8; k++ {
		f := float64(k) / 7.0
		img := image.NewRGBA(image.Rect(0, 0, 64, 64))
		for y := 0; y < 64; y++ {
			for x := 0; x < 64; x++ {
				r0, g0, b0, _ := src.At(x, y).RGBA()
				r := int(float64(r0>>8)*f) + int(float64(fog[0])*(1-f))
				g := int(float64(g0>>8)*f) + int(float64(fog[1])*(1-f))
				b := int(float64(b0>>8)*f) + int(float64(fog[2])*(1-f))
				img.Set(x, y, color.RGBA{R: clamp8(r), G: clamp8(g), B: clamp8(b), A: 255})
			}
		}
		out = append(out, ebiten.NewImageFromImage(img))
	}
	return out
}

// ---- sprite assets ----

func drawImp(t *ebiten.Image, leg bool, attack bool) {
	skin := hex(0x7a2f52)
	dark := hex(0x4e1c38)
	l1, l2 := float32(0), float32(0)
	if leg {
		l1, l2 = -3, 3
	}
	vector.DrawFilledRect(t, 24+l1, 46, 6, 16, dark, false)
	vector.DrawFilledRect(t, 34+l2, 46, 6, 16, dark, false)
	vector.DrawFilledCircle(t, 32, 40, 14, skin, false)
	vector.DrawFilledCircle(t, 32, 42, 9, dark, false)
	if attack {
		vector.DrawFilledRect(t, 12, 30, 10, 6, skin, false)
		vector.DrawFilledRect(t, 42, 30, 10, 6, skin, false)
		vector.DrawFilledCircle(t, 10, 33, 3, hex(0xffcf5c), false)
		vector.DrawFilledCircle(t, 54, 33, 3, hex(0xffcf5c), false)
	} else {
		vector.DrawFilledRect(t, 16, 34, 6, 16, skin, false)
		vector.DrawFilledRect(t, 42, 34, 6, 16, skin, false)
	}
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
	levelPalettes = loadPalettes()
	brickTex = makeBrickRGBA()
	stoneTex = makeStoneRGBA()
	techTex = makeTechRGBA()
	rebuildFogSets()
	impFrames = makeImpFrames()
	impDead = makeImpDead()
	heartSprite = makeHeartSprite()
	ammoSprite = makeAmmoSprite()
	beamSprite = makeBeamSprite()
	background = makeBackground()
}
