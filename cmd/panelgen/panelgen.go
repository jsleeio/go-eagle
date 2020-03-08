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
	"strings"

	"github.com/jsleeio/go-eagle/pkg/eagle"
	"github.com/jsleeio/go-eagle/pkg/format/eurorack"
	"github.com/jsleeio/go-eagle/pkg/format/intellijel"
	"github.com/jsleeio/go-eagle/pkg/format/pulplogic"
	filespec "github.com/jsleeio/go-eagle/pkg/format/spec"
	"github.com/jsleeio/go-eagle/pkg/panel"

	"github.com/jsleeio/go-eagle/internal/boardops/standard"
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

type config struct {
	Width        *int
	Format       *string
	Output       *string
	RefBoard     *string
	OutlineLayer *string
	SpecFile     *string
}

func configureFromFlags() (*config, error) {
	formatList := "(" + strings.Join([]string{FormatEurorack, FormatPulplogic, FormatIntellijel, FormatSpec}, ",") + ")"
	c := &config{
		Width:        flag.Int("width", 4, "width of the panel, in integer units appropriate for the format"),
		Format:       flag.String("format", FormatEurorack, "panel format to create "+formatList),
		RefBoard:     flag.String("reference-board", "", "reference Eagle board file to read layer information from"),
		Output:       flag.String("output", "newpanel.brd", "filename to write new Eagle board file to"),
		OutlineLayer: flag.String("outline-layer", "Dimension", "layer to draw board outline in"),
		SpecFile:     flag.String("spec-file", "", "filename to read YAML panel spec from"),
	}
	flag.Parse()
	if *c.RefBoard == "" {
		return nil, fmt.Errorf("a reference board file (-reference-board option) is required to acquire a list of Eagle layers")
	}
	return c, nil
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
	if err := standard.ApplyStandardBoardOperations(panel, spec); err != nil {
		return fmt.Errorf("error creating panel features: %v", err)
	}
	if err := panel.WriteFile(*cfg.Output); err != nil {
		return fmt.Errorf("can't write output board: %v", err)
	}
	return nil
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
	case FormatSpec:
		spec, err = filespec.LoadSpec(*cfg.SpecFile)
		if err != nil {
			fmt.Printf("error loading YAML panel spec from '%v': %v", *cfg.SpecFile, err)
			os.Exit(1)
		}
	default:
		fmt.Printf("unsupported format: %s\n", *cfg.Format)
		os.Exit(3)
	}
	if err := generatePanelBoardFile(cfg, spec); err != nil {
		fmt.Printf("error generating panel: %v\n", err)
		os.Exit(2)
	}
}
