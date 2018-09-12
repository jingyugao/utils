package hash

func rot(x, k uint32) uint32 {
	return (((x) << (k)) | ((x) >> (32 - (k))))
}

// * mix -- mix 3 32-bit values reversibly.
func mix(a, b, c uint32) (uint32, uint32, uint32) {
	a -= c
	a ^= rot(c, 4)
	c += b
	b -= a
	b ^= rot(a, 6)
	a += c
	c -= b
	c ^= rot(b, 8)
	b += a
	a -= c
	a ^= rot(c, 16)
	c += b
	b -= a
	b ^= rot(a, 19)
	a += c
	c -= b
	c ^= rot(b, 4)
	b += a
	return a, b, c
}

func DobbHash(k []byte, initval uint32) uint32 {

	/* Set up the internal state */
	totalSize := len(k)
	size := len(k)

	var a uint32 = 0x9e3779b9
	var b uint32 = 0x9e3779b9 /* the golden ratio; an arbitrary value */
	c := initval              /* the previous hash value */

	/*---------------------------------------- handle most of the key */
	for size >= 12 {
		a += (uint32(k[0]) + uint32(k[1])<<8 + uint32(k[2])<<16 + uint32(k[3])<<24)
		b += (uint32(k[4]) + uint32(k[5])<<8 + uint32(k[6])<<16 + uint32(k[7])<<24)
		c += (uint32(k[8]) + uint32(k[9])<<8 + uint32(k[10])<<16 + uint32(k[11])<<24)
		a, b, c = mix(a, b, c)
		k = k[12:]
		size -= 12
	}

	/*------------------------------------- handle the last 11 bytes */
	c += uint32(totalSize)
	switch size { /* all the case statements fall through */
	case 11:
		c += uint32(k[10]) << 24
	case 10:
		c += uint32(k[9]) << 16
	case 9:
		c += uint32(k[8]) << 8
		/* the first byte of c is reserved for the  sizegth */
	case 8:
		b += uint32(k[7]) << 24
	case 7:
		b += uint32(k[6]) << 16
	case 6:
		b += uint32(k[5]) << 8
	case 5:
		b += uint32(k[4])
	case 4:
		a += uint32(k[3]) << 24
	case 3:
		a += uint32(k[2]) << 16
	case 2:
		a += uint32(k[1]) << 8
	case 1:
		a += uint32(k[0])
		/* case 0: nothing left to add */
	}
	a, b, c = mix(a, b, c)
	/*-------------------------------------------- report the result */
	return c
}
