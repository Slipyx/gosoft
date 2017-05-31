package main

import (
	"unsafe"
)

// components are BGRA
// to be used for an RGB sdl texture.update
// endian pls
type Bitmap struct {
	Width int
	Height int
	Comp []byte
}

func NewBitmap( w, h int ) *Bitmap {
	nbm := &Bitmap{ w, h, nil }
	nbm.Comp = make( []byte, w * h * 4 )
	return nbm
}

func (self *Bitmap) Clear( shade byte ) {
	for i := range self.Comp { self.Comp[i] = shade }
}

// DrawPixel is always r, g, b, a
func (self *Bitmap) DrawPixel( x, y int, r, g, b, a byte ) {
	ix := (y * self.Width + x) * 4
	self.Comp[ix] = b
	self.Comp[ix + 1] = g
	self.Comp[ix + 2] = r
	self.Comp[ix + 3] = a
}

const ( BM_MAX_SZ = 512 )

// hopefully unused since we can pass the full Comp pointer to texture.update
func (self *Bitmap) CopyToPtr( dst unsafe.Pointer ) {
	for i := 0; i < self.Width * self.Height; i++ {
		a := self.Comp[i * 4]
		r := self.Comp[i * 4 + 1]
		g := self.Comp[i * 4 + 2]
		b := self.Comp[i * 4 + 3]
		(*[BM_MAX_SZ*BM_MAX_SZ]int)(dst)[i] =
			(int(a) << 24) | (int(r) << 16) | (int(g) << 8) | (int(b))
	}
	//var pptr unsafe.Pointer
	//var ppitch int
	//bmtex.Lock( nil, &pptr, &ppitch )
	//bmtest.CopyToPtr( pptr )
	//bmtex.Unlock()
}

