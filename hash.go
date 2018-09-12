// Based on the MurmurHash2.cpp source from SMHasher & MurmurHash,
// https://code.google.com/p/smhasher/

package utils

import (
	"hash"
)

// Mixing constants; generated offline.
const (
	M    = 0x5bd1e995
	BIGM = 0xc6a4a7935bd1e995
	R    = 24
	BIGR = 47
)

// 32-bit mixing function.
func mix2(h uint32, k uint32) (uint32, uint32) {
	k *= M
	k ^= k >> R
	k *= M
	h *= M
	h ^= k
	return h, k
}

// The original MurmurHash2 32-bit algorithm by Austin Appleby.
func MurmurHash2(data []byte, seed uint32) (h uint32) {
	var k uint32

	// Initialize the hash to a 'random' value
	h = seed ^ uint32(len(data))

	// Mix 4 bytes at a time into the hash
	for l := len(data); l >= 4; l -= 4 {
		k = uint32(data[0]) | uint32(data[1])<<8 | uint32(data[2])<<16 | uint32(data[3])<<24
		h, k = mix2(h, k)
		data = data[4:]
	}

	// Handle the last few bytes of the input array
	switch len(data) {
	case 3:
		h ^= uint32(data[2]) << 16
		fallthrough
	case 2:
		h ^= uint32(data[1]) << 8
		fallthrough
	case 1:
		h ^= uint32(data[0])
		h *= M
	}

	// Do a few final mixes of the hash to ensure the last few bytes are well incorporated
	h ^= h >> 13
	h *= M
	h ^= h >> 15

	return
}

// -----------------------------------------------------------------------------

// MurmurHash64A (64-bit) algorithm by Austin Appleby.
func MurmurHash64A(data []byte, seed uint64) (h uint64) {
	var k uint64

	h = seed ^ uint64(len(data))*BIGM

	for l := len(data); l >= 8; l -= 8 {
		k = uint64(data[0]) | uint64(data[1])<<8 | uint64(data[2])<<16 | uint64(data[3])<<24 |
			uint64(data[4])<<32 | uint64(data[5])<<40 | uint64(data[6])<<48 | uint64(data[7])<<56

		k *= BIGM
		k ^= k >> BIGR
		k *= BIGM

		h ^= k
		h *= BIGM

		data = data[8:]
	}

	switch len(data) {
	case 7:
		h ^= uint64(data[6]) << 48
		fallthrough
	case 6:
		h ^= uint64(data[5]) << 40
		fallthrough
	case 5:
		h ^= uint64(data[4]) << 32
		fallthrough
	case 4:
		h ^= uint64(data[3]) << 24
		fallthrough
	case 3:
		h ^= uint64(data[2]) << 16
		fallthrough
	case 2:
		h ^= uint64(data[1]) << 8
		fallthrough
	case 1:
		h ^= uint64(data[0])
		h *= BIGM
	}

	h ^= h >> BIGR
	h *= BIGM
	h ^= h >> BIGR

	return
}

// -----------------------------------------------------------------------------

// MurmurHash2A (32-bit) algorithm by Austin Appleby.
func MurmurHash2A(data []byte, seed uint32) (h uint32) {
	var k, t, ln uint32

	ln = uint32(len(data))
	h = seed

	for l := len(data); l >= 4; l -= 4 {
		k = uint32(data[0]) | uint32(data[1])<<8 | uint32(data[2])<<16 | uint32(data[3])<<24
		h, k = mix2(h, k)
		data = data[4:]
	}

	switch len(data) {
	case 3:
		t ^= uint32(data[2]) << 16
		fallthrough
	case 2:
		t ^= uint32(data[1]) << 8
		fallthrough
	case 1:
		t ^= uint32(data[0])
	}

	h, _ = mix2(h, t)
	h, _ = mix2(h, ln)

	h ^= h >> 13
	h *= M
	h ^= h >> 15

	return
}

// -----------------------------------------------------------------------------
// Based on the implementation of CMurmurHash2A by Austin Appleby.
// Designed to work incrementally.

type (
	murmur32 struct {
		seed  uint32
		hash  uint32
		tail  uint32
		count uint32
		size  uint32
	}
)

// New32 returns a new 32-bit MurmurHash2.
func New32(seed uint32) hash.Hash32 {
	return &murmur32{seed, seed, 0, 0, 0}
}

func (m *murmur32) mixTail(data []byte) []byte {
	// Mix in a slice of less than 4 bytes. Or, if we've previously mixed in some
	// trailing bytes, add more until we've mixed in 4 bytes.
	for l := len(data); l > 0 && (l < 4 || m.count > 0); l-- {
		m.tail |= uint32(data[0]) << (m.count * 8)
		m.count++
		if m.count == 4 {
			m.hash, _ = mix2(m.hash, m.tail)
			m.tail = 0
			m.count = 0
		}
		data = data[1:]
	}
	return data
}

// Reset the hash to its initial state.
func (m *murmur32) Reset() {
	m.hash = m.seed
	m.tail = 0
	m.count = 0
	m.size = 0
}

// Add some data to the running hash.
func (m *murmur32) Write(data []byte) (n int, err error) {
	n = len(data)
	m.size += uint32(n)

	data = m.mixTail(data)

	for l := len(data); l >= 4; l -= 4 {
		k := uint32(data[0]) | uint32(data[1])<<8 | uint32(data[2])<<16 | uint32(data[3])<<24
		m.hash, k = mix2(m.hash, k)
		data = data[4:]
	}

	m.mixTail(data) // data should be length 0 after this
	return
}

// Get the hash result.
func (m *murmur32) Sum32() (hash uint32) {
	hash, _ = mix2(m.hash, m.tail)
	hash, _ = mix2(hash, m.size)

	hash ^= hash >> 13
	hash *= M
	hash ^= hash >> 15

	return
}

func (m *murmur32) Size() int { return 4 }

func (m *murmur32) BlockSize() int { return 4 }

func (m *murmur32) Sum(in []byte) []byte {
	v := m.Sum32()
	return append(in, byte(v>>24), byte(v>>16), byte(v>>8), byte(v))
}
