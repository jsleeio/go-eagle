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

package eagle

// Grid object
type Grid struct {
	Distance        float64 `xml:"distance,attr"`
	UnitDistance    string  `xml:"unitdist,attr"`
	Unit            string  `xml:"unit,attr"`
	Style           string  `xml:"style,attr"`
	Multiple        float64 `xml:"multiple,attr"`
	Display         string  `xml:"display,attr"`
	AltDistance     float64 `xml:"altdistance,attr"`
	AltUnitDistance string  `xml:"altunitdist,attr"`
	AltUnit         string  `xml:"altunit,attr"`
}

// Layer object
type Layer struct {
	Number  int    `xml:"number,attr"`
	Name    string `xml:"name,attr"`
	Color   int    `xml:"color,attr"`
	Fill    int    `xml:"fill,attr"`
	Visible string `xml:"visible,attr"`
	Active  string `xml:"active,attr"`
}

// Wire object
type Wire struct {
	X1    float64 `xml:"x1,attr"`
	Y1    float64 `xml:"y1,attr"`
	X2    float64 `xml:"x2,attr"`
	Y2    float64 `xml:"y2,attr"`
	Width float64 `xml:"width,attr"`
	Layer int     `xml:"layer,attr"`
	Style string  `xml:"style,attr,omitempty"`
	Cap   string  `xml:"cap,attr,omitempty"`
	Curve float64 `xml:"curve,attr,omitempty"`
}

// Rectangle object
type Rectangle struct {
	X1     float64 `xml:"x1,attr"`
	Y1     float64 `xml:"y1,attr"`
	X2     float64 `xml:"x2,attr"`
	Y2     float64 `xml:"y2,attr"`
	Layer  int     `xml:"layer,attr"`
	Rotate string  `xml:"rot,attr,omitempty"`
}

// Vertex object, used only in Polygon
type Vertex struct {
	X     float64 `xml:"x,attr"`
	Y     float64 `xml:"y,attr"`
	Curve float64 `xml:"curve,attr,omitempty"`
}

// Polygon object
type Polygon struct {
	Vertices []Vertex `xml:"vertices>vertex"`
	Isolate  string   `xml:"isolate,omitempty"`
	Pour     string   `xml:"pour,omitempty"`
	Orphans  string   `xml:"orphans,omitempty"`
	Layer    int      `xml:"layer,attr"`
	Rank     int      `xml:"rank,attr,omitempty"`
	Spacing  string   `xml:"spacing,attr,omitempty"`
	Thermals string   `xml:"thermals,attr,omitempty"`
	Width    float64  `xml:"width,attr"`
}

// Hole object
type Hole struct {
	X     float64 `xml:"x,attr"`
	Y     float64 `xml:"y,attr"`
	Drill float64 `xml:"drill,attr"`
}

// Circle object
type Circle struct {
	X      float64 `xml:"x,attr"`
	Y      float64 `xml:"y,attr"`
	Radius float64 `xml:"radius,attr"`
	Width  float64 `xml:"width,attr"`
	Layer  int     `xml:"layer,attr"`
}

// Pad object
type Pad struct {
	Name   string  `xml:"name,attr"`
	X      float64 `xml:"x,attr"`
	Y      float64 `xml:"y,attr"`
	Drill  float64 `xml:"drill,attr"`
	Shape  string  `xml:"shape,attr,omitempty"`
	Rotate string  `xml:"rot,attr,omitempty"`
}

// Text object
type Text struct {
	Text     string  `xml:",chardata"`
	X        float64 `xml:"x,attr"`
	Y        float64 `xml:"y,attr"`
	Size     float64 `xml:"size,attr"`
	Layer    int     `xml:"layer,attr"`
	Ratio    int     `xml:"ratio,attr,omitempty"`
	Font     string  `xml:"font,attr,omitempty"`
	Align    string  `xml:"align,attr,omitempty"`
	Distance float64 `xml:"distance,attr,omitempty"`
	Rotate   string  `xml:"rot,attr,omitempty"`
}

