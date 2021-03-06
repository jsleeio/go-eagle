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
	"strings"

	"github.com/jsleeio/go-eagle/pkg/eagle"
	"github.com/jsleeio/go-eagle/pkg/format/eurorack"
	"github.com/jsleeio/go-eagle/pkg/format/intellijel"
	"github.com/jsleeio/go-eagle/pkg/format/pulplogic"
	filespec "github.com/jsleeio/go-eagle/pkg/format/spec"
	"github.com/jsleeio/go-eagle/pkg/geometry"
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
	// FormatSpec is the YAML-derived panel specification
	FormatSpec = "spec"
)

// wrap up all of the context required for creating panel features
// into one place to simplify and reduce error
type panelLayoutContext struct {
	format string
	bc     outline.BoardCoords
	board  *eagle.Eagle
	panel  *eagle.Eagle
	cfg    config
	spec   panel.Panel
	// legendSkipRe is pulled from the board global attribute PANEL_LEGEND_SKIP_RE.
	// If a component name matches this regexp, it will NOT have a panel legend
	// text object created.
	legendSkipRe *regexp.Regexp
	legendLayer  string
	headerLayer  string
	footerLayer  string
}

func (plc *panelLayoutContext) panelSpecForFormat() (err error) {
	err = nil
	switch *plc.cfg.Format {
	case FormatEurorack:
		plc.spec = eurorack.NewEurorack(plc.bc.HP)
	case FormatPulplogic:
		plc.spec = pulplogic.NewPulplogic(plc.bc.HP)
	case FormatIntellijel:
		plc.spec = intellijel.NewIntellijel(plc.bc.HP)
	case FormatSpec:
		plc.spec, err = filespec.LoadSpec(*plc.cfg.SpecFile)
		if err != nil {
			err = fmt.Errorf("error loading YAML panel spec from '%v': %v", *plc.cfg.SpecFile, err)
		}
	default:
		err = fmt.Errorf("unsupported format: %s", *plc.cfg.Format)
	}
	return
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
	err := plc.panelSpecForFormat()
	if err != nil {
		return panelLayoutContext{}, err
	}
	plc.panel = plc.board.CloneEmpty()
	if err := standard.ApplyStandardBoardOperations(plc.panel, plc.spec); err != nil {
		return panelLayoutContext{}, fmt.Errorf("error creating panel features: %v", err)
	}
	// centre the board on the panel
	plc.bc.XOffset += (plc.spec.Width()-plc.bc.Width())/2 + plc.spec.HorizontalFit()/2
	plc.bc.YOffset += (plc.spec.Height() - plc.bc.Height()) / 2
	return plc, nil
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
	headerloc := plc.spec.HeaderLocation()
	header := eagle.Text{
		X:     headerloc.X + offsets["PANEL_HEADER_OFFSET_X"],
		Y:     headerloc.Y + offsets["PANEL_HEADER_OFFSET_Y"],
		Align: "center",
		Size:  3.0,
		Text:  eagle.AttributeString(plc.board.Board, "PANEL_HEADER_TEXT", "<HEADER>"),
		Layer: plc.panel.LayerByName(plc.headerLayer),
	}
	footerloc := plc.spec.FooterLocation()
	plc.panel.Board.Plain.Texts = append(plc.panel.Board.Plain.Texts, header)
	footer := eagle.Text{
		X:     footerloc.X + offsets["PANEL_FOOTER_OFFSET_X"],
		Y:     footerloc.Y + offsets["PANEL_FOOTER_OFFSET_Y"],
		Align: "center",
		Size:  3.0,
		Text:  eagle.AttributeString(plc.board.Board, "PANEL_FOOTER_TEXT", "<FOOTER>"),
		Layer: plc.panel.LayerByName(plc.footerLayer),
	}
	plc.panel.Board.Plain.Texts = append(plc.panel.Board.Plain.Texts, footer)
}

type elementConfig struct {
	legendOffsetX, legendOffsetY   float64
	legendLocationFactor           float64
	legendAlign                    string
	legend                         string
	ticks, ticksLabels             bool
	ticksStartAngle, ticksEndAngle float64
	ticksLength, ticksWidth        float64
	ticksCount                     int
	ticksLabelsTexts               []string
}

