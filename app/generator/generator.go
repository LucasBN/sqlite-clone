package generator

import (
	"github/com/lucasbn/sqlite-clone/app/machine/instructions"

	"github.com/xwb1989/sqlparser"
)

func Generate(_ sqlparser.Statement) []instructions.Instruction {
	return []instructions.Instruction{
		instructions.OpenRead{RootPage: 5, CursorID: 0},
		instructions.Column{Cursor: 0, Column: 0, Register: 1},
		instructions.Column{Cursor: 0, Column: 1, Register: 2},
		instructions.Column{Cursor: 0, Column: 2, Register: 3},
		instructions.ResultRow{FromRegister: 1, ToRegister: 3},
		instructions.Next{Cursor: 0, FromAddress: 1},
		instructions.Halt{},
	}
}
