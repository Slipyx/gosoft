package main

type Gradients struct {
	//col []Vec4
	//colXStep, colYStep Vec4
	texCoordX, texCoordY []float32
	texCoordXXStep, texCoordXYStep float32
	texCoordYXStep, texCoordYYStep float32
	// perspective
	oneOverZ []float32
	oneOverZXStep, oneOverZYStep float32
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

	// perspective
	ng.oneOverZ = make( []float32, 3 )

	ng.oneOverZ[0] = 1.0 / minY.Pos.W
	ng.oneOverZ[1] = 1.0 / midY.Pos.W
	ng.oneOverZ[2] = 1.0 / maxY.Pos.W

	ng.oneOverZXStep = CalcXStep( ng.oneOverZ, minY, midY, maxY, oneOverDX )
	ng.oneOverZYStep = CalcYStep( ng.oneOverZ, minY, midY, maxY, oneOverDY )

	// texcoord
	ng.texCoordX = make( []float32, 3 )
	ng.texCoordY = make( []float32, 3 )

	ng.texCoordX[0] = minY.TexCoord.X * ng.oneOverZ[0]
	ng.texCoordX[1] = midY.TexCoord.X * ng.oneOverZ[1]
	ng.texCoordX[2] = maxY.TexCoord.X * ng.oneOverZ[2]

	ng.texCoordY[0] = minY.TexCoord.Y * ng.oneOverZ[0]
	ng.texCoordY[1] = midY.TexCoord.Y * ng.oneOverZ[1]
	ng.texCoordY[2] = maxY.TexCoord.Y * ng.oneOverZ[2]

	ng.texCoordXXStep = CalcXStep( ng.texCoordX, minY, midY, maxY, oneOverDX )
	ng.texCoordXYStep = CalcYStep( ng.texCoordX, minY, midY, maxY, oneOverDY )
	ng.texCoordYXStep = CalcXStep( ng.texCoordY, minY, midY, maxY, oneOverDX )
	ng.texCoordYYStep = CalcYStep( ng.texCoordY, minY, midY, maxY, oneOverDY )

	/*ng.colXStep = (((ng.col[1].Sub( ng.col[2] )).Mul(
		(minY.Pos.Y - maxY.Pos.Y))).Sub( ((ng.col[0].Sub(
		ng.col[2] )).Mul( (midY.Pos.Y - maxY.Pos.Y) )))).Mul( oneOverDX )

	ng.colYStep = (((ng.col[1].Sub( ng.col[2] )).Mul(
		(minY.Pos.X - maxY.Pos.X))).Sub( ((ng.col[0].Sub(
		ng.col[2] )).Mul( (midY.Pos.X - maxY.Pos.X) )))).Mul( oneOverDY )*/

	return ng
}

func CalcXStep( values []float32, minY, midY, maxY Vertex, oneOverDX float32 ) float32 {
	return (((values[1] - values[2]) * (minY.Pos.Y - maxY.Pos.Y)) -
		((values[0] - values[2]) * (midY.Pos.Y - maxY.Pos.Y))) * oneOverDX
}

func CalcYStep( values []float32, minY, midY, maxY Vertex, oneOverDY float32 ) float32 {
	return (((values[1] - values[2]) * (minY.Pos.X - maxY.Pos.X)) -
		((values[0] - values[2]) * (midY.Pos.X - maxY.Pos.X))) * oneOverDY
}