// Package object
type Package struct {
	Name        string      `xml:"name,attr"`
	Urn         string      `xml:"urn,attr,omitempty"`
	Version     string      `xml:"library_version,attr,omitempty"`
	Description string      `xml:"description"`
	Pads        []Pad       `xml:"pad"`
	Rectangles  []Rectangle `xml:"rectangle"`
	Circles     []Circle    `xml:"circle"`
	Texts       []Text      `xml:"text"`
	Wires       []Wire      `xml:"wire"`
}

// PackageInstance object
// FIXME: should be able to eliminate this and directly include in
// Package3D object via a slice
type PackageInstance struct {
	Name string `xml:"name,attr"`
}

// Package3D object
type Package3D struct {
	Name        string            `xml:"name,attr"`
	Urn         string            `xml:"urn,attr,omitempty"`
	Type        string            `xml:"type,attr"`
	Version     string            `xml:"library_version,attr,omitempty"`
	Description string            `xml:"description"`
	Instances   []PackageInstance `xml:"packageinstances>packageinstance"`
}

// Library object
type Library struct {
	Name        string    `xml:"name,attr"`
	Urn         string    `xml:"urn,attr,omitempty"`
	Description string    `xml:"description"`
	Packages    []Package `xml:"packages>package"`
}

// Attribute object
type Attribute struct {
	Name     string  `xml:"name,attr"`
	Value    string  `xml:"value,attr"`
	X        float64 `xml:"x,attr,omitempty"`
	Y        float64 `xml:"y,attr,omitempty"`
	Size     float64 `xml:"size,attr,omitempty"`
	Display  string  `xml:"display,attr,omitempty"`
	Constant string  `xml:"constant,attr,omitempty"`
	Font     string  `xml:"font,attr,omitempty"`
	Ratio    int     `xml:"ratio,attr,omitempty"`
	Rotate   string  `xml:"rot,attr,omitempty"`
}

// Element object
type Element struct {
	Name         string      `xml:"name,attr"`
	Value        string      `xml:"value,attr"`
	X            float64     `xml:"x,attr"`
	Y            float64     `xml:"y,attr"`
	Library      string      `xml:"library,attr"`
	LibraryUrn   string      `xml:"library_urn,attr,omitempty"`
	Package      string      `xml:"package,attr"`
	Package3dUrn string      `xml:"package3d_urn,attr,omitempty"`
	Smashed      string      `xml:"smashed,attr"`
	Rotate       string      `xml:"rot,attr,omitempty"`
	Attributes   []Attribute `xml:"attribute"`
}

func (e Element) GetAttributes() []Attribute {
	return e.Attributes
}

// Plain object
type Plain struct {
	Holes      []Hole      `xml:"hole"`
	Circles    []Circle    `xml:"circle"`
	Rectangles []Rectangle `xml:"rectangle"`
	Polygons   []Polygon   `xml:"polygon"`
	Texts      []Text      `xml:"text"`
	Wires      []Wire      `xml:"wire"`
}

// NewPlain constructs a new Plain object with empty slices; this
// should be substantially more convenient to use.
func NewPlain() Plain {
	return Plain{
		Holes:      []Hole{},
		Circles:    []Circle{},
		Rectangles: []Rectangle{},
		Polygons:   []Polygon{},
		Texts:      []Text{},
		Wires:      []Wire{},
	}
}

// Board object
type Board struct {
	Libraries  []Library   `xml:"libraries>library"`
	Elements   []Element   `xml:"elements>element"`
	Plain      Plain       `xml:"plain"`
	Attributes []Attribute `xml:"attributes>attribute"`
}

func (b Board) GetAttributes() []Attribute {
	return b.Attributes
}

// NewBoard constructs a new empty Board object.
func NewBoard() Board {
	return Board{
		Libraries:  []Library{},
		Elements:   []Element{},
		Attributes: []Attribute{},
		Plain:      NewPlain(),
	}
}

// Eagle object
type Eagle struct {
	Version string  `xml:"version,attr"`
	Grid    Grid    `xml:"drawing>grid"`
	Layers  []Layer `xml:"drawing>layers>layer"`
	Board   Board   `xml:"drawing>board"`
}
