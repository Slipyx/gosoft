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

func (r *RenderContext) DrawScanLine( grad Gradients, left, right Edge, j int ) {
	xmin := int(math.Ceil( float64(left.x) ))
	xmax := int(math.Ceil( float64(right.x) ))

	xPreStep := float32(xmin) - left.x

	minCol := left.col.Add( grad.colXStep.Mul( xPreStep ) )
	maxCol := right.col.Add( grad.colXStep.Mul( xPreStep ) )

	lerpAmt := float32(0.0)
	lerpStep := 1.0 / float32(xmax - xmin)

	for i := xmin; i < xmax; i++ {
		col := minCol.Lerp( maxCol, lerpAmt )

		r.Bm.DrawPixel( i, j, byte(col.X * 0xff),
			byte(col.Y * 0xff), byte(col.Z * 0xff), 0xff )
		lerpAmt += lerpStep
	}
	//r.Bm.DrawPixel( xmin, j, 0xff, 0x00, 0x00, 0xff )
	//r.Bm.DrawPixel( xmax, j, 0x00, 0x00, 0xff, 0xff )
}

func (r *RenderContext) ScanTriangle( minY, midY, maxY Vertex, hand bool ) {
	grad := NewGradients( minY, midY, maxY )

	topToBot := NewEdge( grad, minY, maxY, 0 )
	topToMid := NewEdge( grad, minY, midY, 0 )
	midToBot := NewEdge( grad, midY, maxY, 1 )

	// scan edges apparently needs to modify
	// the passed in edge so it continues correctly
	// on subsequent calls
	r.ScanEdges( grad, &topToBot, &topToMid, hand )
	r.ScanEdges( grad, &topToBot, &midToBot, hand )
}

// passed in edges need to be mutable so it can continue
// on subsequent calls
func (r *RenderContext) ScanEdges( grad Gradients, a, b *Edge, hand bool ) {
	//var left *Edge = a
	//var right *Edge = b

	yStart, yEnd := b.yStart, b.yEnd

	if hand { a, b = b, a }

	for j := yStart; j < yEnd; j++ {
		r.DrawScanLine( grad, *a, *b, j )
		a.Step()
		b.Step()
	}
}
