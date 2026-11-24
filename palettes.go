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
