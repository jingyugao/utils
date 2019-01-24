package utils

import "sync"

const (
	kMaxHeight = 12
)

type Node struct {
	Key  string
	Val  interface{}
	Next *Node
	mu   sync.Mutex
}

type SkipList struct {
	head  *Node
	level int
}

func (n *Node) findLessThan(key string) *Node {
	for p := n; p != nil; p = p.Next {
		if p.Key >= key {
			return p
		}
	}
	return nil
}

func (sl *SkipList) insert(key string, val interface{}) {
	if sl.head == nil {
		sl.head = new(Node)
	}

}
