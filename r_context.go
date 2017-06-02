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

func (r *RenderContext) DrawMesh( mesh *Mesh, transform Mat4, texture *Bitmap ) {
	for i := 0; i < len( mesh.Indices ); i += 3 {
		r.FillTriangle( mesh.Vertices[mesh.Indices[i]].Transform( transform ),
			mesh.Vertices[mesh.Indices[i + 1]].Transform( transform ),
			mesh.Vertices[mesh.Indices[i + 2]].Transform( transform ), texture )
	}
}

func (r *RenderContext) FillTriangle( v1, v2, v3 Vertex, texture *Bitmap ) {
	var sstf Mat4;
	sstf.InitScreenSpaceTransform( float32(r.Bm.Width) / 2.0, float32(r.Bm.Height) / 2.0 )

	minY := v1.Transform( sstf ).PerspectiveDivide()
	midY := v2.Transform( sstf ).PerspectiveDivide()
	maxY := v3.Transform( sstf ).PerspectiveDivide()

	if minY.TriangleArea2( maxY, midY ) >= 0 { return }

	if maxY.Pos.Y < midY.Pos.Y { maxY, midY = midY, maxY }
	if midY.Pos.Y < minY.Pos.Y { midY, minY = minY, midY }
	if maxY.Pos.Y < midY.Pos.Y { maxY, midY = midY, maxY }

	r.ScanTriangle( minY, midY, maxY, minY.TriangleArea2( maxY, midY ) >= 0, texture )
}

func (r *RenderContext) DrawScanLine( left, right Edge, j int, texture *Bitmap ) {
	xmin := int(math.Ceil( float64(left.x) ))
	xmax := int(math.Ceil( float64(right.x) ))

	xPreStep := float32(xmin) - left.x

	xDist := right.x - left.x

	texCoordXXStep := (right.texCoordX - left.texCoordX) / xDist
	texCoordYXStep := (right.texCoordY - left.texCoordY) / xDist
	oneOverZXStep := (right.oneOverZ - left.oneOverZ) / xDist

	texCoordX := left.texCoordX + texCoordXXStep * xPreStep
	texCoordY := left.texCoordY + texCoordYXStep * xPreStep
	oneOverZ := left.oneOverZ + oneOverZXStep * xPreStep

	//minCol := left.col.Add( grad.colXStep.Mul( xPreStep ) )
	//maxCol := right.col.Add( grad.colXStep.Mul( xPreStep ) )

	for i := xmin; i < xmax; i++ {
		z := 1.0 / oneOverZ
		srcX := int((texCoordX * z) * float32(texture.Width - 1) + 0.5)
		srcY := int((texCoordY * z) * float32(texture.Height - 1) + 0.5)

		r.Bm.CopyPixel( i, j, srcX, srcY, texture )

		texCoordX += texCoordXXStep
		texCoordY += texCoordYXStep
		oneOverZ += oneOverZXStep
	}
}

func (r *RenderContext) ScanTriangle( minY, midY, maxY Vertex, hand bool, texture *Bitmap ) {
	grad := NewGradients( minY, midY, maxY )

	topToBot := NewEdge( grad, minY, maxY, 0 )
	topToMid := NewEdge( grad, minY, midY, 0 )
	midToBot := NewEdge( grad, midY, maxY, 1 )

	// scan edges apparently needs to modify
	// the passed in edge so it continues correctly
	// on subsequent calls
	r.ScanEdges( &topToBot, &topToMid, hand, texture )
	r.ScanEdges( &topToBot, &midToBot, hand, texture )
}

// passed in edges need to be mutable so it can continue
// on subsequent calls
func (r *RenderContext) ScanEdges( a, b *Edge, hand bool, texture *Bitmap ) {
	//var left *Edge = a
	//var right *Edge = b

	yStart, yEnd := b.yStart, b.yEnd

	if hand { a, b = b, a }

	for j := yStart; j < yEnd; j++ {
		r.DrawScanLine( *a, *b, j, texture )
		a.Step()
		b.Step()
	}
}
