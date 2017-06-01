package main

import (
	"math"
	"fmt"
)

const (
	// 3d rt
	VID_W = 320
	VID_H = 200

	// 2d rt
	VID2D_W = 640
	VID2D_H = 360
)

// render context
type RenderContext struct {
	Bm *Bitmap
}

// creates and initializes new 3d render context
// using the const VID dimensions
func CreateRenderContext() *RenderContext {
	fmt.Println( "new ctx bro" )
	ctx := &RenderContext{}
	ctx.Bm = NewBitmap( VID_W, VID_H )
	return ctx
}

func (r *RenderContext) FillTriangle( v1, v2, v3 Vertex ) {
	var sstf Mat4;
	sstf.InitScreenSpaceTransform( float32(r.Bm.Width) / 2.0, float32(r.Bm.Height) / 2.0 )

	minY := v1.Transform( sstf ).PerspectiveDivide()
	midY := v2.Transform( sstf ).PerspectiveDivide()
	maxY := v3.Transform( sstf ).PerspectiveDivide()

	if maxY.Pos.Y < midY.Pos.Y { maxY, midY = midY, maxY }
	if midY.Pos.Y < minY.Pos.Y { midY, minY = minY, midY }
	if maxY.Pos.Y < midY.Pos.Y { maxY, midY = midY, maxY }

	r.ScanTriangle( minY, midY, maxY, minY.TriangleArea2( maxY, midY ) >= 0 )
}

func (r *RenderContext) DrawScanLine( left, right Edge, j int ) {
	xmin := int(math.Ceil( float64(left.x) ))
	xmax := int(math.Ceil( float64(right.x) ))

	for i := xmin; i < xmax; i++ {
		r.Bm.DrawPixel( i, j, 0xff, 0xff, 0xff, 0xff )
	}
	//r.Bm.DrawPixel( xmin, j, 0xff, 0x00, 0x00, 0xff )
	//r.Bm.DrawPixel( xmax, j, 0x00, 0x00, 0xff, 0xff )
}

func (r *RenderContext) ScanTriangle( minY, midY, maxY Vertex, hand bool ) {
	topToBot := NewEdge( minY, maxY )
	topToMid := NewEdge( minY, midY )
	midToBot := NewEdge( midY, maxY )

	// scan edges apparently needs to modify
	// the passed in edge so it continues correctly
	// on subsequent calls
	r.ScanEdges( &topToBot, &topToMid, hand )
	r.ScanEdges( &topToBot, &midToBot, hand )
}

// passed in edges need to be mutable so it can continue
// on subsequent calls
func (r *RenderContext) ScanEdges( a, b *Edge, hand bool ) {
	//var left *Edge = a
	//var right *Edge = b

	yStart, yEnd := b.yStart, b.yEnd

	if hand { a, b = b, a }

	for j := yStart; j < yEnd; j++ {
		r.DrawScanLine( *a, *b, j )
		a.Step()
		b.Step()
	}
}