// extract all the per-element config into a nice structure. Later this should help
// with refactoring the currently-ugly elementOp() into a bunch of separate operations
func elementConfigFromElement(elem eagle.Element) (elementConfig, error) {
	var err error
	ec := elementConfig{
		legend:           eagle.AttributeString(elem, "PANEL_LEGEND", elem.Name),
		ticksLabelsTexts: strings.Split(eagle.AttributeString(elem, "PANEL_LEGEND_TICKS_LABELS_TEXTS", ""), ","),
	}
	if ec.legendOffsetX, err = eagle.AttributeFloat(elem, "PANEL_LEGEND_OFFSET_X", 0.0); err != nil {
		return ec, err
	}
	if ec.legendOffsetY, err = eagle.AttributeFloat(elem, "PANEL_LEGEND_OFFSET_Y", 0.0); err != nil {
		return ec, err
	}
	legendlocation := eagle.AttributeString(elem, "PANEL_LEGEND_LOCATION", "above")
	switch legendlocation {
	case "above":
		ec.legendLocationFactor = 1
		ec.legendAlign = "bottom-center"
	case "below":
		ec.legendLocationFactor = -1
		ec.legendAlign = "top-center"
	default:
		return ec, fmt.Errorf("invalid value %q for attribute PANEL_LEGEND_LOCATION on object %q: must be 'above' or 'below'", legendlocation, elem.Name)
	}
	if ec.ticks, err = eagle.AttributeBool(elem, "PANEL_LEGEND_TICKS", false); err != nil {
		return ec, err
	}
	if ec.ticksLabels, err = eagle.AttributeBool(elem, "PANEL_LEGEND_TICKS_LABELS", false); err != nil {
		return ec, err
	}
	if ec.ticksLength, err = eagle.AttributeFloat(elem, "PANEL_LEGEND_TICKS_LENGTH", 1.5); err != nil {
		return ec, err
	}
	if ec.ticksWidth, err = eagle.AttributeFloat(elem, "PANEL_LEGEND_TICKS_WIDTH", 0.25); err != nil {
		return ec, err
	}
	// default values for start and end angles suit a typical single-turn potentiometer
	// with a 300-degree rotation, like Alpha 9mm vertical pots
	// https://www.thonk.co.uk/documents/alpha/9mm/Alpha%209mm%20Vertical%20-%20Linear%20Taper%20B1K-B500K.pdf
	if ec.ticksStartAngle, err = eagle.AttributeFloat(elem, "PANEL_LEGEND_TICKS_START_ANGLE", -60.0); err != nil {
		return ec, err
	}
	if ec.ticksEndAngle, err = eagle.AttributeFloat(elem, "PANEL_LEGEND_TICKS_END_ANGLE", 240.0); err != nil {
		return ec, err
	}
	// everything should go up to (at least) 11
	if ec.ticksCount, err = eagle.AttributeInt(elem, "PANEL_LEGEND_TICKS_COUNT", 11); err != nil {
		return ec, err
	}
	if ec.ticksLabels && len(ec.ticksLabelsTexts) != ec.ticksCount {
		return ec, fmt.Errorf("incorrect number of tick labels provided for object %q: ticks = %v, labels = %v", elem.Name, ec.ticksCount, len(ec.ticksLabelsTexts))
	}
	log.Printf("element config for %s: %+v", elem.Name, ec)
	return ec, nil
}

func elementOp(plc panelLayoutContext, elem eagle.Element) {
	hole, needHole, err := holeForPanelElement(elem)
	if err != nil {
		log.Fatalf("can't find drill size for element %q: %v", elem.Name, err)
	}
	if !needHole {
		return
	}
	// derive the per-element config
	elementConfig, err := elementConfigFromElement(elem)
	if err != nil {
		log.Fatalf("error extracting per-element config from attributes: %v", err)
	}
	// the hole was generated with coordinates from the source board, now
	// adjust them to be in the right place on the panel
	tstop := plc.panel.LayerByName("tStop")
	hole.X += plc.bc.XOffset
	hole.Y += plc.bc.YOffset
	plc.panel.Board.Plain.Holes = append(plc.panel.Board.Plain.Holes, hole)
	text := eagle.Text{
		X:     hole.X + elementConfig.legendOffsetX,
		Y:     hole.Y + ((elementConfig.legendOffsetY + (hole.Drill / 2.0) + *plc.cfg.TextSpacing) * elementConfig.legendLocationFactor),
		Size:  *plc.cfg.TextSize,
		Layer: plc.panel.LayerByName(plc.legendLayer),
		Text:  elementConfig.legend,
		Align: elementConfig.legendAlign,
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
		X: hole.X, Y: hole.Y,
		Radius: hole.Drill / 2.0,
		Width:  hsw,
		Layer:  tstop,
	}
	plc.panel.Board.Plain.Circles = append(plc.panel.Board.Plain.Circles, stop)
	if elementConfig.ticks {
		rpg := geometry.RadialPointGenerator{
			X: hole.X, Y: hole.Y,
			StartAngle: elementConfig.ticksStartAngle,
			EndAngle:   elementConfig.ticksEndAngle,
			Count:      elementConfig.ticksCount,
		}
		tickstarts := rpg.GenerateAtRadius(hole.Drill/2.0 + *plc.cfg.HoleStopRadius)
		tickends := rpg.GenerateAtRadius(hole.Drill/2.0 + *plc.cfg.HoleStopRadius + elementConfig.ticksLength)
		textorigins := rpg.GenerateAtRadius(hole.Drill/2.0 + *plc.cfg.HoleStopRadius + elementConfig.ticksLength + 2.0)
		for index, inner := range tickstarts {
			plc.panel.Board.Plain.Wires = append(plc.panel.Board.Plain.Wires, eagle.Wire{
				X1: inner.X, Y1: inner.Y,
				X2: tickends[index].X, Y2: tickends[index].Y,
				Width: elementConfig.ticksWidth,
				Layer: tstop,
			})
			if elementConfig.ticksLabels {
				plc.panel.Board.Plain.Texts = append(plc.panel.Board.Plain.Texts, eagle.Text{
					X: textorigins[index].X, Y: textorigins[index].Y,
					Align: "center",
					Size:  1.5,
					Text:  strings.TrimSpace(elementConfig.ticksLabelsTexts[index]),
					Layer: tstop,
				})
			}
		}
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
	SpecFile       *string
}

func configureFromFlags() config {
	formatList := "(" + strings.Join([]string{FormatEurorack, FormatPulplogic, FormatIntellijel, FormatSpec}, ",") + ")"
	cfg := config{
		Format:         flag.String("format", FormatEurorack, "panel format to create "+formatList),
		TextSpacing:    flag.Float64("text-spacing", 3.5, "spacing between a hole and its related label"),
		TextSize:       flag.Float64("text-size", 2.25, "label text size"),
		HoleStopRadius: flag.Float64("hole-stop-radius", 2.0, "Radius to pull back soldermask around a hole"),
		SpecFile:       flag.String("spec-file", "", "filename to read YAML panel spec from"),
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
