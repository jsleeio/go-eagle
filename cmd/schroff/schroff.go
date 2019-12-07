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
	"regexp"

	"github.com/jsleeio/go-eagle/pkg/eagle"
	"github.com/jsleeio/go-eagle/pkg/format/eurorack"
	"github.com/jsleeio/go-eagle/pkg/format/intellijel"
	"github.com/jsleeio/go-eagle/pkg/format/pulplogic"
	"github.com/jsleeio/go-eagle/pkg/panel"

	"github.com/jsleeio/go-eagle/internal/boardops/standard"
	"github.com/jsleeio/go-eagle/internal/outline"
)

const (
	// FormatEurorack is the Doepfer-defined 3U specification. Not Eurocard!
	FormatEurorack = "eurorack"
	// FormatPulplogic is the PulpLogic-defined 1U specification
	FormatPulplogic = "pulplogic"
	// FormatIntellijel is the Intellijel-defined 1U specification
	FormatIntellijel = "intellijel"
)

// wrap up all of the context required for creating panel features
// into one place to simplify and reduce error
type panelLayoutContext struct {
	bc    outline.BoardCoords
	board *eagle.Eagle
	panel *eagle.Eagle
	cfg   config
	spec  panel.Panel
	// legendSkipRe is pulled from the board global attribute PANEL_LEGEND_SKIP_RE.
	// If a component name matches this regexp, it will NOT have a panel legend
	// text object created.
	legendSkipRe *regexp.Regexp
	legendLayer  string
	headerLayer  string
	footerLayer  string
}

func setupPanelLayoutContext(board *eagle.Eagle, c config) (panelLayoutContext, error) {
	plc := panelLayoutContext{
		cfg:          c,
		board:        board,
		bc:           outline.DeriveBoardCoords(board),
		legendLayer:  eagle.AttributeString(board.Board, "PANEL_LEGEND_LAYER", "tStop"),
		headerLayer:  eagle.AttributeString(board.Board, "PANEL_HEADER_LAYER", "tStop"),
		footerLayer:  eagle.AttributeString(board.Board, "PANEL_FOOTER_LAYER", "tStop"),
		legendSkipRe: nil,
	}
	if lsre := eagle.AttributeString(board.Board, "PANEL_LEGEND_SKIP_RE", ""); lsre != "" {
		plc.legendSkipRe = regexp.MustCompile(lsre)
	}
	spec, err := panelSpecForFormat(plc.bc.HP, *plc.cfg.Format)
	if err != nil {
		return panelLayoutContext{}, err
	}
	plc.spec = spec
	plc.panel = plc.board.CloneEmpty()
	if err := standard.ApplyStandardBoardOperations(plc.panel, plc.spec); err != nil {
		return panelLayoutContext{}, fmt.Errorf("error creating panel features: %v", err)
	}
	// centre the board on the panel
	plc.bc.XOffset += (plc.spec.Width()-plc.bc.Width())/2 + plc.spec.HorizontalFit()/2
	plc.bc.YOffset += (plc.spec.Height() - plc.bc.Height()) / 2
	return plc, nil
}

func panelSpecForFormat(width int, format string) (panel.Panel, error) {
	var spec panel.Panel
	switch format {
	case FormatEurorack:
		spec = eurorack.NewEurorack(width)
	case FormatPulplogic:
		spec = pulplogic.NewPulplogic(width)
	case FormatIntellijel:
		spec = intellijel.NewIntellijel(width)
	default:
		return nil, fmt.Errorf("unsupported format: %s", format)
	}
	return spec, nil
}

func headerOp(plc panelLayoutContext) {
	offsets := make(map[string]float64)
	offsets["PANEL_HEADER_OFFSET_X"] = 0.0
	offsets["PANEL_HEADER_OFFSET_Y"] = 0.0
	offsets["PANEL_FOOTER_OFFSET_X"] = 0.0
	offsets["PANEL_FOOTER_OFFSET_Y"] = 0.0
	for k, defval := range offsets {
		if v, err := eagle.AttributeFloat(plc.board.Board, k, defval); err == nil {
			offsets[k] = v
		} else {
			log.Fatalf("invalid global attribute numeric value: %s: %v", err)
		}
	}
	// add the header and footer
	header := eagle.Text{
		X:     plc.spec.Width()/2.0 + offsets["PANEL_HEADER_OFFSET_X"],
		Y:     plc.spec.MountingHoleTopY() + offsets["PANEL_HEADER_OFFSET_Y"],
		Align: "center",
		Size:  3.0,
		Text:  eagle.AttributeString(plc.board.Board, "PANEL_HEADER_TEXT", "<HEADER>"),
		Layer: plc.panel.LayerByName(plc.headerLayer),
	}
	plc.panel.Board.Plain.Texts = append(plc.panel.Board.Plain.Texts, header)
	footer := eagle.Text{
		X:     plc.spec.Width()/2.0 + offsets["PANEL_FOOTER_OFFSET_X"],
		Y:     plc.spec.MountingHoleBottomY() + offsets["PANEL_FOOTER_OFFSET_Y"],
		Align: "center",
		Size:  3.0,
		Text:  eagle.AttributeString(plc.board.Board, "PANEL_FOOTER_TEXT", "<FOOTER>"),
		Layer: plc.panel.LayerByName(plc.footerLayer),
	}
	plc.panel.Board.Plain.Texts = append(plc.panel.Board.Plain.Texts, footer)
}

