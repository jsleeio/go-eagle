package boardops

import (
	"github.com/jsleeio/go-eagle/pkg/eagle"
	"github.com/jsleeio/go-eagle/pkg/panel"
)

// BoardOperation functions do ... something ... to an Eagle board
type BoardOperation func(*eagle.Eagle, panel.Panel) error

// ApplyBoardOperations works through a list of board operation
// functions and applies them to an Eagle board object
func ApplyBoardOperations(board *eagle.Eagle, spec panel.Panel, ops []BoardOperation) error {
	for _, op := range ops {
		if err := op(board, spec); err != nil {
			return err
		}
	}
	return nil
}
