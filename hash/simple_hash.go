package hash

// Additive Hash
// This takes 5n+3 instructions. There is no mixing step. The combining step handles one byte at a time. Input bytes commute. The table length must be prime, and can't be much bigger than one byte because the value of variable hash is never much bigger than one byte.

func additive(key []byte, prime uint32) uint32 {
	var hash, i uint32
	size := uint32(len(key))
	hash = size
	for i = 0; i < size; i++ {
		hash += uint32(key[i])
	}
	return (hash % prime)
}

// RotateHash takes 8n+3 instructions. This is the same as the additive hash, except it has a mixing step (a circular shift by 4) and the combining step is exclusive-or instead of addition. The table size is a prime, but the prime can be any size. On machines with a rotate (such as the Intel x86 line) this is 6n+2 instructions. I have seen the (hash % prime) replaced with
//   hash = (hash ^ (hash>>10) ^ (hash>>20)) & mask;
// eliminating the % and allowing the table size to be a power of 2, making this 6n+6 instructions. % can be very slow, for example it is 230 times slower than addition on a Sparc 20.
func RotateHash(key []byte, prime uint32) uint32 {
	var hash, i uint32
	size := uint32(len(key))
	hash = size
	for i = 0; i < size; i++ {
		hash = (hash << 4) ^ (hash >> 28) ^ uint32(key[i])
	}
	return (hash % prime)
}

// func one_at_a_time(key []byte) uint32 {
// 	var hash, i uint32
// 	hash = 0
// 	size := uint32(len(key))
// 	for i = 0; i < size; i++ {
// 		hash += uint32(key[i])
// 		hash += (hash << 10)
// 		hash ^= (hash >> 6)
// 	}
// 	hash += (hash << 3)
// 	hash ^= (hash >> 11)
// 	hash += (hash << 15)
// 	return (hash & mask)
// }
