package testdata

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

var (
	db *sqlx.DB
)

func f2(name string) string {
	sql1 := fmt.Sprintf("select id from tb where name = %s and 2 > 1", name)
	db.Exec(sql1)
	return "sql1"
}

// func f1() {
// 	db := &sqlx.DB{}
// 	name := "jake"
// 	sql1 := fmt.Sprintf("select id from tb where name = %s and 2 > 1", name)
// 	sql2 := "select id from tb where name = " + name + "and 2 > 1"
// 	sql3 := "select id from tb where name = " + f2(name) + "and 2 > 1"
// 	db.Exec(sql1)
// 	db.Exec(sql2)
// 	db.Exec(sql3)
// 	sql4 := sql3
// 	_ = sql4
// }
