package entities

// toIsometric converts 2D game coordinates to isometric projection
// This creates a 30-60-90 degree isometric view
func toIsometric(x, y float64) (float64, float64) {
	isoX := (x - y) * 0.866 // cos(30°) ≈ 0.866
	isoY := (x + y) * 0.5   // sin(30°) / 2
	return isoX, isoY
}

// fromIsometric converts isometric coordinates back to world coordinates
func fromIsometric(isoX, isoY float64) (float64, float64) {
	x := (isoX/0.866 + isoY*2) / 2
	y := (-isoX/0.866 + isoY*2) / 2
	return x, y
}
