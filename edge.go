package main

import ( "math" )

type Edge struct {
	x, xStep float32
	yStart, yEnd int
	col, colStep Vec4
}

func NewEdge( grad Gradients, minY, maxY Vertex, minYIndex int ) Edge {
	ne := Edge{}

	ne.yStart = int(math.Ceil( float64(minY.Pos.Y) ))
	ne.yEnd = int(math.Ceil( float64(maxY.Pos.Y) ))

	yDist := maxY.Pos.Y - minY.Pos.Y
	xDist := maxY.Pos.X - minY.Pos.X

	ne.xStep = float32(xDist) / float32(yDist)
	yPreStep := float32(ne.yStart) - minY.Pos.Y
	ne.x = minY.Pos.X + yPreStep * ne.xStep

	xPreStep := ne.x - minY.Pos.X

	ne.col = grad.col[minYIndex].Add(
		grad.colYStep.Mul( yPreStep ) ).Add(
		grad.colXStep.Mul( xPreStep ) )

	ne.colStep = grad.colYStep.Add( grad.colXStep.Mul( ne.xStep ) )

	return ne
}

func (e *Edge) Step() {
	e.x += e.xStep
	e.col = e.col.Add( e.colStep )
}

