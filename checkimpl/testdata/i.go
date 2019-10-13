package testdata

type i1 interface {
	f1()
}

type i2 i1

type i3 interface {
	i1
}

type i4 interface {
	i2
}

type i5 interface {
	f2()
}
