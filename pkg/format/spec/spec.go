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

package spec

import (
	"fmt"
	"io/ioutil"
	"sort"

	"gopkg.in/yaml.v2"

	"github.com/jsleeio/go-eagle/pkg/panel"
)

// Spec implements the panel.Panel interface and encapsulates the physical
// characteristics of a Spec panel
type Spec struct {
	SpecName                 string        `yaml:"name"`
	SpecWidth                float64       `yaml:"width"`
	SpecHeight               float64       `yaml:"height"`
	SpecMountingHoles        []panel.Point `yaml:"mountingHoles"`
	SpecMountingHoleDiameter float64       `yaml:"mountingHoleDiameter"`
	SpecHorizontalFit        float64       `yaml:"horizontalFit"`
}

type PanelSpecError struct {
	s string
}

func (e *PanelSpecError) Error() string {
	return fmt.Sprintf("PanelSpecError: %s", e.s)
}

func NewPanelSpecError(s string) error {
	return &PanelSpecError{s: s}
}

// LoadSpec constructs a new Spec object according to a YAML file definition
func LoadSpec(filename string) (*Spec, error) {
	yamltext, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var sp Spec
	if err := yaml.Unmarshal(yamltext, &sp); err != nil {
		return nil, err
	}
	if len(sp.SpecMountingHoles) < 1 {
		return nil, NewPanelSpecError("need at least one mounting hole")
	}
	sort.Slice(sp.SpecMountingHoles, func(i, j int) bool {
		return sp.SpecMountingHoles[i].Y < sp.SpecMountingHoles[j].Y
	})
	return &sp, nil
}

// Width returns the width of a Spec panel, in millimetres
func (s Spec) Width() float64 {
	return s.SpecWidth
}

// Height returns the height of a Spec panel, in millimetres
func (s Spec) Height() float64 {
	return s.SpecHeight
}

// MountingHoleDiameter returns the Spec system mounting hole size, in
// millimetres
func (s Spec) MountingHoleDiameter() float64 {
	return s.SpecMountingHoleDiameter
}

// MountingHoles generates a set of Point objects representing the mounting
// hole locations of a Spec panel
func (s Spec) MountingHoles() []panel.Point {
	return s.SpecMountingHoles
}

// HorizontalFit indicates the panel tolerance adjustment for the format
func (s Spec) HorizontalFit() float64 {
	return s.SpecHorizontalFit
}

// RailHeightFromMountingHole doesn't really directly apply to YAML-spec
// panels where the keepout area is more likely to be a ring around each
// mounting hole, with the thickness of the ring likely varying with the
// mounting hole diameter. For now, keep it simple and return 0. The PCB
// designer will need to be careful regardless in order to fit their
// design within the given enclosure's envelope.
//
// FIXME: fix the interface design here to better facilitate non-system
//        panels (without making it unnecessarily awful!)
func (s Spec) RailHeightFromMountingHole() float64 {
	return 0.0
}

// MountingHoleTopY returns the Y coordinate for the top row of mounting
// holes
func (s Spec) MountingHoleTopY() float64 {
	return s.SpecMountingHoles[0].Y
}

// MountingHoleBottomY returns the Y coordinate for the bottom row of
// mounting holes
func (s Spec) MountingHoleBottomY() float64 {
	return s.SpecMountingHoles[len(s.SpecMountingHoles)-1].Y
}

// HeaderLocation returns the location of the header text. Spec panels
// may not have mounting rails so this is entirely arbitrary
func (s Spec) HeaderLocation() panel.Point {
	return panel.Point{X: s.Width() / 2, Y: s.MountingHoleTopY()}
}

// FooterLocation returns the location of the footer text. Spec panels
// may not have mounting rails so this is entirely arbitrary
func (s Spec) FooterLocation() panel.Point {
	return panel.Point{X: s.Width() / 2, Y: s.MountingHoleBottomY()}
}
