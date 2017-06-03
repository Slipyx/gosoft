package main

import (
	"unsafe"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/sdl_image"
)

// internal components are BGRA
// compatible when directly copying from
// an (A)RGB sdl surface or texture
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

// use sdl_image to load a surface from file
// and convert to Bitmap
func NewBitmapFromFile( file string ) *Bitmap {
	tImg, _ := img.Load( file )
	// ensure pixel format will be compatible when
	// copying directly to Bitmap's component array
	texImg, _ := tImg.ConvertFormat( sdl.PIXELFORMAT_ARGB8888, 0 )
	tImg.Free()

	tex := NewBitmap( int(texImg.W), int(texImg.H) )
	copy( tex.Comp, texImg.Pixels() )
	texImg.Free()

	return tex
}

func (b *Bitmap) Clear( shade byte ) {
	for i := range b.Comp { b.Comp[i] = shade }
}

// DrawPixel is always r, g, b, a
func (bm *Bitmap) DrawPixel( x, y int, r, g, b, a byte ) {
	ix := (y * bm.Width + x) * 4
	bm.Comp[ix] = b
	bm.Comp[ix + 1] = g
	bm.Comp[ix + 2] = r
	bm.Comp[ix + 3] = a
}

func (bm *Bitmap) CopyPixel( dx, dy, sx, sy int, src *Bitmap ) {
	dix := (dy * bm.Width + dx) * 4
	six := (sy * src.Width + sx) * 4

	bm.Comp[dix] = src.Comp[six]
	bm.Comp[dix + 1] = src.Comp[six + 1]
	bm.Comp[dix + 2] = src.Comp[six + 2]
	bm.Comp[dix + 3] = src.Comp[six + 3]
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
