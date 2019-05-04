package util

import (
	"github.com/jsleeio/go-eagle/pkg/eagle"
)

// WireRectangle generates a rectangle
func WireRectangle(x1, y1, x2, y2 float64, layer int, width float64) []eagle.Wire {
	return []eagle.Wire{
		{X1: x1, Y1: y1, X2: x2, Y2: y1, Layer: layer, Width: width}, // bottom
		{X1: x1, Y1: y2, X2: x2, Y2: y2, Layer: layer, Width: width}, // top
		{X1: x1, Y1: y1, X2: x1, Y2: y2, Layer: layer, Width: width}, // left
		{X1: x2, Y1: y1, X2: x2, Y2: y2, Layer: layer, Width: width}, // right
	}
}
