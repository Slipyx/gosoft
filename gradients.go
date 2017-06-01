package main

type Gradients struct {
	//col []Vec4
	//colXStep, colYStep Vec4
	texCoordX, texCoordY []float32
	texCoordXXStep, texCoordXYStep float32
	texCoordYXStep, texCoordYYStep float32
}

func NewGradients( minY, midY, maxY Vertex ) Gradients {
	ng := Gradients{}

	//ng.col = make( []Vec4, 3 )

	//ng.col[0] = minY.Col
	//ng.col[1] = midY.Col
	//ng.col[2] = maxY.Col

	var oneOverDX float32 = 1.0 /
		(((midY.Pos.X - maxY.Pos.X) * (minY.Pos.Y - maxY.Pos.Y)) -
		((minY.Pos.X - maxY.Pos.X) * (midY.Pos.Y - maxY.Pos.Y)))

	oneOverDY := -oneOverDX

	ng.texCoordX = make( []float32, 3 )
	ng.texCoordY = make( []float32, 3 )

	ng.texCoordX[0] = minY.TexCoord.X
	ng.texCoordX[1] = midY.TexCoord.X
	ng.texCoordX[2] = maxY.TexCoord.X

	ng.texCoordY[0] = minY.TexCoord.Y
	ng.texCoordY[1] = midY.TexCoord.Y
	ng.texCoordY[2] = maxY.TexCoord.Y

	ng.texCoordXXStep = (((ng.texCoordX[1] - ng.texCoordX[2]) *
		(minY.Pos.Y - maxY.Pos.Y)) -
		((ng.texCoordX[0] - ng.texCoordX[2]) *
		(midY.Pos.Y - maxY.Pos.Y))) * oneOverDX

	ng.texCoordXYStep = (((ng.texCoordX[1] - ng.texCoordX[2]) *
		(minY.Pos.X - maxY.Pos.X)) -
		((ng.texCoordX[0] - ng.texCoordX[2]) *
		(midY.Pos.X - maxY.Pos.X))) * oneOverDY

	ng.texCoordYXStep = (((ng.texCoordY[1] - ng.texCoordY[2]) *
		(minY.Pos.Y - maxY.Pos.Y)) -
		((ng.texCoordY[0] - ng.texCoordY[2]) *
		(midY.Pos.Y - maxY.Pos.Y))) * oneOverDX

	ng.texCoordYYStep = (((ng.texCoordY[1] - ng.texCoordY[2]) *
		(minY.Pos.X - maxY.Pos.X)) -
		((ng.texCoordY[0] - ng.texCoordY[2]) *
		(midY.Pos.X - maxY.Pos.X))) * oneOverDY

	/*ng.colXStep = (((ng.col[1].Sub( ng.col[2] )).Mul(
		(minY.Pos.Y - maxY.Pos.Y))).Sub( ((ng.col[0].Sub(
		ng.col[2] )).Mul( (midY.Pos.Y - maxY.Pos.Y) )))).Mul( oneOverDX )

	ng.colYStep = (((ng.col[1].Sub( ng.col[2] )).Mul(
		(minY.Pos.X - maxY.Pos.X))).Sub( ((ng.col[0].Sub(
		ng.col[2] )).Mul( (midY.Pos.X - maxY.Pos.X) )))).Mul( oneOverDY )*/

	return ng
}

