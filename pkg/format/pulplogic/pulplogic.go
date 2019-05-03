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

package pulplogic

import (
	"github.com/jsleeio/go-eagle/pkg/panel"
)

// based on http://pulplogic.com/1u_tiles/

const (
	inch = 25.4

	// PanelHeight1U represents the total height of a Pulplogic panel, in
	// millimetres
	PanelHeight1U = 1.70 * inch

	// MountingHolesLeftOffset represents the distance of the first mounting
	// hole from the left edge of the panel, in millimetres
	MountingHolesLeftOffset = 0.2 * inch

	// MountingHolesRightOffset represents the distance of the first mounting
	// hole from the right edge of the panel, in millimetres
	MountingHolesRightOffset = 0.2 * inch

	// MountingHoleTopY1U represents the Y value for the top row of 1U mounting
	// holes, in millimetres
	MountingHoleTopY1U = PanelHeight1U - (0.118 * inch)

	// MountingHoleBottomY1U represents the Y value for the bottom row of 1U
	// mounting holes, in millimetres
	MountingHoleBottomY1U = 0.118 * inch

	// MountingHoleDiameter represents the diameter of a Eurorack system
	// mounting hole, in millimetres
	MountingHoleDiameter = 0.125 * inch

	// HP represents horizontal pitch in a Eurorack frame, in millimetres
	HP = 5.08
)

// Pulplogic implements the panel.Panel interface and encapsulates the physical
// characteristics of a Pulplogic panel
type Pulplogic struct {
	HP int
}

// NewPulplogic constructs a new Pulplogic object
func NewPulplogic(hp int) *Pulplogic {
	return &Pulplogic{HP: hp}
}

// Width returns the width of a Pulplogic panel, in millimetres
func (e Pulplogic) Width() float64 {
	return HP * float64(e.HP)
}

// Height returns the height of a Pulplogic panel, in millimetres
func (e Pulplogic) Height() float64 {
	return PanelHeight1U
}

// MountingHoleDiameter returns the Pulplogic system mounting hole size, in
// millimetres
func (e Pulplogic) MountingHoleDiameter() float64 {
	return MountingHoleDiameter
}

// MountingHoles generates a set of Point objects representing the mounting
// hole locations of a Pulplogic panel
func (e Pulplogic) MountingHoles() []panel.Point {
	holes := []panel.Point{
		{X: MountingHolesLeftOffset, Y: MountingHoleBottomY1U},
		{X: MountingHolesLeftOffset, Y: MountingHoleTopY1U},
		{X: e.Width() - MountingHolesRightOffset, Y: MountingHoleBottomY1U},
		{X: e.Width() - MountingHolesRightOffset, Y: MountingHoleTopY1U},
	}
	return holes
}
