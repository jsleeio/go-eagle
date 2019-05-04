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

package main

import (
	"flag"
	"fmt"
	"log"
	"path/filepath"
	"strconv"

	"github.com/jsleeio/go-eagle/internal/outline"
	"github.com/jsleeio/go-eagle/pkg/eagle"
)

func printElements(e *eagle.Eagle) {
	for _, elem := range e.Board.Elements {
		fmt.Printf("%s is a %s::%s at (%v,%v)\n",
			elem.Name, elem.Library, elem.Package, elem.X, elem.Y)
	}
}

func schroffSystemHoles(hp int) []eagle.Hole {
	holes := []eagle.Hole{
		eagle.Hole{X: 7.5, Y: 3.0, Drill: 3.2},
		eagle.Hole{X: 7.5, Y: 125.5, Drill: 3.2},
	}
	if hp > 8 {
		// http://www.doepfer.de/a100_man/a100m_e.htm
		rhsX := 7.5 + (5.08 * float64(hp-3))
		holes = append(holes, eagle.Hole{X: rhsX, Y: 3.0, Drill: 3.2})
		holes = append(holes, eagle.Hole{X: rhsX, Y: 125.5, Drill: 3.2})
	}
	return holes
}

func wireRectangle(x1, y1, x2, y2 float64, layer int) []eagle.Wire {
	return []eagle.Wire{
		{X1: x1, Y1: y1, X2: x2, Y2: y1, Layer: layer, Width: 1}, // bottom
		{X1: x1, Y1: y2, X2: x2, Y2: y2, Layer: layer, Width: 1}, // top
		{X1: x1, Y1: y1, X2: x1, Y2: y2, Layer: layer, Width: 1}, // left
		{X1: x2, Y1: y1, X2: x2, Y2: y2, Layer: layer, Width: 1}, // right
	}
}

func schroffPanelForBoard(board *eagle.Eagle, cfg config) (*eagle.Eagle, outline.BoardCoords) {
	bc := outline.DeriveBoardCoords(board)
	panel := board.CloneEmpty()
	// http://www.doepfer.de/a100_man/a100m_e.htm
	// Doepfer spec does includes a table of some width corrections
	// but 0.25mm should be fine for all likely panel sizes really
	panelWidth := 5.08*float64(bc.HP) - 0.25
	panelHeight := 128.5
	// centre the board on the panel
	bc.XOffset += (panelWidth - bc.Width()) / 2
	bc.YOffset += (panelHeight - bc.Height()) / 2
	dimension := board.LayerByName("Dimension")
	// outline of the panel board
	for _, wire := range eagle.BoardOutlineWires(panelWidth, panelHeight, dimension) {
		panel.Board.Plain.Wires = append(panel.Board.Plain.Wires, wire)
	}
	// place a board rectangle in tDocu layer to indicate its alignment
	// with the panel
	tdocu := board.LayerByName("tDocu")
	indicator := wireRectangle(
		bc.XMin+bc.XOffset, bc.YMin+bc.YOffset,
		bc.XMax+bc.XOffset, bc.YMax+bc.YOffset,
		tdocu,
	)
	for _, wire := range indicator {
		panel.Board.Plain.Wires = append(panel.Board.Plain.Wires, wire)
	}
	// place an outline of the rail areas in tKeepout layer to indicate their
	// size and location. The magic number 8 here is derived from the Doepfer
	// spec (system hole centre is 3mm from top or bottom edge of panel) but
	// not defined in it explicitly, as enclosures do not all use the same rails.
	// All known-used Eurorack rails are <= 10mm wide, however, so adding
	// or subtracting 5mm is a pretty good heuristic.
	bRail := eagle.Rectangle{X1: 0, Y1: 0, X2: panelWidth, Y2: 8, Layer: tdocu}
	tRail := eagle.Rectangle{X1: 0, Y1: panelHeight - 8, X2: panelWidth, Y2: panelHeight, Layer: tdocu}
	panel.Board.Plain.Rectangles = append(panel.Board.Plain.Rectangles, bRail)
	panel.Board.Plain.Rectangles = append(panel.Board.Plain.Rectangles, tRail)
	// add the header and footer
	headertext, headerok := board.Board.AttributeByName("HEADER_TEXT")
	header := eagle.Text{
		X:     panelWidth / 2.0,
		Y:     panelHeight - 3.0,
		Align: "center",
		Size:  3.0,
		Text:  "<<HEADER_TEXT>>",
		Layer: board.LayerByName("tSilk"),
	}
	if headerok {
		log.Printf("board: found HEADER_TEXT attribute with value %q", headertext)
		header.Text = headertext
	}
	panel.Board.Plain.Texts = append(panel.Board.Plain.Texts, header)
	footertext, footerok := board.Board.AttributeByName("FOOTER_TEXT")
	footer := eagle.Text{
		X:     panelWidth / 2.0,
		Y:     3.0,
		Align: "center",
		Size:  3.0,
		Text:  "<<FOOTER_TEXT>>",
		Layer: board.LayerByName("tSilk"),
	}
	if footerok {
		log.Printf("board: found FOOTER_TEXT attribute with value %q", footertext)
		footer.Text = footertext
	}
	panel.Board.Plain.Texts = append(panel.Board.Plain.Texts, footer)
	// fill the non-rail areas with copper
	copper := eagle.Rectangle{
		X1:    *cfg.CopperPullback,
		Y1:    8 + *cfg.CopperPullback,
		X2:    panelWidth - *cfg.CopperPullback,
		Y2:    panelHeight - 8 - *cfg.CopperPullback,
		Layer: panel.LayerByName("Top"),
	}
	panel.Board.Plain.Rectangles = append(panel.Board.Plain.Rectangles, copper)
	copper.Layer = panel.LayerByName("Bottom")
	panel.Board.Plain.Rectangles = append(panel.Board.Plain.Rectangles, copper)
	// add the Schroff system holes
	for _, hole := range schroffSystemHoles(bc.HP) {
		panel.Board.Plain.Holes = append(panel.Board.Plain.Holes, hole)
	}
	return panel, bc
}

