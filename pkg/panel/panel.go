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

package panel

// Point defines a single metric coordinate in a 2D space.
type Point struct {
	X, Y float64
}

// Panel types encapsulate physical characteristics of a rail-mounted panel
// All coordinates, distances and sizes are indicated in millimetres
type Panel interface {
	// MountingHoles returns a list of Points indicating mounting hole locations
	MountingHoles() []Point

	// MountingHoleDiameter returns the appropriate mounting hole diameter for
	// the panel format
	MountingHoleDiameter() float64

	// Height returns the Y-dimension size for the panel, eg. 128.5mm for panels
	// in the Eurorack system. This does NOT include tolerance adjustments
	Height() float64

	// Height returns the X-dimension size for the panel. This does not include
	// tolerance adjustments
	Width() float64

	// HorizontalFit returns the panel tolerance amount in the horizontal axis.
	// When creating panel outlines, this tolerance amount should be added to
	// the left-edge X coordinate, and subtracted from the right-edge coordinate,
	// resulting in the panel being slightly narrower than the "correct" width.
	//
	// As this is intended to improve panel fit in a system with panels of varying
	// tolerances, this adjustment should only be applied to the left and right
	// edges of the outline, and NOT to the X coordinates of any other features
	// of the panel. (and especially not the mounting holes!)
	HorizontalFit() float64

	// RailHeightFromMountingHole indicates how far up (from centre of bottom
	// mounting hole) or down (from centre of top mounting hole) the mounting
	// rail extends. This can be used to define KeepOut areas on the panel
	//
	// eg. For all Eurorack-related formats this is likely to be around 5mm,
	// though the exact figure will differ with rail type.
	//
	// This is primarily used to determine how much empty space there is between
	// the mounting rails, so best to err on the side of larger than smaller
	RailHeightFromMountingHole() float64

	// MountingHoleTopY returns the Y coordinate for the top row of mounting
	// holes
	MountingHoleTopY() float64

	// MountingHoleBottomY returns the Y coordinate for the bottom row of
	// mounting holes
	MountingHoleBottomY() float64
}

func LeftX(spec Panel) float64 {
	return spec.HorizontalFit() / 2
}

func RightX(spec Panel) float64 {
	return spec.Width() - spec.HorizontalFit()/2
}

func TopY(spec Panel) float64 {
	return spec.Height()
}

func BottomY(spec Panel) float64 {
	return 0
}

func TopLeft(spec Panel) Point {
	return Point{X: LeftX(spec), Y: TopY(spec)}
}

func TopRight(spec Panel) Point {
	return Point{X: RightX(spec), Y: TopY(spec)}
}

func BottomLeft(spec Panel) Point {
	return Point{X: LeftX(spec), Y: BottomY(spec)}
}

func BottomRight(spec Panel) Point {
	return Point{X: RightX(spec), Y: BottomY(spec)}
}
