package utils

type Bitmap struct {
	Array []byte
	Len   uint32
}

func NewBitmap(max uint32) *Bitmap {

	len := max/8 + 1
	return &Bitmap{
		Len:   len,
		Array: make([]byte, len),
	}
}

func (bitmap *Bitmap) Set(i uint32) {

	idx := i / 8
	pos := i % 8
	bitmap.Array[idx] |= 1 << pos
	return
}

func (bitmap *Bitmap) UnSet(i uint32) {

	idx := i / 8
	pos := i % 8
	bitmap.Array[idx] &= ^(1 << pos)
	return
}

func (bitmap *Bitmap) Test(i uint32) bool {

	idx := i / 8
	pos := i % 8
	return (bitmap.Array[idx] & (1 << pos)) != 0
}
