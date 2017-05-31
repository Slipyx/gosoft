package main

import ( "math" )

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
	scanBuf []int
}

// creates and initializes new 3d render context
// using the const VID dimensions
func CreateRenderContext() *RenderContext {
	ctx := &RenderContext{}
	ctx.Bm = NewBitmap( VID_W, VID_H )
	ctx.scanBuf = make( []int, VID_H * 2 )
	return ctx
}

func (r *RenderContext) DrawScanBuffer( y, xmin, xmax int ) {
	r.scanBuf[y * 2] = xmin
	r.scanBuf[y * 2 + 1] = xmax
}

func (r *RenderContext) FillShape( ymin, ymax int ) {
	for j := ymin; j < ymax; j++ {
		xmin := r.scanBuf[j * 2]
		xmax := r.scanBuf[j * 2 + 1]

		for i := xmin; i < xmax; i++ {
			r.Bm.DrawPixel( i, j, 0xff, 0xff, 0xff, 0xff )
		}
	}
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

	hand := 1
	if minY.TriangleArea2( maxY, midY ) < 0 { hand = 0 }

	r.ScanConvertTriangle( minY, midY, maxY, hand )
	r.FillShape( int(math.Ceil( float64(minY.Pos.Y) )), int(math.Ceil( float64(maxY.Pos.Y) )) )
}

func (r *RenderContext) ScanConvertTriangle( minY, midY, maxY Vertex, hand int ) {
	r.ScanConvertLine( minY, maxY, 0 + hand )
	r.ScanConvertLine( minY, midY, 1 - hand )
	r.ScanConvertLine( midY, maxY, 1 - hand )
}

func (r *RenderContext) ScanConvertLine( minY, maxY Vertex, side int ) {
	startY, endY := int(math.Ceil( float64(minY.Pos.Y) )), int(math.Ceil( float64(maxY.Pos.Y) ))
	//startX, endX := int(math.Ceil( float64(minY.Pos.X) )), int(math.Ceil( float64(maxY.Pos.X) ))

	distY, distX := maxY.Pos.Y - minY.Pos.Y, maxY.Pos.X - minY.Pos.X

	if distY <= 0 { return }

	xStep := float32(distX) / float32(distY)
	yPreStep := float32(startY) - minY.Pos.Y
	curX := minY.Pos.X + yPreStep * xStep

	for j := startY; j < endY; j++ {
		r.scanBuf[j * 2 + side] = int(math.Ceil( float64(curX) ))
		curX += xStep
	}
}