// generate a panel hole for a single element, if necessary
func holeForPanelElement(elem eagle.Element) (eagle.Hole, bool, error) {
	hole := eagle.Hole{X: elem.X, Y: elem.Y}
	// see if the element has a drill attribute first
	if elemdrill, found := elem.AttributeByName("PANEL_DRILL_MM"); found {
		log.Printf("%s: found PANEL_DRILL_MM attribute with value %q", elem.Name, elemdrill)
		f, err := strconv.ParseFloat(elemdrill, 64)
		if err != nil {
			err = fmt.Errorf("element %s has an unparseable PANEL_DRILL_MM attribute %q: %v", elemdrill, err)
			return eagle.Hole{}, false, err
		}
		hole.Drill = f
		return hole, true, nil
	}
	// if not, check our defaults...
	defaultDrills := map[string]float64{
		// "MusicThingModular::9MM_SNAP-IN_POT":    7.0,
		// "MusicThingModular::WQP-PJ301M-12_JACK": 6.0,
	}
	if defdrill, ok := defaultDrills[elem.Library+"::"+elem.Package]; ok {
		hole.Drill = defdrill
		return hole, true, nil
	}
	return eagle.Hole{}, false, nil
}

type config struct {
	TextSpacing    *float64
	TextSize       *float64
	HoleStopRadius *float64
	CopperPullback *float64
}

func configureFromFlags() config {
	cfg := config{
		TextSpacing:    flag.Float64("text-spacing", 3.5, "spacing between a hole and its related label"),
		TextSize:       flag.Float64("text-size", 2.25, "label text size"),
		HoleStopRadius: flag.Float64("hole-stop-radius", 2.0, "Radius to pull back soldermask around a hole"),
		CopperPullback: flag.Float64("copper-pullback", 0.1, "Distance to pull back the copper pour from the panel edge"),
	}
	flag.Parse()
	return cfg
}

func main() {
	config := configureFromFlags()
	for _, filename := range flag.Args() {
		board, err := eagle.LoadEagleFile(filename)
		if err != nil {
			log.Fatalf("can't load input file %q: %v", filename, err)
		}
		panel, bc := schroffPanelForBoard(board, config)
		tstop := panel.LayerByName("tStop")
		for _, elem := range board.Board.Elements {
			hole, needHole, err := holeForPanelElement(elem)
			if err != nil {
				log.Fatalf("can't find drill size for element %q: %v", elem.Name, err)
			}
			if needHole {
				// offset to allow for centreing the board on the panel
				hole.X += bc.XOffset
				hole.Y += bc.YOffset
				panel.Board.Plain.Holes = append(panel.Board.Plain.Holes, hole)
				text := eagle.Text{
					X:     hole.X,
					Y:     hole.Y + (hole.Drill / 2.0) + *config.TextSpacing,
					Size:  *config.TextSize,
					Layer: tstop,
					Text:  elem.Name,
					Align: "bottom-center",
					Font:  "vector",
				}
				panel.Board.Plain.Texts = append(panel.Board.Plain.Texts, text)
				stop := eagle.Circle{
					X:      hole.X,
					Y:      hole.Y,
					Radius: hole.Drill / 2.0,
					Width:  *config.HoleStopRadius,
					Layer:  tstop,
				}
				panel.Board.Plain.Circles = append(panel.Board.Plain.Circles, stop)
			}
		}
		outFilename := filepath.Base(filename) + ".panel.brd"
		if err := panel.WriteFile(outFilename); err != nil {
			log.Fatalf("can't write output file %q: %v", outFilename, err)
		}
	}
}
