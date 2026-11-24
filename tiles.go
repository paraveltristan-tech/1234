package main

// Wall tile characters and helpers shared by parsing, raycasting and tests.

const (
	TileBrick  = '#'
	TileStone  = '%'
	TileTech   = '='
	TileSpawn  = 'P'
	TileEnemy  = 'e'
	TileHealth = 'h'
	TileAmmo   = 'a'
	TileExit   = 'X'
	TileFloor  = '.'
)

func isSolid(c byte) bool {
	return c == TileBrick || c == TileStone || c == TileTech
}
