package utils

import (
	"math"
)

// type HashFunc func([]byte) uint32
var hashFuncs = []func([]byte) uint32{
	func(x []byte) uint32 {
		return 1
	},
	func(x []byte) uint32 {

		return 1
	},
	func(x []byte) uint32 {

		return 1
	},
	func(x []byte) uint32 {

		return 1
	},
}

type BloomFilter struct {
	// *Bitmap
	m     int
	k     int
	array []byte
	cnt   int
	seed  uint64
}

func NewBloomFilter(memSize int, total_elems int, seed uint64) *BloomFilter {

	// var k int = int(float64(size/uint64(elemCnt)) * math.Log(2))

	k := len(hashFuncs)
	// size = int(math.Max(float64(size), float64(1024 * 1024)))
	return &BloomFilter{
		m:     memSize,
		k:     k,
		array: make([]byte, memSize),
	}

}

func (bf *BloomFilter) Set(data []byte) {

	for _, hf := range hashFuncs {
		i := hf(data)
		idx := i / 8
		pos := i % 8
		bf.array[idx] |= 1 << pos
	}

	return
}

func (bf *BloomFilter) Test(data []byte) bool {

	for _, hf := range hashFuncs {
		i := hf(data)
		idx := i / 8
		pos := i % 8
		if (bf.array[idx] & (1 << pos)) == 0 {
			return false
		}
	}

	return true
}

func (bf *BloomFilter) ErrRate() float64 {

	m := float64(bf.m)
	n := float64(bf.cnt)
	k := float64(bf.k)

	return errRate(m, n, k)
}

func errRate(m, n, k float64) float64 {

	return 1 - math.Pow(1-1/m, k*n)
}

func init() {

}