func elementOp(plc panelLayoutContext, elem eagle.Element) {
	hole, needHole, err := holeForPanelElement(elem)
	if err != nil {
		log.Fatalf("can't find drill size for element %q: %v", elem.Name, err)
	}
	if needHole {
		// the hole was generated with coordinates from the source board, now
		// adjust them to be in the right place on the panel
		tstop := plc.panel.LayerByName("tStop")
		hole.X += plc.bc.XOffset
		hole.Y += plc.bc.YOffset
		aox, err := eagle.AttributeFloat(elem, "PANEL_LEGEND_OFFSET_X", 0.0)
		if err != nil {
			log.Fatal(err)
		}
		aoy, err := eagle.AttributeFloat(elem, "PANEL_LEGEND_OFFSET_Y", 0.0)
		if err != nil {
			log.Fatal(err)
		}
		plc.panel.Board.Plain.Holes = append(plc.panel.Board.Plain.Holes, hole)
		text := eagle.Text{
			X:     aox + hole.X,
			Y:     aoy + hole.Y + (hole.Drill / 2.0) + *plc.cfg.TextSpacing,
			Size:  *plc.cfg.TextSize,
			Layer: plc.panel.LayerByName(plc.legendLayer),
			Text:  eagle.AttributeString(elem, "PANEL_LEGEND", elem.Name),
			Align: "bottom-center",
			Font:  "vector",
		}
		if text.Text != "" && (plc.legendSkipRe == nil || !plc.legendSkipRe.MatchString(elem.Name)) {
			plc.panel.Board.Plain.Texts = append(plc.panel.Board.Plain.Texts, text)
		} else {
			log.Printf("%s: skipping legend\n", elem.Name)
		}
		hsw, err := eagle.AttributeFloat(elem, "PANEL_HOLE_STOP_WIDTH", *plc.cfg.HoleStopRadius)
		if err != nil {
			log.Fatal(err)
		}
		stop := eagle.Circle{
			X:      hole.X,
			Y:      hole.Y,
			Radius: hole.Drill / 2.0,
			Width:  hsw,
			Layer:  tstop,
		}
		plc.panel.Board.Plain.Circles = append(plc.panel.Board.Plain.Circles, stop)
	}
}

// generate a panel hole for a single element, if necessary
func holeForPanelElement(elem eagle.Element) (eagle.Hole, bool, error) {
	hole := eagle.Hole{X: elem.X, Y: elem.Y}
	// negative default drill size => no drill unless PANEL_DRILL_MM present
	drillmm, err := eagle.AttributeFloat(elem, "PANEL_DRILL_MM", -1.0)
	if err != nil {
		return eagle.Hole{}, false, err
	}
	if drillmm < 0.0 { // no drill size found, do nothing
		return eagle.Hole{}, false, nil
	}
	hole.Drill = drillmm
	log.Printf("%s: found PANEL_DRILL_MM attribute with value %v", elem.Name, drillmm)
	return hole, true, nil
}

type config struct {
	Format         *string
	TextSpacing    *float64
	TextSize       *float64
	HoleStopRadius *float64
}

func configureFromFlags() config {
	cfg := config{
		Format:         flag.String("format", FormatEurorack, "panel format to create (eurorack, pulplogic, intellijel)"),
		TextSpacing:    flag.Float64("text-spacing", 3.5, "spacing between a hole and its related label"),
		TextSize:       flag.Float64("text-size", 2.25, "label text size"),
		HoleStopRadius: flag.Float64("hole-stop-radius", 2.0, "Radius to pull back soldermask around a hole"),
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
		// panel, bc := schroffPanelForBoard(board, config)
		plc, err := setupPanelLayoutContext(board, config)
		if err != nil {
			log.Fatalf("can't setup panel layout context: %v", err)
		}
		headerOp(plc)
		for _, elem := range plc.board.Board.Elements {
			elementOp(plc, elem)
		}
		outFilename := filepath.Base(filename) + ".panel.brd"
		if err := plc.panel.WriteFile(outFilename); err != nil {
			log.Fatalf("can't write output file %q: %v", outFilename, err)
		}
	}
}
