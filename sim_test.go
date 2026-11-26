package main

import (
	"math"
	"testing"
)

// TestSimulationNoPanic runs seconds of chaotic gameplay per level: nothing
// may panic and the world must stay consistent.
func TestSimulationNoPanic(t *testing.T) {
	for level := 1; level <= 6; level++ {
		g := &Game{}
		g.startLevel(level)
		g.state = statePlay
		w := g.world
		for f := 0; f < 60*20; f++ {
			p := &w.Player
			p.Dir += 0.02
			p.X, p.Y = tryMove(w.Grid, w.W, w.H, p.X, p.Y, math.Cos(p.Dir)*0.05, math.Sin(p.Dir)*0.05, playerRadius)
			if f%20 == 0 {
				w.T += 0.4
				w.fire(g)
			}
			w.update(1.0/60.0, g)
			if g.state == stateDead {
				break
			}
		}
	}
}

func TestExitCompletesLevel(t *testing.T) {
	for level := 1; level <= 6; level++ {
		g := &Game{}
		g.startLevel(level)
		g.state = statePlay
		w := g.world
		for _, e := range w.Enemies {
			if e.Kind == TileBoss {
				e.Dead = true
			}
		}
		ex := findKind(w, TileExit)
		w.Player.X, w.Player.Y = ex.X, ex.Y
		w.update(1.0/60.0, g)
		if !w.Cleared {
			t.Fatalf("level %d: stepping on exit should complete the level", level)
		}
	}
}

func TestEnemyContactKillsWoundedPlayer(t *testing.T) {
	g := &Game{}
	g.startLevel(1)
	g.state = statePlay
	w := g.world
	p := &w.Player
	p.HP = 10
	w.Enemies = append(w.Enemies, &Enemy{X: p.X + 0.3, Y: p.Y, HP: 100})
	for f := 0; f < 120 && g.state == statePlay; f++ {
		w.update(1.0/60.0, g)
	}
	if g.state != stateDead {
		t.Fatal("adjacent enemy should kill a wounded player")
	}
}

func TestSurvivalSimNoPanic(t *testing.T) {
	g := &Game{}
	g.mode = modeSurvive
	g.startLevel(1)
	g.state = statePlay
	maxWave := 1
	for f := 0; f < 60*40; f++ {
		w := g.world
		p := &w.Player
		p.Dir += 0.03
		p.X, p.Y = tryMove(w.Grid, w.W, w.H, p.X, p.Y, math.Cos(p.Dir)*0.05, math.Sin(p.Dir)*0.05, playerRadius)
		if f%30 == 0 {
			w.T += 0.5
			wi := (f / 30) % 5
			w.Player.Weapon = wi
			p.Ammo[wi] = []int{12, 6, 30, 5, 4}[wi]
			w.fire(g)
		}
		if f%60 == 0 {
			for _, e := range w.Enemies {
				if !e.Dead {
					w.damageEnemy(e, 1000, g)
					break
				}
			}
		}
		w.update(1.0/60.0, g)
		if w.Wave > maxWave {
			maxWave = w.Wave
		}
	}
	if maxWave <= 1 {
		t.Fatal("chaotic survival sim should clear at least wave 1")
	}
}
