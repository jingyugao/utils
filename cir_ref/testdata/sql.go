package testdata

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

var (
	db *sqlx.DB
)

func f5(str string) string {
	return str
}

const (
	strc = "qwer"
)

func f2(name string) string {
	names := []string{}
	str0 := "sql"
	str1 := "q1" + "q2"
	str2 := str0 + str1
	str3 := fmt.Sprintf("%s %s", str2, str0)
	str4 := f5(str3) + strc
	str5 := str4 + strings.Join(names, ",")

	// str5 := "xx" + fmt.Sprintf("%s %s", "sql"+"q1"+name, fmt.Sprintf(name)) + strings.Join(names, ",")

	// sql1, xx := fmt.Sprintf("select id from tb where name = %s and 2 > 1", name), "qwe"
	db.Exec(str5)
	db.Exec(str4)
	db.Exec(str3)
	db.Exec(str2)
	db.Exec(str1)
	db.Exec(str0)
	// _ = xx

	return "xx"
}
