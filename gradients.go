package main

type Gradients struct {
	col []Vec4
	colXStep, colYStep Vec4
}

func NewGradients( minY, midY, maxY Vertex ) Gradients {
	ng := Gradients{}

	ng.col = make( []Vec4, 3 )

	ng.col[0] = minY.Col
	ng.col[1] = midY.Col
	ng.col[2] = maxY.Col

	var oneOverDX float32 = 1.0 /
		(((midY.Pos.X - maxY.Pos.X) * (minY.Pos.Y - maxY.Pos.Y)) -
		((minY.Pos.X - maxY.Pos.X) * (midY.Pos.Y - maxY.Pos.Y)))

	oneOverDY := -oneOverDX

	ng.colXStep = (((ng.col[1].Sub( ng.col[2] )).Mul(
		(minY.Pos.Y - maxY.Pos.Y))).Sub( ((ng.col[0].Sub(
		ng.col[2] )).Mul( (midY.Pos.Y - maxY.Pos.Y) )))).Mul( oneOverDX )

	ng.colYStep = (((ng.col[1].Sub( ng.col[2] )).Mul(
		(minY.Pos.X - maxY.Pos.X))).Sub( ((ng.col[0].Sub(
		ng.col[2] )).Mul( (midY.Pos.X - maxY.Pos.X) )))).Mul( oneOverDY )

	return ng
}

