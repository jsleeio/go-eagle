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

package outline

import (
	"math"

	"github.com/jsleeio/go-eagle/pkg/eagle"
)

// FindBoardOutlineWires searches through Wires in the Plain section of
// the board for zero-width wires in the Dimension layer
func FindBoardOutlineWires(e *eagle.Eagle) []eagle.Wire {
	wires := []eagle.Wire{}
	dimension := e.LayerByName("Dimension")
	for _, wire := range e.Board.Plain.Wires {
		if wire.Layer == dimension && wire.Width == 0.0 {
			wires = append(wires, wire)
		}
	}
	return wires
}

// BoardCoords holds information about a board outline and its place in
// the coordinate space. This is used to determine panel width and to
// correctly align the board with the panel
type BoardCoords struct {
	XMin, XMax, YMin, YMax float64
	XOffset, YOffset       float64
	HP                     int
}

// Width returns the width of the board outline
func (bc BoardCoords) Width() float64 {
	return bc.XMax - bc.XMin
}

// Height returns the height of the board outline
func (bc BoardCoords) Height() float64 {
	return bc.YMax - bc.YMin
}

// DeriveBoardCoords creates a BoardCoords object from the discovered outline
// wires in the Plain section of a board
func DeriveBoardCoords(e *eagle.Eagle) BoardCoords {
	bc := BoardCoords{}
	for _, wire := range FindBoardOutlineWires(e) {
		txmin, txmax := fsort2(wire.X1, wire.X2)
		tymin, tymax := fsort2(wire.Y1, wire.Y2)
		bc.XMin = math.Min(bc.XMin, txmin)
		bc.XMax = math.Max(bc.XMax, txmax)
		bc.YMin = math.Min(bc.YMin, tymin)
		bc.YMax = math.Max(bc.YMax, tymax)
	}
	if bc.XMin != 0 {
		bc.XOffset = -bc.XMin
	}
	if bc.YMin != 0 {
		bc.YOffset = -bc.YMin
	}
	bc.HP = int(math.Ceil(math.Ceil(bc.XMax-bc.XMin) / 5.08))
	return bc
}

func fsort2(a, b float64) (float64, float64) {
	if a > b {
		return b, a
	}
	return a, b
}
