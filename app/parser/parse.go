package parser

import (
	"github.com/xwb1989/sqlparser"
)

func MustParse(sql string) sqlparser.Statement {
	stmt, err := sqlparser.Parse(sql)
	if err != nil {
		panic(err)
	}
	return stmt
}
