package systems

import (
	"image"
	"image/color"
	"math"
)

// ============================================================================
// DRAWING HELPER FUNCTIONS
// ============================================================================

func drawFilledCircle(img *image.RGBA, cx, cy, radius int, c color.RGBA) {
	for y := cy - radius; y <= cy+radius; y++ {
		for x := cx - radius; x <= cx+radius; x++ {
			dx := x - cx
			dy := y - cy
			if dx*dx+dy*dy <= radius*radius {
				if x >= 0 && x < img.Bounds().Dx() && y >= 0 && y < img.Bounds().Dy() {
					img.Set(x, y, c)
				}
			}
		}
	}
}

func drawFilledEllipse(img *image.RGBA, cx, cy, rx, ry int, c color.RGBA) {
	for y := cy - ry; y <= cy+ry; y++ {
		for x := cx - rx; x <= cx+rx; x++ {
			dx := float64(x - cx)
			dy := float64(y - cy)
			if (dx*dx)/(float64(rx*rx))+(dy*dy)/(float64(ry*ry)) <= 1.0 {
				if x >= 0 && x < img.Bounds().Dx() && y >= 0 && y < img.Bounds().Dy() {
					img.Set(x, y, c)
				}
			}
		}
	}
}

func drawCircleOutline(img *image.RGBA, cx, cy, radius, thickness int, c color.RGBA) {
	for y := cy - radius - thickness; y <= cy+radius+thickness; y++ {
		for x := cx - radius - thickness; x <= cx+radius+thickness; x++ {
			dx := x - cx
			dy := y - cy
			dist := dx*dx + dy*dy
			innerRadius := (radius - thickness) * (radius - thickness)
			outerRadius := (radius + thickness) * (radius + thickness)
			if dist >= innerRadius && dist <= outerRadius {
				if x >= 0 && x < img.Bounds().Dx() && y >= 0 && y < img.Bounds().Dy() {
					img.Set(x, y, c)
				}
			}
		}
	}
}

func drawDottedCircle(img *image.RGBA, cx, cy, radius, dotSize int, c color.RGBA) {
	numDots := 24
	for i := 0; i < numDots; i++ {
		angle := float64(i) * math.Pi * 2 / float64(numDots)
		x := cx + int(math.Cos(angle)*float64(radius))
		y := cy + int(math.Sin(angle)*float64(radius))
		drawFilledCircle(img, x, y, dotSize, c)
	}
}

func drawFilledRect(img *image.RGBA, x, y, w, h int, c color.RGBA) {
	for py := y; py < y+h; py++ {
		for px := x; px < x+w; px++ {
			if px >= 0 && px < img.Bounds().Dx() && py >= 0 && py < img.Bounds().Dy() {
				img.Set(px, py, c)
			}
		}
	}
}

func drawRectOutline(img *image.RGBA, x, y, w, h, thickness int, c color.RGBA) {
	// Top
	drawFilledRect(img, x, y, w, thickness, c)
	// Bottom
	drawFilledRect(img, x, y+h-thickness, w, thickness, c)
	// Left
	drawFilledRect(img, x, y, thickness, h, c)
	// Right
	drawFilledRect(img, x+w-thickness, y, thickness, h, c)
}

func drawLine(img *image.RGBA, x1, y1, x2, y2, thickness int, c color.RGBA) {
	dx := x2 - x1
	dy := y2 - y1
	steps := int(math.Sqrt(float64(dx*dx + dy*dy)))

	for i := 0; i <= steps; i++ {
		t := float64(i) / float64(steps)
		x := x1 + int(float64(dx)*t)
		y := y1 + int(float64(dy)*t)

		for tx := -thickness / 2; tx <= thickness/2; tx++ {
			for ty := -thickness / 2; ty <= thickness/2; ty++ {
				px := x + tx
				py := y + ty
				if px >= 0 && px < img.Bounds().Dx() && py >= 0 && py < img.Bounds().Dy() {
					img.Set(px, py, c)
				}
			}
		}
	}
}

func drawFilledTriangle(img *image.RGBA, x1, y1, x2, y2, x3, y3 int, c color.RGBA) {
	// Find bounding box
	minX := min(x1, min(x2, x3))
	maxX := max(x1, max(x2, x3))
	minY := min(y1, min(y2, y3))
	maxY := max(y1, max(y2, y3))

	// Check each pixel in bounding box
	for y := minY; y <= maxY; y++ {
		for x := minX; x <= maxX; x++ {
			if pointInTriangle(x, y, x1, y1, x2, y2, x3, y3) {
				if x >= 0 && x < img.Bounds().Dx() && y >= 0 && y < img.Bounds().Dy() {
					img.Set(x, y, c)
				}
			}
		}
	}
}

func pointInTriangle(px, py, x1, y1, x2, y2, x3, y3 int) bool {
	// Barycentric coordinate method
	denominator := float64((y2-y3)*(x1-x3) + (x3-x2)*(y1-y3))
	if denominator == 0 {
		return false
	}

	a := float64((y2-y3)*(px-x3)+(x3-x2)*(py-y3)) / denominator
	b := float64((y3-y1)*(px-x3)+(x1-x3)*(py-y3)) / denominator
	c := 1 - a - b

	return a >= 0 && a <= 1 && b >= 0 && b <= 1 && c >= 0 && c <= 1
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
