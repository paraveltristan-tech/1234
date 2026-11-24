package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type gameState int

const (
	stateTitle gameState = iota
	statePlay
	statePause
	stateDead
	stateWin
)

const (
	InternalW = 480
	InternalH = 300
	ScaleF    = 2.0
)

type Game struct {
	state       gameState
	level       int
	score       int
	kills       int
	world       *World
	clearT      float64
	time        float64
	ShowMinimap bool
	ShowBigMap  bool
}

var (
	faceTiny, faceSmall, faceMed, faceBig, faceTitle text.Face
	assetsInited                                     bool
)

func main() {
	initFaces()
	ebiten.SetWindowSize(int(InternalW*ScaleF), int(InternalH*ScaleF))
	ebiten.SetWindowTitle("PROTOCOLE NEON - FPS")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	g := &Game{state: stateTitle, level: 1}
	if err := ebiten.RunGame(g); err != nil {
		panic(err)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return int(InternalW * ScaleF), int(InternalH * ScaleF)
}
