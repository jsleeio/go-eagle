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
	"github.com/jsleeio/go-eagle/pkg/format/eurorack"
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
	HP = eurorack.HP

	// HorizontalFit indicates the panel tolerance adjustment for the format
	HorizontalFit = eurorack.HorizontalFit

	// RailHeightFromMountingHole is used to determine how much space exists.
	// See discussion in github.com/jsleeio/pkg/panel.
	//
	// NOT using the Eurorack figure here; see the Pulplogic spec. Instead we use
	// half the height of Vector T-strut rails. With this measurement, the
	// Pulplogic-recommended maximum PCB size (1.130") will fit between a pair of
	// keepout areas extending this distance beyond the mounting hole centres.
	RailHeightFromMountingHole = (0.291 / 2.0) * inch
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
func (p Pulplogic) Width() float64 {
	return HP * float64(p.HP)
}

// Height returns the height of a Pulplogic panel, in millimetres
func (p Pulplogic) Height() float64 {
	return PanelHeight1U
}

// MountingHoleDiameter returns thp Pulplogic system mounting hole size, in
// millimetres
func (p Pulplogic) MountingHoleDiameter() float64 {
	return MountingHoleDiameter
}

// MountingHoles generates a set of Point objects representing the mounting
// hole locations of a Pulplogic panel
func (p Pulplogic) MountingHoles() []panel.Point {
	holes := []panel.Point{
		{X: MountingHolesLeftOffset, Y: MountingHoleBottomY1U},
		{X: MountingHolesLeftOffset, Y: MountingHoleTopY1U},
		{X: p.Width() - MountingHolesRightOffset, Y: MountingHoleBottomY1U},
		{X: p.Width() - MountingHolesRightOffset, Y: MountingHoleTopY1U},
	}
	return holes
}

// HorizontalFit indicates the panel tolerance adjustment for the format
func (p Pulplogic) HorizontalFit() float64 {
	return HorizontalFit
}

// RailHeightFromMountingHole is used to calculate space between rails
func (p Pulplogic) RailHeightFromMountingHole() float64 {
	return RailHeightFromMountingHole
}

// MountingHoleTopY returns the Y coordinate for the top row of mounting
// holes
func (p Pulplogic) MountingHoleTopY() float64 {
	return MountingHoleTopY1U
}

// MountingHoleBottomY returns the Y coordinate for the bottom row of
// mounting holes
func (p Pulplogic) MountingHoleBottomY() float64 {
	return MountingHoleBottomY1U
}

// HeaderLocation returns the location of the header text. Pulplogic 1U has
// mounting rails so this is typically aligned with the top mounting screw
func (p Pulplogic) HeaderLocation() panel.Point {
	return panel.Point{X: p.Width() / 2.0, Y: p.MountingHoleTopY()}
}

// FooterLocation returns the location of the footer text. Pulplogic 1U has
// mounting rails so this is typically aligned with the bottom mounting screw
func (p Pulplogic) FooterLocation() panel.Point {
	return panel.Point{X: p.Width() / 2.0, Y: p.MountingHoleBottomY()}
}
