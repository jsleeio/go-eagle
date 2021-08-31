// Copyright 2019 John Slee <jslee@jslee.io>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to
// deal in the Software without restriction, including without limitation the
// rights to use, copy, modify, merge, publish, distribute, sublicense, and/or
// sell copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS
// IN THE SOFTWARE.

package util

import (
	"github.com/jsleeio/go-eagle/pkg/eagle"
)

// WireRectangle generates a rectangle
func WireRectangle(x1, y1, x2, y2 float64, layer int, width float64, radius float64) []eagle.Wire {
	segments := []eagle.Wire{
		{X1: x1 + radius, Y1: y1, X2: x2 - radius, Y2: y1, Layer: layer, Width: width}, // bottom
		{X1: x1 + radius, Y1: y2, X2: x2 - radius, Y2: y2, Layer: layer, Width: width}, // top
		{X1: x1, Y1: y1 + radius, X2: x1, Y2: y2 - radius, Layer: layer, Width: width}, // left
		{X1: x2, Y1: y1 + radius, X2: x2, Y2: y2 - radius, Layer: layer, Width: width}, // right
	}
	if radius > 0.0 {
		segments = append(segments,
			eagle.Wire{Curve: -90.0, Layer: layer, X1: x1 + radius, Y1: y1, X2: x1, Y2: y1 + radius, Width: width}, // bottom-left
			eagle.Wire{Curve: -90.0, Layer: layer, X1: x1, Y1: y2 - radius, X2: x1 + radius, Y2: y2, Width: width}, // top-left
			eagle.Wire{Curve: -90.0, Layer: layer, X1: x2 - radius, Y1: y2, X2: x2, Y2: y2 - radius, Width: width}, // top-right
			eagle.Wire{Curve: -90.0, Layer: layer, X1: x2, Y1: y1 + radius, X2: x2 - radius, Y2: y1, Width: width}) // bottom-right
	}
	return segments
}
