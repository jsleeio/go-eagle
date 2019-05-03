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

package intellijel

import (
	"github.com/jsleeio/go-eagle/pkg/format/eurorack"
	"github.com/jsleeio/go-eagle/pkg/panel"
)

// based on https://intellijel.com/support/1u-technical-specifications/

const (
	// PanelHeight1U represents the total height of an Intellijel 1U panel, in
	// millimetres
	PanelHeight1U = 39.65

	// MountingHolesLeftOffset represents the distance of the first mounting
	// hole from the left edge of the panel, in millimetres
	MountingHolesLeftOffset = eurorack.MountingHolesLeftOffset

	// MountingHoleTopY1U represents the Y value for the top row of 1U mounting
	// holes, in millimetres
	MountingHoleTopY1U = PanelHeight1U - 3.00

	// MountingHoleBottomY1U represents the Y value for the bottom row of 1U
	// mounting holes, in millimetres
	MountingHoleBottomY1U = 3.00

	// MountingHoleDiameter represents the diameter of a Eurorack system
	// mounting hole, in millimetres
	MountingHoleDiameter = eurorack.MountingHoleDiameter

	// HP represents horizontal pitch in a Eurorack frame, in millimetres
	HP = eurorack.HP
)

// Intellijel implements the panel.Panel interface and encapsulates the physical
// characteristics of a Intellijel panel
type Intellijel struct {
	HP int
}

// NewIntellijel constructs a new Intellijel object
func NewIntellijel(hp int) *Intellijel {
	return &Intellijel{HP: hp}
}

// Width returns the width of a Intellijel panel, in millimetres
func (e Intellijel) Width() float64 {
	return HP * float64(e.HP)
}

// Height returns the height of a Intellijel panel, in millimetres
func (e Intellijel) Height() float64 {
	return PanelHeight1U
}

// MountingHoleDiameter returns the Intellijel system mounting hole size, in
// millimetres
func (e Intellijel) MountingHoleDiameter() float64 {
	return MountingHoleDiameter
}

// MountingHoles generates a set of Point objects representing the mounting
// hole locations of a Intellijel panel
func (e Intellijel) MountingHoles() []panel.Point {
	rhsx := MountingHolesLeftOffset + HP*(float64(e.HP-3))
	holes := []panel.Point{
		{X: MountingHolesLeftOffset, Y: MountingHoleBottomY1U},
		{X: MountingHolesLeftOffset, Y: MountingHoleTopY1U},
		{X: rhsx, Y: MountingHoleBottomY1U},
		{X: rhsx, Y: MountingHoleTopY1U},
	}
	return holes
}
