package instructions

import (
	"github/com/lucasbn/sqlite-clone/app/btree"
	"github/com/lucasbn/sqlite-clone/app/machine/state"
)

type Next struct {
	Cursor      uint32
	FromAddress uint32
}

var _ Instruction = Next{}

func (next Next) Execute(s *state.MachineState, p *btree.BTreeProcessor) [][]int {
	didAdvanced := p.GetCursor(next.Cursor).Next()

	if didAdvanced {
		s.CurrentAddress = int(next.FromAddress)
	} else {
		s.CurrentAddress++
	}

	return [][]int{}
}
