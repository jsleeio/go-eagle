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
	"os"

	"github.com/jsleeio/go-eagle/pkg/eagle"
	"github.com/jsleeio/go-eagle/pkg/format/eurorack"
	"github.com/jsleeio/go-eagle/pkg/format/intellijel"
	"github.com/jsleeio/go-eagle/pkg/format/pulplogic"
	"github.com/jsleeio/go-eagle/pkg/panel"
)

const (
	// FormatEurorack is the Doepfer-defined 3U specification. Not Eurocard!
	FormatEurorack = "eurorack"
	// FormatPulplogic is the PulpLogic-defined 1U specification
	FormatPulplogic = "pulplogic"
	// FormatIntellijel is the Intellijel-defined 1U specification
	FormatIntellijel = "intellijel"
)

type config struct {
	Width          *int
	Format         *string
	Output         *string
	RefBoard       *string
	OutlineLayer   *string
	CopperPullback *float64
}

func configureFromFlags() (*config, error) {
	c := &config{
		Width:          flag.Int("width", 4, "width of the panel, in integer units appropriate for the format"),
		Format:         flag.String("format", FormatEurorack, "panel format to create (eurorack, pulplogic, intellijel)"),
		RefBoard:       flag.String("reference-board", "", "reference Eagle board file to read layer information from"),
		Output:         flag.String("output", "newpanel.brd", "filename to write new Eagle board file to"),
		OutlineLayer:   flag.String("outline-layer", "Dimension", "layer to draw board outline in"),
		CopperPullback: flag.Float64("copper-pullback", 0.1, "Distance to pull back the copper pour from the panel edge"),
	}
	flag.Parse()
	if *c.RefBoard == "" {
		return nil, fmt.Errorf("a reference board file (-reference-board option) is required to acquire a list of Eagle layers")
	}
	return c, nil
}

func boardOutline(x1, y1, x2, y2 float64, layer int) []eagle.Wire {
	return []eagle.Wire{
		{X1: x1, Y1: y1, X2: x2, Y2: y1, Layer: layer, Width: 0}, // bottom
		{X1: x1, Y1: y2, X2: x2, Y2: y2, Layer: layer, Width: 0}, // top
		{X1: x1, Y1: y1, X2: x1, Y2: y2, Layer: layer, Width: 0}, // left
		{X1: x2, Y1: y1, X2: x2, Y2: y2, Layer: layer, Width: 0}, // right
	}
}

func generatePanelBoardFile(cfg *config, spec panel.Panel) error {
	// the user very likely already has an Eagle board file nearby, so use it to
	// acquire a list of layers --- avoids hardcoding them, lets users use their
	// own mix/subset of layers if desired
	ref, err := eagle.LoadEagleFile(*cfg.RefBoard)
	if err != nil {
		return fmt.Errorf("can't load reference board: %v", err)
	}
	panel := ref.CloneEmpty()
	if err := applyBoardOperations(cfg, panel, spec, standardPanelOperations()); err != nil {
		return fmt.Errorf("error creating panel features: %v", err)
	}
	if err := panel.WriteFile(*cfg.Output); err != nil {
		return fmt.Errorf("can't write output board: %v", err)
	}
	return nil
}

type boardOperation func(*config, *eagle.Eagle, panel.Panel) error

func applyBoardOperations(cfg *config, board *eagle.Eagle, spec panel.Panel, ops []boardOperation) error {
	for _, op := range ops {
		if err := op(cfg, board, spec); err != nil {
			return err
		}
	}
	return nil
}

func outlineWiresOp(cfg *config, board *eagle.Eagle, spec panel.Panel) error {
	adjust := spec.HorizontalFit() / 2 // half on left edge, half on right edge
	outline := boardOutline(
		0+adjust,
		0,
		spec.Width()-adjust,
		spec.Height(),
		board.LayerByName(*cfg.OutlineLayer),
	)
	for _, wire := range outline {
		board.Board.Plain.Wires = append(board.Board.Plain.Wires, wire)
	}
	return nil
}

func mountingHolesOp(cfg *config, board *eagle.Eagle, spec panel.Panel) error {
	for _, hole := range spec.MountingHoles() {
		board.Board.Plain.Holes = append(board.Board.Plain.Holes, eagle.Hole{
			X:     hole.X,
			Y:     hole.Y,
			Drill: spec.MountingHoleDiameter(),
		})
	}
	return nil
}

func copperFillOp(cfg *config, board *eagle.Eagle, spec panel.Panel) error {
	top := eagle.Rectangle{
		X1:    *cfg.CopperPullback,
		Y1:    *cfg.CopperPullback,
		X2:    spec.Width() - *cfg.CopperPullback,
		Y2:    spec.Height() - *cfg.CopperPullback,
		Layer: board.LayerByName("Top"),
	}
	bottom := eagle.Rectangle{
		X1:    *cfg.CopperPullback,
		Y1:    *cfg.CopperPullback,
		X2:    spec.Width() - *cfg.CopperPullback,
		Y2:    spec.Height() - *cfg.CopperPullback,
		Layer: board.LayerByName("Bottom"),
	}
	board.Board.Plain.Rectangles = append(board.Board.Plain.Rectangles, top)
	board.Board.Plain.Rectangles = append(board.Board.Plain.Rectangles, bottom)
	return nil
}

func standardPanelOperations() []boardOperation {
	return []boardOperation{
		outlineWiresOp,
		mountingHolesOp,
		copperFillOp,
	}
}

func main() {
	cfg, err := configureFromFlags()
	if err != nil {
		fmt.Printf("configuration error: %v\n", err)
		os.Exit(1)
	}
	var spec panel.Panel
	switch *cfg.Format {
	case FormatEurorack:
		spec = eurorack.NewEurorack(*cfg.Width)
	case FormatPulplogic:
		spec = pulplogic.NewPulplogic(*cfg.Width)
	case FormatIntellijel:
		spec = intellijel.NewIntellijel(*cfg.Width)
	default:
		fmt.Printf("unsupported format: %s\n", *cfg.Format)
		os.Exit(3)
	}
	if err := generatePanelBoardFile(cfg, spec); err != nil {
		fmt.Printf("error generating panel: %v\n", err)
		os.Exit(2)
	}
}
