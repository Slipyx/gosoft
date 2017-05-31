package main

import (
	"math"
	"math/rand"
)

type Stars3D struct {
	spread float32
	speed float32

	starX []float32
	starY []float32
	starZ []float32
}

func NewStars3D( num int, spread, speed float32 ) *Stars3D {
	ns := &Stars3D{}

	ns.spread = spread
	ns.speed = speed

	ns.starX = make( []float32, num )
	ns.starY = make( []float32, num )
	ns.starZ = make( []float32, num )

	for i := 0; i < num; i++ {
		ns.InitStar( i )
	}

	return ns
}

func (self *Stars3D) InitStar( ix int ) {
	self.starX[ix] = (rand.Float32() - 0.5) * 2.0 * self.spread
	self.starY[ix] = (rand.Float32() - 0.5) * 2.0 * self.spread
	self.starZ[ix] = (rand.Float32() + 0.00001) * self.spread
}

func (self *Stars3D) UpdateAndRender( ctx *RenderContext, dt float32 ) {
	//target.Clear( 0x10 )

	tanHalfFOV := float32(math.Tan( math.Pi / 180.0 * 66.0 / 2.0 ))

	halfW := float32(ctx.Bm.Width / 2.0)
	halfH := float32(ctx.Bm.Height / 2.0)

	var tv [3]Vertex

	for i := 0; i < len(self.starX); i++ {
		self.starZ[i] -= dt * self.speed
		if self.starZ[i] <= 0 { self.InitStar( i ) }

		x := (self.starX[i] / (self.starZ[i] * tanHalfFOV)) * halfW + halfW
		y := (self.starY[i] / (self.starZ[i] * tanHalfFOV)) * halfH + halfH

		if x < 0 || x >= float32(ctx.Bm.Width) || y < 0 || y >= float32(ctx.Bm.Height) {
			self.InitStar( i )
		} else {
			tv[i].Pos = Vec4{x,y,0,1}
		}
	}

	ctx.FillTriangle( tv[0], tv[1], tv[2] )
}
