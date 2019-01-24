package utils

type Cache interface {
	Insert(key string, val interface{}, deleter func(string, interface{})) (err error)
	Lookup(key string) (val interface{}, err error)
}
