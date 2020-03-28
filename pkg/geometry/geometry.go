// Copyright 2020 John Slee <jslee@jslee.io>
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

package geometry

import "math"

// Point holds a Cartesian point and also an angle in degrees.
type RadialPoint struct {
	Angle float64
	X, Y  float64
}

// RadialPointGenerator generates a series of Cartesian points around a segment
// of a circle.
type RadialPointGenerator struct {
	// StartAngle and EndAngle indicate the span in degrees around the circle for
	// which points are generated. Zero degrees is at 9-o'clock and increases in
	// the clockwise direction.
	StartAngle, EndAngle float64
	// X and Y indicate the centre of the circle (origin)
	X, Y float64
	// Count indicates how many ticks are generated
	Count int
}

// GenerateAtRadius returns a set of evenly-distributed Cartesian points around
// a circle at a supplied radius.
func (rpg RadialPointGenerator) GenerateAtRadius(r float64) []RadialPoint {
	var points []RadialPoint
	interval := (rpg.EndAngle - rpg.StartAngle) / float64(rpg.Count-1)
	for i := 0; i < rpg.Count; i++ {
		angle := rpg.StartAngle + interval*float64(i)
		radians := angle * math.Pi / 180.0
		point := RadialPoint{
			Angle: angle,
			X:     rpg.X - r*math.Cos(radians),
			Y:     rpg.Y + r*math.Sin(radians),
		}
		points = append(points, point)
	}
	return points
}
