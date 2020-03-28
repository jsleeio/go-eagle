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

package eurorack

import (
	"github.com/jsleeio/go-eagle/pkg/panel"
)

const (
	// PanelHeight3U represents the total height of a Eurorack panel. Note in
	// particular that this is NOT the same as the Eurocard standard, as the
	// latter does not use lipped rails
	PanelHeight3U = 128.5

	// ExtraMountingHolesThreshold represents the panel width threshold beyond
	// which additional mounting holes are required
	ExtraMountingHolesThreshold = 8

	// MountingHolesLeftOffset represents the distance of the first mounting
	// hole from the left edge of the panel
	MountingHolesLeftOffset = 7.5

	// MountingHoleTopY3U represents the Y value for the top row of 3U mounting
	// holes
	MountingHoleTopY3U = PanelHeight3U - 3.00

	// MountingHoleBottomY3U represents the Y value for the bottom row of 3U
	// mounting holes
	MountingHoleBottomY3U = 3.00

	// MountingHoleDiameter represents the diameter of a Eurorack system
	// mounting hole, in millimetres
	MountingHoleDiameter = 3.2

	// HP represents horizontal pitch in a Eurorack frame, in millimetres
	HP = 5.08

	// HorizontalFit indicates the panel tolerance adjustment for the format
	HorizontalFit = 0.25

	// RailHeightFromMountingHole is used to determine how much space exists.
	// See discussion in github.com/jsleeio/pkg/panel. 5mm is a good safe
	// figure for all known-used Eurorack rail types
	RailHeightFromMountingHole = 5.0
)

// Eurorack implements the panel.Panel interface and encapsulates the physical
// characteristics of a Eurorack panel
type Eurorack struct {
	HP int
}

// NewEurorack constructs a new Eurorack object
func NewEurorack(hp int) *Eurorack {
	return &Eurorack{HP: hp}
}

// Width returns the width of a Eurorack panel, in millimetres
func (e Eurorack) Width() float64 {
	return HP * float64(e.HP)
}

// Height returns the height of a Eurorack panel, in millimetres
func (e Eurorack) Height() float64 {
	return PanelHeight3U
}

// MountingHoleDiameter returns the Eurorack system mounting hole size, in
// millimetres
func (e Eurorack) MountingHoleDiameter() float64 {
	return MountingHoleDiameter
}

// MountingHoles generates a set of Point objects representing the mounting
// hole locations of a Eurorack panel
func (e Eurorack) MountingHoles() []panel.Point {
	holes := []panel.Point{
		{X: MountingHolesLeftOffset, Y: MountingHoleBottomY3U},
		{X: MountingHolesLeftOffset, Y: MountingHoleTopY3U},
	}
	if e.HP > ExtraMountingHolesThreshold {
		rhsx := MountingHolesLeftOffset + HP*(float64(e.HP-3))
		holes = append(holes, panel.Point{X: rhsx, Y: MountingHoleBottomY3U})
		holes = append(holes, panel.Point{X: rhsx, Y: MountingHoleTopY3U})
	}
	return holes
}

// HorizontalFit indicates the panel tolerance adjustment for the format
func (e Eurorack) HorizontalFit() float64 {
	return HorizontalFit
}

// RailHeightFromMountingHole is used to calculate space between rails
func (e Eurorack) RailHeightFromMountingHole() float64 {
	return RailHeightFromMountingHole
}

// MountingHoleTopY returns the Y coordinate for the top row of mounting
// holes
func (e Eurorack) MountingHoleTopY() float64 {
	return MountingHoleTopY3U
}

// MountingHoleBottomY returns the Y coordinate for the bottom row of
// mounting holes
func (e Eurorack) MountingHoleBottomY() float64 {
	return MountingHoleBottomY3U
}

// HeaderLocation returns the location of the header text. Eurorack has
// mounting rails so this is typically aligned with the top mounting screw
func (e Eurorack) HeaderLocation() panel.Point {
	return panel.Point{X: e.Width() / 2, Y: e.MountingHoleTopY()}
}

// FooterLocation returns the location of the footer text. Eurorack has
// mounting rails so this is typically aligned with the bottom mounting screw
func (e Eurorack) FooterLocation() panel.Point {
	return panel.Point{X: e.Width() / 2, Y: e.MountingHoleBottomY()}
}
