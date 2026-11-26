package main

// Wall tile characters and helpers shared by parsing, raycasting and tests.

const (
	TileBrick    = '#'
	TileStone    = '%'
	TileTech     = '='
	TileSpawn    = 'P'
	TileEnemy    = 'e' // imp de base
	TileSpectre  = 's' // rapide et fragile
	TileGolem    = 'g' // lent et blindé
	TileSpitter  = 'c' // attaque à distance
	TileBoss     = 'B' // Gardien, boss du niveau 5
	TileColossus = 'C' // Colosse, second boss (charge + grosses salves)
	TileHealth   = 'h'
	TileAmmo     = 'a'
	TileExit     = 'X'
	TileFloor    = '.'
)

func isSolid(c byte) bool {
	return c == TileBrick || c == TileStone || c == TileTech
}

func isEnemy(c byte) bool {
	switch c {
	case TileEnemy, TileSpectre, TileGolem, TileSpitter, TileBoss, TileColossus:
		return true
	}
	return false
}
