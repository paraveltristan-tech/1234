package main

import "math"

// castStep performs one DDA raycast against the grid. It returns the
// perpendicular distance to the first solid cell, which side was hit
// (0 = vertical wall, 1 = horizontal wall), the tile character and the
// fractional position of the impact along the wall.
func castStep(grid [][]byte, w, h int, px, py, rdx, rdy float64) (dist float64, side int, tile byte, wallX float64) {
	mapX, mapY := int(px), int(py)
	deltaX := math.Abs(1 / rdx)
	deltaY := math.Abs(1 / rdy)
	var stepX, stepY int
	var sideX, sideY float64

	if rdx < 0 {
		stepX = -1
		sideX = (px - float64(mapX)) * deltaX
	} else {
		stepX = 1
		sideX = (float64(mapX) + 1 - px) * deltaX
	}
	if rdy < 0 {
		stepY = -1
		sideY = (py - float64(mapY)) * deltaY
	} else {
		stepY = 1
		sideY = (float64(mapY) + 1 - py) * deltaY
	}

	for steps := 0; steps < 512; steps++ {
		if sideX < sideY {
			sideX += deltaX
			mapX += stepX
			side = 0
		} else {
			sideY += deltaY
			mapY += stepY
			side = 1
		}
		if mapX < 0 || mapX >= w || mapY < 0 || mapY >= h {
			return 1e9, side, TileBrick, 0
		}
		tile = grid[mapY][mapX]
		if isSolid(tile) || (tile == TileSecret && !secretSeen[[2]int{mapX, mapY}]) {
			if side == 0 {
				dist = sideX - deltaX
				wallX = py + dist*rdy
			} else {
				dist = sideY - deltaY
				wallX = px + dist*rdx
			}
			wallX -= math.Floor(wallX)
			if dist < 1e-4 {
				dist = 1e-4
			}
			return dist, side, tile, wallX
		}
	}
	return 1e9, side, TileBrick, 0
}

// losClear reports whether a straight line between two points crosses no
// solid cell.
func losClear(grid [][]byte, w, h int, x0, y0, x1, y1 float64) bool {
	dx, dy := x1-x0, y1-y0
	d := math.Hypot(dx, dy)
	if d < 1e-6 {
		return true
	}
	steps := int(d / 0.08)
	for i := 1; i < steps; i++ {
		t := float64(i) / float64(steps)
		gx, gy := int(x0+dx*t), int(y0+dy*t)
		if gx < 0 || gx >= w || gy < 0 || gy >= h || isSolid(grid[gy][gx]) {
			return false
		}
	}
	return true
}

// tryMove applies axis-separated movement with collision against solid
// cells, so entities slide along walls instead of sticking.
func tryMove(grid [][]byte, w, h int, x, y, dx, dy, r float64) (float64, float64) {
	if !blocked(grid, w, h, x+dx, y, r) {
		x += dx
	}
	if !blocked(grid, w, h, x, y+dy, r) {
		y += dy
	}
	return x, y
}

func blocked(grid [][]byte, w, h int, x, y, r float64) bool {
	for _, cx := range [2]float64{x - r, x + r} {
		for _, cy := range [2]float64{y - r, y + r} {
			gx, gy := int(cx), int(cy)
			if gx < 0 || gx >= w || gy < 0 || gy >= h || isSolid(grid[gy][gx]) {
				return true
			}
		}
	}
	return false
}

func angDiff(a, b float64) float64 {
	d := a - b
	for d > math.Pi {
		d -= 2 * math.Pi
	}
	for d < -math.Pi {
		d += 2 * math.Pi
	}
	return d
}
