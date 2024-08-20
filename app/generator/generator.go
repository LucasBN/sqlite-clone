package generator

import (
	"github/com/lucasbn/sqlite-clone/app/machine/instructions"

	"github.com/xwb1989/sqlparser"
)

func Generate(_ sqlparser.Statement) []instructions.Instruction {
	return []instructions.Instruction{
		instructions.Integer{Register: 1, Value: 1},
		instructions.Integer{Register: 2, Value: 2},
		instructions.ResultRow{FromRegister: 1, ToRegister: 2},
		instructions.Halt{},
	}
}
