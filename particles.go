package main

import "math"

// Effets de particules : sang, impacts muraux, eclats d'acide, ions,
// scintillements de ramassage. Tout est billboardise comme les sprites.

const maxParts = 360

var (
	bloodTint = map[byte][3]float32{
		TileEnemy:   {0.55, 0.10, 0.20},
		TileSpectre: {0.75, 0.30, 0.10},
		TileGolem:   {0.28, 0.45, 0.28},
		TileSpitter: {0.45, 0.75, 0.20},
		TileBoss:    {0.60, 0.08, 0.10},
	}
	defaultBlood = [3]float32{0.5, 0.12, 0.15}
	sparkCol     = [3]float32{1.0, 0.85, 0.40}
	dustCol      = [3]float32{0.45, 0.42, 0.36}
	acidCol      = [3]float32{0.45, 0.85, 0.15}
	ionCol       = [3]float32{0.35, 0.95, 0.60}
	bossCol      = [3]float32{0.85, 0.15, 0.15}
	healCol      = [3]float32{0.35, 0.95, 0.50}
	ammoCol      = [3]float32{0.95, 0.80, 0.20}
	plasmaCol    = [3]float32{0.35, 0.85, 0.95}
)

type Particle struct {
	X, Y, VX, VY float64
	Z, VZ        float64
	T, Life      float64
	Size         float64
	Col          [3]float32
	Glow         bool
}

func (w *World) spawnBurst(x, y, z, spd float64, n int, col [3]float32, glow bool, life, size float64) {
	for i := 0; i < n; i++ {
		if len(w.Parts) >= maxParts {
			return
		}
		a := w.rng.Float64() * 2 * math.Pi
		v := spd * (0.2 + 0.8*w.rng.Float64())
		w.Parts = append(w.Parts, &Particle{
			X: x, Y: y, Z: z,
			VX: math.Cos(a) * v, VY: math.Sin(a) * v,
			VZ:   (w.rng.Float64() - 0.25) * spd * 0.8,
			Life: life * (0.6 + 0.7*w.rng.Float64()),
			Size: size * (0.6 + 0.8*w.rng.Float64()),
			Col:  col, Glow: glow,
		})
	}
}

func (w *World) spawnBlood(x, y float64, kind byte, n int) {
	col, ok := bloodTint[kind]
	if !ok {
		col = defaultBlood
	}
	w.spawnBurst(x, y, 0.55, 2.2, n, col, false, 0.6, 0.050)
}

func (w *World) spawnSparks(x, y float64, n int) {
	w.spawnBurst(x, y, 0.6, 1.8, n, sparkCol, true, 0.35, 0.035)
	w.spawnBurst(x, y, 0.6, 1.1, n/2+1, dustCol, false, 0.7, 0.060)
}

func (w *World) spawnPickupGlow(x, y float64, col [3]float32) {
	w.spawnBurst(x, y, 0.4, 1.2, 8, col, true, 0.5, 0.045)
}

func (w *World) updateParticles(dt float64) {
	kept := w.Parts[:0]
	for _, p := range w.Parts {
		p.T += dt
		if p.T >= p.Life {
			continue
		}
		p.X += p.VX * dt
		p.Y += p.VY * dt
		damp := 1 - 2.0*dt
		p.VX *= damp
		p.VY *= damp
		p.Z += p.VZ * dt
		p.VZ -= 2.6 * dt
		if p.Z < 0.02 {
			p.Z = 0.02
			p.VZ *= -0.35
			p.VX *= 0.6
			p.VY *= 0.6
		}
		if w.solidAt(int(p.X), int(p.Y)) {
			continue
		}
		kept = append(kept, p)
	}
	w.Parts = kept
}
