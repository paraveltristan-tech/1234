package main

import (
	"bytes"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
)

// Moteur son 100% procedural : tous les bruitages et les deux musiques sont
// synthetises en memoire (PCM float32 stereo), sans aucun fichier externe.

const audioRate = 44100

var snd *sndSys

type sndSys struct {
	ctx       *audio.Context
	ones      map[string][]byte
	music     *audio.Player
	musicCur  string
	engine    *audio.Player
	engineCur byte
	live      []*audio.Player
	muted     bool
}

// initAudio pre-rend tous les sons ; appele depuis main, jamais en tests.
func initAudio() {
	s := &sndSys{ctx: audio.NewContext(audioRate)}
	s.ones = buildSounds()
	snd = s
}

func playSnd(name string) {
	if snd == nil || snd.muted {
		return
	}
	buf, ok := snd.ones[name]
	if !ok {
		return
	}
	p, err := snd.ctx.NewPlayer(bytes.NewReader(buf))
	if err != nil {
		return
	}
	p.Play()
	snd.live = append(snd.live, p)
	keep := snd.live[:0]
	for _, lp := range snd.live {
		if lp.IsPlaying() {
			lp.Close()
		} else {
			keep = append(keep, lp)
		}
	}
	snd.live = keep
}

// tickAudio pilote la musique selon l'etat du jeu et le moteur du vehicule.
func (g *Game) tickAudio() {
	if snd == nil {
		return
	}
	snd.muted = g.muted
	name := "musique_combat"
	switch g.state {
	case stateTitle, stateWin, stateDead:
		name = "musique_calme"
	}
	snd.setMusic(name)
	kind := byte(0)
	moving := false
	if g.world != nil && g.world.VehMounted && g.state == statePlay {
		kind = g.world.VehKind
		moving = ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyUp) ||
			ebiten.IsKeyPressed(ebiten.KeyS) || ebiten.IsKeyPressed(ebiten.KeyDown) ||
			ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyD)
	}
	snd.setEngine(kind, moving)
}

func (s *sndSys) vol(base float64) float64 {
	if s.muted {
		return 0
	}
	return base
}

func (s *sndSys) setMusic(name string) {
	if s.musicCur == name && s.music != nil {
		s.music.SetVolume(s.vol(0.30))
		return
	}
	if s.music != nil {
		s.music.Close()
		s.music = nil
	}
	s.musicCur = name
	buf := s.ones[name]
	if buf == nil {
		return
	}
	p, err := s.ctx.NewPlayer(audio.NewInfiniteLoop(bytes.NewReader(buf), int64(len(buf))))
	if err != nil {
		return
	}
	p.SetVolume(s.vol(0.30))
	p.Play()
	s.music = p
}

func (s *sndSys) setEngine(kind byte, moving bool) {
	if kind != s.engineCur {
		if s.engine != nil {
			s.engine.Close()
			s.engine = nil
		}
		s.engineCur = kind
		if kind != 0 {
			names := map[byte]string{'V': "moteur_vtt", 'T': "moteur_tank", 'H': "moteur_heli"}
			if buf := s.ones[names[kind]]; buf != nil {
				p, err := s.ctx.NewPlayer(audio.NewInfiniteLoop(bytes.NewReader(buf), int64(len(buf))))
				if err == nil {
					p.Play()
					s.engine = p
				}
			}
		}
	}
	if s.engine != nil {
		v := 0.0
		if kind != 0 {
			if moving {
				v = 1.0
			} else {
				v = 0.42
			}
		}
		s.engine.SetVolume(s.vol(v))
	}
}

// ---- synthese ----

func frac(x float64) float64 { return x - math.Floor(x) }

func hnoise(t float64) float64 { return 2*frac(math.Sin(t*91713.7)*43758.5453) - 1 }

func render(dur float64, fn func(t float64) float64) []float64 {
	n := int(audioRate * dur)
	out := make([]float64, n)
	for i := range out {
		out[i] = fn(float64(i) / audioRate)
	}
	return out
}

