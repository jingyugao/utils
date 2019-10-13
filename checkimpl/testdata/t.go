package testdata

type t1 struct{}

type t2 struct{}

type t3 t1

type t4 struct {
	t1
}

type t5 struct {
	*t1
}

type t6 t2

type t7 struct {
	t2
}

type t8 struct {
	*t2
}

func (t1) f1()
