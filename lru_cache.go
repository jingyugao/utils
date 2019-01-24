package utils

import (
	"sync"
)

type lru_node struct {
	next *lru_node
	prev *lru_node
	key  string
}

type LRU struct {
	items map[string]*lru_node
	head  *lru_node
	tail  *lru_node
	size  int
	mu    sync.Mutex
}

func NewLRU(size int) *LRU {
	head := new(lru_node)
	tail := new(lru_node)
	head.next = tail
	head.prev = tail
	tail.next = head
	tail.prev = head
	return &LRU{
		items: map[string]*lru_node{},
		head:  head,
		tail:  tail,
		size:  size,
		mu:    sync.Mutex{},
	}
}
func (l *LRU) moveToTail(one *lru_node) {

	if one == l.head || one == l.tail {
		panic("move head or tail")
		return
	}
	one.prev.next = one.next
	one.next.prev = one.prev

	l.tail.prev.next = one
	one.prev = l.tail.prev

	one.next = l.tail
	l.tail.prev = one

}

func (l *LRU) put(key string) {
	if one, ok := l.items[key]; ok {
		l.moveToTail(one)
		return
	}
	if len(l.items) == l.size {
		one := l.head.next
		l.moveToTail(one)
		delete(l.items, one.key)
		one.key = key
		l.items[key] = one
	} else {
		one := new(lru_node)
		l.tail.prev.next = one
		one.prev = l.tail.prev

		one.next = l.tail
		l.tail.prev = one

		one.key = key
		l.items[key] = one
	}
}