func tone(dur, f0, f1, attack, decay, vol float64, wave func(float64) float64) []float64 {
	return render(dur, func(t float64) float64 {
		p := t / dur
		f := f0 + (f1-f0)*p
		env := math.Min(1, t/attack) * math.Exp(-decay*t)
		return wave(2*math.Pi*f*t) * env * vol
	})
}

func noiseBurst(dur, decay, vol float64) []float64 {
	return render(dur, func(t float64) float64 {
		return hnoise(t*97.31) * math.Exp(-t*decay) * vol
	})
}

func sine(ph float64) float64 { return math.Sin(ph) }

func square(ph float64) float64 {
	if math.Sin(ph) > 0 {
		return 0.5
	}
	return -0.5
}

func saw(ph float64) float64 { return 2*frac(ph/(2*math.Pi)) - 1 }

func mix(parts ...[]float64) []float64 {
	n := 0
	for _, p := range parts {
		if len(p) > n {
			n = len(p)
		}
	}
	out := make([]float64, n)
	for _, p := range parts {
		for i, v := range p {
			out[i] += v
		}
	}
	return out
}

func addAt(dst, src []float64, offset float64) {
	start := int(offset * audioRate)
	for i, v := range src {
		j := start + i
		if j >= 0 && j < len(dst) {
			dst[j] += v
		}
	}
}

func toLE(m []float64) []byte {
	b := make([]byte, 4*len(m))
	i := 0
	for _, v := range m {
		if v > 1 {
			v = 1
		}
		if v < -1 {
			v = -1
		}
		bits := int16(v * 32000)
		b[i] = byte(bits)
		b[i+1] = byte(bits >> 8)
		b[i+2] = byte(bits)
		b[i+3] = byte(bits >> 8)
		i += 4
	}
	return b
}

func engineLoop(hz, rough float64) []float64 {
	return render(1.0, func(t float64) float64 {
		ph := 2 * math.Pi * hz * t
		v := math.Sin(ph)*0.6 + math.Sin(2*ph)*0.25 + math.Sin(3*ph)*0.15
		v += hnoise(t*61.3) * 0.09 * rough
		wob := 0.75 + 0.25*math.Sin(2*math.Pi*3.1*t)
		return v * wob * 0.5
	})
}

func rotorLoop() []float64 {
	return render(1.0, func(t float64) float64 {
		pulse := math.Pow(0.5+0.5*math.Cos(2*math.Pi*10.5*t), 2)
		v := hnoise(t*53.7)*0.8 + math.Sin(2*math.Pi*120*t)*0.15
		return v * pulse * 0.55
	})
}

func calmMusic() []float64 {
	out := make([]float64, int(audioRate*8))
	roots := []float64{110, 87.31, 130.81, 98}
	mels := [][]float64{{440, 523.25}, {349.23, 523.25}, {523.25, 659.25}, {392, 587.33}}
	for ci, root := range roots {
		start := int(float64(ci) * 2 * audioRate)
		for i := 0; i < int(2*audioRate) && start+i < len(out); i++ {
			t := float64(i) / audioRate
			env := math.Min(1, t*2) * (0.65 + 0.35*math.Cos(2*math.Pi*t/4))
			v := math.Sin(2*math.Pi*root*t)*0.20 + math.Sin(2*math.Pi*root*1.5*t)*0.09
			mn := mels[ci][0]
			if t > 1 {
				mn = mels[ci][1]
			}
			v += math.Sin(2*math.Pi*mn*t) * 0.11
			out[start+i] += v * env
		}
	}
	return out
}

func combatMusic() []float64 {
	step := 0.2
	bass := []float64{55, 55, 65.41, 55, 82.41, 82.41, 98, 82.41,
		55, 55, 65.41, 55, 73.42, 73.42, 98, 73.42,
		55, 55, 65.41, 55, 82.41, 98, 110, 98}
	mel := []float64{220, 0, 261.63, 0, 329.63, 261.63, 0, 220,
		246.94, 0, 293.66, 0, 329.63, 0, 246.94, 220,
		261.63, 0, 329.63, 392, 0, 329.63, 261.63, 0}
	out := make([]float64, int(audioRate*float64(len(bass))*step))
	for s := 0; s < len(bass); s++ {
		start := int(float64(s) * step * audioRate)
		for i := 0; i < int(step*audioRate) && start+i < len(out); i++ {
			t := float64(i) / audioRate
			env := math.Exp(-t * 6)
			v := square(2*math.Pi*bass[s]*t) * 0.16 * env
			if mel[s] > 0 {
				v += square(2*math.Pi*mel[s]*t) * 0.09 * env
			}
			out[start+i] += v
		}
	}
	return out
}

