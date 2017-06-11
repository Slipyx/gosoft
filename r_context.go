package main

import (
	"math"
	"fmt"
)

const (
	// 3d rt
	VID_W = 320
	VID_H = 224

	// 2d rt
	VID2D_W = 640
	VID2D_H = 360
)

// render context
type RenderContext struct {
	Bm *Bitmap
	zBuffer []float32
	sstf Mat4
}

// creates and initializes new 3d render context
// using the const VID dimensions
func CreateRenderContext() *RenderContext {
	fmt.Println( "new ctx bro" )
	ctx := &RenderContext{}
	ctx.Bm = NewBitmap( VID_W, VID_H )

	ctx.zBuffer = make( []float32, ctx.Bm.Width * ctx.Bm.Height )

	ctx.sstf.InitScreenSpaceTransform(
		float32(ctx.Bm.Width) / 2.0, float32(ctx.Bm.Height) / 2.0 )

	return ctx
}

func (r *RenderContext) ClearDepthBuffer() {
	for i := range r.zBuffer { r.zBuffer[i] = math.MaxFloat32 }
}

func (r *RenderContext) DrawTriangle( v1, v2, v3 Vertex, texture *Bitmap ) {
	v1In := v1.IsInViewFrustum()
	v2In := v2.IsInViewFrustum()
	v3In := v3.IsInViewFrustum()

	if v1In && v2In && v3In {
		r.FillTriangle( v1, v2, v3, texture )
		return
	}

	// skip only if entire mesh is outside
	//if !v1In && !v2In && !v3In { return }

	vertices := make( []Vertex, 0 )
	auxList := make( []Vertex, 0 )

	vertices = append( vertices, v1, v2, v3 )

	if ClipPolygonAxis( &vertices, &auxList, 0 ) &&
		ClipPolygonAxis( &vertices, &auxList, 1 ) &&
		ClipPolygonAxis( &vertices, &auxList, 2 ) {
		iVert := vertices[0]

		for i := 1; i < len( vertices ) - 1; i++ {
			r.FillTriangle( iVert, vertices[i], vertices[i + 1], texture )
		}
	}
}

func ClipPolygonAxis( vertices, auxList *[]Vertex, compIndex int ) bool {
	ClipPolygonComponent( vertices, compIndex, 1, auxList )
	*vertices = nil

	if len( *auxList ) == 0 { return false }

	ClipPolygonComponent( auxList, compIndex, -1, vertices )
	*auxList = nil

	return !(len( *vertices ) == 0)
}

func ClipPolygonComponent( vertices *[]Vertex, compIndex int, compFactor float32, result *[]Vertex ) {
	pVertex := (*vertices)[len( *vertices ) - 1]
	pComp := pVertex.GetPosI( compIndex ) * compFactor
	pInside := pComp <= pVertex.Pos.W

	for _, v := range *vertices {
		cComp := v.GetPosI( compIndex ) * compFactor
		cInside := cComp <= v.Pos.W

		if cInside != pInside {
			lerpAmt := (pVertex.Pos.W - pComp) / ((pVertex.Pos.W - pComp) - (v.Pos.W - cComp))
			*result = append( *result, pVertex.Lerp( v, lerpAmt ) )
		}

		if cInside {
			*result = append( *result, v )
		}

		pVertex = v
		pComp = cComp
		pInside = cInside
	}
}

func (r *RenderContext) FillTriangle( v1, v2, v3 Vertex, texture *Bitmap ) {
	//var sstf Mat4;
	//sstf.InitScreenSpaceTransform( float32(r.Bm.Width) / 2.0, float32(r.Bm.Height) / 2.0 )

	minY := v1.Transform( r.sstf ).PerspectiveDivide()
	midY := v2.Transform( r.sstf ).PerspectiveDivide()
	maxY := v3.Transform( r.sstf ).PerspectiveDivide()

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
	depthXStep := (right.depth - left.depth) / xDist

	texCoordX := left.texCoordX + texCoordXXStep * xPreStep
	texCoordY := left.texCoordY + texCoordYXStep * xPreStep
	oneOverZ := left.oneOverZ + oneOverZXStep * xPreStep
	depth := left.depth + depthXStep * xPreStep

	//minCol := left.col.Add( grad.colXStep.Mul( xPreStep ) )
	//maxCol := right.col.Add( grad.colXStep.Mul( xPreStep ) )

	for i := xmin; i < xmax; i++ {
		zindex := j * r.Bm.Width + i

		if depth < r.zBuffer[zindex] {
			r.zBuffer[zindex] = depth

			z := 1.0 / oneOverZ
			srcX := int((texCoordX * z) * float32(texture.Width - 1) + 0.5)
			srcY := int((texCoordY * z) * float32(texture.Height - 1) + 0.5)

			// texture repeat
			if srcX >= texture.Width || srcX < 0 {
				srcX -= (srcX / texture.Width * texture.Width)
			}
			if srcY >= texture.Height || srcY < 0 {
				srcY -= (srcY / texture.Height * texture.Height)
			}

			r.Bm.CopyPixel( i, j, srcX, srcY, texture )
		}

		texCoordX += texCoordXXStep
		texCoordY += texCoordYXStep
		oneOverZ += oneOverZXStep
		depth += depthXStep
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
	var left *Edge = a
	var right *Edge = b

	//yStart, yEnd := b.yStart, b.yEnd

	if hand { left, right = right, left }

	for j := b.yStart; j < b.yEnd; j++ {
		r.DrawScanLine( *left, *right, j, texture )
		left.Step()
		right.Step()
	}
}
