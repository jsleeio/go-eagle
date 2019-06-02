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
	"strconv"

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
		legendLayer:  "tStop",
		headerLayer:  "tStop",
		footerLayer:  "tStop",
		legendSkipRe: nil,
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
	if legendSkipRe, ok := plc.board.Board.AttributeByName("PANEL_LEGEND_SKIP_RE"); ok {
		if legendSkipRe != "" {
			plc.legendSkipRe = regexp.MustCompile(legendSkipRe)
		}
	}
	if headerLayer, ok := plc.board.Board.AttributeByName("PANEL_HEADER_LAYER"); ok {
		plc.headerLayer = headerLayer
	}
	if footerLayer, ok := plc.board.Board.AttributeByName("PANEL_FOOTER_LAYER"); ok {
		plc.footerLayer = footerLayer
	}
	if legendLayer, ok := plc.board.Board.AttributeByName("PANEL_LEGEND_LAYER"); ok {
		plc.legendLayer = legendLayer
	}
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
	// add the header and footer
	headertext, headerok := plc.board.Board.AttributeByName("PANEL_HEADER_TEXT")
	header := eagle.Text{
		X:     plc.spec.Width() / 2.0,
		Y:     plc.spec.MountingHoleTopY(),
		Align: "center",
		Size:  3.0,
		Text:  "<<HEADER_TEXT>>",
		Layer: plc.panel.LayerByName(plc.headerLayer),
	}
	if headerok {
		log.Printf("board: found PANEL_HEADER_TEXT attribute with value %q", headertext)
		header.Text = headertext
	}
	plc.panel.Board.Plain.Texts = append(plc.panel.Board.Plain.Texts, header)
	footertext, footerok := plc.board.Board.AttributeByName("PANEL_FOOTER_TEXT")
	footer := eagle.Text{
		X:     plc.spec.Width() / 2.0,
		Y:     plc.spec.MountingHoleBottomY(),
		Align: "center",
		Size:  3.0,
		Text:  "<<FOOTER_TEXT>>",
		Layer: plc.panel.LayerByName(plc.footerLayer),
	}
	if footerok {
		log.Printf("board: found PANEL_FOOTER_TEXT attribute with value %q", footertext)
		footer.Text = footertext
	}
	plc.panel.Board.Plain.Texts = append(plc.panel.Board.Plain.Texts, footer)
}

func elementFloatAttribute(elem eagle.Element, attribute string, def float64) (float64, error) {
	if s, ok := elem.AttributeByName(attribute); ok {
		if f, err := strconv.ParseFloat(s, 64); err != nil {
			return 0, fmt.Errorf("unparseable numeric attribute %q: %v", attribute, err)
		} else {
			return f, nil
		}
	}
	return def, nil
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
		aox, err := elementFloatAttribute(elem, "PANEL_LEGEND_OFFSET_X", 0.0)
		if err != nil {
			log.Fatalf("element %s has an unparseable PANEL_LEGEND_OFFSET_X attribute: %v", elem.Name, err)
		}
		aoy, err := elementFloatAttribute(elem, "PANEL_LEGEND_OFFSET_Y", 0.0)
		if err != nil {
			log.Fatalf("element %s has an unparseable PANEL_LEGEND_OFFSET_Y attribute: %v", elem.Name, err)
		}
		plc.panel.Board.Plain.Holes = append(plc.panel.Board.Plain.Holes, hole)
		text := eagle.Text{
			X:     aox + hole.X,
			Y:     aoy + hole.Y + (hole.Drill / 2.0) + *plc.cfg.TextSpacing,
			Size:  *plc.cfg.TextSize,
			Layer: plc.panel.LayerByName(plc.legendLayer),
			Text:  elem.Name,
			Align: "bottom-center",
			Font:  "vector",
		}
		if legend, ok := elem.AttributeByName("PANEL_LEGEND"); ok {
			text.Text = legend
		}
		if text.Text != "" && (plc.legendSkipRe == nil || !plc.legendSkipRe.MatchString(elem.Name)) {
			plc.panel.Board.Plain.Texts = append(plc.panel.Board.Plain.Texts, text)
		}
		hsw := *plc.cfg.HoleStopRadius
		if hsws, ok := elem.AttributeByName("PANEL_HOLE_STOP_WIDTH"); ok {
			if hswf, err := strconv.ParseFloat(hsws, 64); err != nil {
				log.Fatalf("element %s has an unparseable PANEL_HOLE_STOP_WIDTH attribute %q: %v", elem.Name, hsws, err)
			} else {
				hsw = hswf
			}
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
	return eagle.Hole{}, false, nil
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
