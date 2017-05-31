package main

// Vector
type Vec2 struct {
	X, Y float32
}

type Vec3 struct {
	X, Y, Z float32
}

type Vec4 struct {
	X, Y, Z, W float32
}

// Vertex
type Vertex struct {
	Pos Vec4
}

func (v Vertex) Transform( mat Mat4 ) Vertex {
	return Vertex{mat.Transform( v.Pos )}
}

func (v Vertex) PerspectiveDivide() Vertex {
	return Vertex{Vec4{ v.Pos.X / v.Pos.W, v.Pos.Y / v.Pos.W, v.Pos.Z / v.Pos.W, v.Pos.W }}
}

func (v Vertex) TriangleArea2( b, c Vertex ) float32 {
	x1 := b.Pos.X - v.Pos.X
	y1 := b.Pos.Y - v.Pos.Y
	x2 := c.Pos.X - v.Pos.X
	y2 := c.Pos.Y - v.Pos.Y

	return x1 * y2 - x2 * y1
}

