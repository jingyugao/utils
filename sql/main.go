package main

var (
	XX = "tb1"
)

func f(name string) string {
	sql1 := `select * from tb`
	if name == "xx" {
		name = "zz"
	} else {
		name = "yy"
	}
	return sql1
}
func main() {
	println(f("x"))
	println(f(XX))
}