func buildSounds() map[string][]byte {
	m := map[string][]byte{}
	m["pistolet"] = toLE(mix(tone(0.13, 95, 50, 0.002, 30, 0.55, sine), noiseBurst(0.13, 42, 0.5)))
	m["pompe"] = toLE(mix(tone(0.28, 70, 38, 0.002, 13, 0.85, sine), noiseBurst(0.28, 15, 0.75)))
	m["plasma"] = toLE(mix(tone(0.17, 980, 240, 0.004, 16, 0.5, sine), tone(0.17, 1470, 360, 0.004, 18, 0.18, square)))
	m["sniper"] = toLE(mix(tone(0.42, 58, 30, 0.001, 7, 0.9, sine), noiseBurst(0.42, 9, 0.85), noiseBurst(0.02, 60, 0.9)))
	m["ions"] = toLE(mix(tone(0.42, 180, 760, 0.01, 5, 0.55, sine), tone(0.42, 720, 90, 0.05, 6, 0.2, square), noiseBurst(0.30, 8, 0.25)))
	m["obus"] = toLE(mix(tone(0.36, 46, 24, 0.002, 8, 1.0, sine), noiseBurst(0.36, 18, 0.8)))
	m["explosion"] = toLE(mix(noiseBurst(0.85, 4.2, 1.0), tone(0.85, 42, 20, 0.004, 3.5, 0.9, sine)))
	m["roquette"] = toLE(mix(noiseBurst(0.30, 7, 0.7), tone(0.30, 160, 60, 0.002, 6, 0.45, sine)))
	m["clac"] = toLE(mix(noiseBurst(0.035, 110, 0.8), render(0.16, func(t float64) float64 {
		if t < 0.09 {
			return 0
		}
		return hnoise(t*87.7) * math.Exp(-(t-0.09)*110) * 0.8
	})))
	m["grelot"] = toLE(mix(tone(0.5, 660, 660, 0.005, 5, 0.4, sine), tone(0.45, 880, 880, 0.005, 5, 0.35, sine), tone(0.4, 1320, 1320, 0.005, 6, 0.3, sine)))
	m["blip"] = toLE(tone(0.09, 520, 520, 0.003, 26, 0.35, sine))
	m["rugissement"] = toLE(mix(tone(0.95, 88, 42, 0.06, 3.2, 0.8, saw), noiseBurst(0.95, 3.5, 0.3), tone(0.95, 44, 22, 0.02, 3, 0.5, sine)))
	m["mort"] = toLE(tone(1.0, 290, 60, 0.01, 2.6, 0.6, saw))
	win := make([]float64, int(audioRate*1.3))
	addAt(win, tone(0.3, 523.25, 523.25, 0.005, 4, 0.4, sine), 0)
	addAt(win, tone(0.3, 659.25, 659.25, 0.005, 4, 0.4, sine), 0.12)
	addAt(win, tone(0.3, 783.99, 783.99, 0.005, 4, 0.4, sine), 0.24)
	addAt(win, tone(0.55, 1046.5, 1046.5, 0.005, 3, 0.45, sine), 0.36)
	addAt(win, tone(0.55, 261.63, 261.63, 0.005, 3, 0.3, sine), 0.36)
	m["victoire"] = toLE(win)
	m["moteur_vtt"] = toLE(engineLoop(76, 0.9))
	m["moteur_tank"] = toLE(engineLoop(47, 1.3))
	m["moteur_heli"] = toLE(rotorLoop())
	m["musique_calme"] = toLE(calmMusic())
	m["musique_combat"] = toLE(combatMusic())
	return m
}
