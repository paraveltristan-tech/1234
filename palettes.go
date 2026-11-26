package main

// Per-level atmosphere: fog color, ceiling and floor gradients.

type levelPalette struct {
	CeilTop    [3]uint8
	CeilBottom [3]uint8
	FloorNear  [3]uint8
	FloorFar   [3]uint8
	Fog        [3]uint8
}

var levelPalettes []levelPalette

func loadPalettes() []levelPalette {
	return []levelPalette{
		{
			CeilTop:    [3]uint8{10, 14, 26},
			CeilBottom: [3]uint8{46, 54, 78},
			FloorNear:  [3]uint8{78, 62, 52},
			FloorFar:   [3]uint8{30, 24, 22},
			Fog:        [3]uint8{26, 18, 22},
		},
		{
			CeilTop:    [3]uint8{6, 14, 20},
			CeilBottom: [3]uint8{28, 50, 62},
			FloorNear:  [3]uint8{52, 66, 72},
			FloorFar:   [3]uint8{16, 26, 30},
			Fog:        [3]uint8{10, 22, 26},
		},
		{
			CeilTop:    [3]uint8{18, 8, 14},
			CeilBottom: [3]uint8{58, 28, 40},
			FloorNear:  [3]uint8{74, 40, 36},
			FloorFar:   [3]uint8{28, 14, 16},
			Fog:        [3]uint8{28, 12, 16},
		},
		{
			CeilTop:    [3]uint8{16, 8, 28},
			CeilBottom: [3]uint8{64, 32, 88},
			FloorNear:  [3]uint8{80, 48, 78},
			FloorFar:   [3]uint8{30, 16, 34},
			Fog:        [3]uint8{26, 14, 36},
		},
		{
			CeilTop:    [3]uint8{22, 14, 6},
			CeilBottom: [3]uint8{92, 58, 26},
			FloorNear:  [3]uint8{98, 66, 32},
			FloorFar:   [3]uint8{34, 20, 12},
			Fog:        [3]uint8{38, 24, 12},
		},
		{
			CeilTop:    [3]uint8{36, 26, 14},
			CeilBottom: [3]uint8{120, 92, 52},
			FloorNear:  [3]uint8{158, 128, 82},
			FloorFar:   [3]uint8{62, 46, 28},
			Fog:        [3]uint8{96, 74, 44},
		},
	}
}

func paletteFor(level int) levelPalette {
	idx := level - 1
	if idx < 0 {
		idx = 0
	}
	if idx >= len(levelPalettes) {
		idx = len(levelPalettes) - 1
	}
	return levelPalettes[idx]
}
