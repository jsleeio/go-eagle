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
