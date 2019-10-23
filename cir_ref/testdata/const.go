package testdata

import (
	_ "github.com/go-sql-driver/mysql"

	"github.com/jmoiron/sqlx"
)

var (
	db2 *sqlx.DB
)

func gg() {
	var f1 string
	db.Exec(f1)
}
func ff() {
	var str1 string
	var a int
	if a == 1 {
		str1 = "x"
	} else {
		str1 = "y"
	}
	db2.Exec(str1)
}
