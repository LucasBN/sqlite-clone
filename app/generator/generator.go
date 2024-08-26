package generator

import (
	"github/com/lucasbn/sqlite-clone/app/machine/instructions"

	"github.com/xwb1989/sqlparser"
)

func Generate[T any](_ sqlparser.Statement) []instructions.Instruction[T] {
	return []instructions.Instruction[T]{
		instructions.OpenRead[T]{RootPage: 5, CursorID: 0},
		instructions.Rewind[T]{Cursor: 0},
		instructions.Column[T]{Cursor: 0, Column: 0, Register: 1},
		instructions.Column[T]{Cursor: 0, Column: 1, Register: 2},
		instructions.Column[T]{Cursor: 0, Column: 2, Register: 3},
		instructions.ResultRow[T]{FromRegister: 1, ToRegister: 3},
		instructions.Next[T]{Cursor: 0, FromAddress: 2},
		instructions.Halt[T]{},
	}
}
