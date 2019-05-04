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

package standard

import (
	"github.com/jsleeio/go-eagle/internal/boardops"
	"github.com/jsleeio/go-eagle/internal/boardops/util"
	"github.com/jsleeio/go-eagle/pkg/eagle"
	"github.com/jsleeio/go-eagle/pkg/panel"
)

const (
	// CopperPullback indicates how far a copper pour will be "pulled
	// back" from the edge of a board. This shouldn't really need to
	// be configurable, so keep it simple and use a constant
	CopperPullback = 0.5
)

// ApplyStandardBoardOperations applies the minimal set of baseline
// board features to an Eagle board: an outline, some mounting holes
// and copper pours in the Top and Bottom copper layers.
func ApplyStandardBoardOperations(board *eagle.Eagle, spec panel.Panel) error {
	ops := []boardops.BoardOperation{
		outlineWiresOp,
		mountingHolesOp,
		copperFillOp,
		railKeepoutsOp,
	}
	return boardops.ApplyBoardOperations(board, spec, ops)
}

func outlineWiresOp(board *eagle.Eagle, spec panel.Panel) error {
	adjust := spec.HorizontalFit() / 2 // half on left edge, half on right edge
	outline := util.WireRectangle(
		0+adjust,
		0,
		spec.Width()-adjust,
		spec.Height(),
		board.LayerByName("Dimension"),
		0, // outline wires must be zero-width
	)
	for _, wire := range outline {
		board.Board.Plain.Wires = append(board.Board.Plain.Wires, wire)
	}
	return nil
}

func mountingHolesOp(board *eagle.Eagle, spec panel.Panel) error {
	for _, hole := range spec.MountingHoles() {
		board.Board.Plain.Holes = append(board.Board.Plain.Holes, eagle.Hole{
			X:     hole.X,
			Y:     hole.Y,
			Drill: spec.MountingHoleDiameter(),
		})
	}
	return nil
}

func railKeepoutsOp(board *eagle.Eagle, spec panel.Panel) error {
	layer := board.LayerByName("tKeepout")
	bRail := eagle.Rectangle{
		X1:    panel.LeftX(spec),
		Y1:    spec.MountingHoleBottomY(),
		X2:    panel.RightX(spec),
		Y2:    spec.MountingHoleBottomY() + spec.RailHeightFromMountingHole(),
		Layer: layer,
	}
	tRail := eagle.Rectangle{
		X1:    panel.LeftX(spec),
		Y1:    spec.MountingHoleTopY() - spec.RailHeightFromMountingHole(),
		X2:    panel.RightX(spec),
		Y2:    spec.MountingHoleTopY(),
		Layer: layer,
	}
	board.Board.Plain.Rectangles = append(board.Board.Plain.Rectangles, bRail)
	board.Board.Plain.Rectangles = append(board.Board.Plain.Rectangles, tRail)
	return nil
}

func copperFillOp(board *eagle.Eagle, spec panel.Panel) error {
	// normally the horizontal fit value shouldn't be used anywhere except the
	// panel outline, but if we don't include it here, the pullback will be
	// wrong, or possibly even completely ineffective. So, include it.
	adjust := spec.HorizontalFit() / 2
	top := eagle.Rectangle{
		X1:    CopperPullback + adjust,
		Y1:    CopperPullback,
		X2:    spec.Width() - (CopperPullback + adjust),
		Y2:    spec.Height() - CopperPullback,
		Layer: board.LayerByName("Top"),
	}
	bottom := eagle.Rectangle{
		X1:    CopperPullback + adjust,
		Y1:    CopperPullback,
		X2:    spec.Width() - (CopperPullback + adjust),
		Y2:    spec.Height() - CopperPullback,
		Layer: board.LayerByName("Bottom"),
	}
	board.Board.Plain.Rectangles = append(board.Board.Plain.Rectangles, top)
	board.Board.Plain.Rectangles = append(board.Board.Plain.Rectangles, bottom)
	return nil
}
