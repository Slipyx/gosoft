package main

import ( "math" )

// Vector
type Vec2 struct {
	X, Y float32
}

// Vec3
type Vec3 struct {
	X, Y, Z float32
}

func (v Vec3) Add( n Vec3 ) Vec3 {
	return Vec3{ v.X + n.X, v.Y + n.Y, v.Z + n.Z }
}

// Vec4
type Vec4 struct {
	X, Y, Z, W float32
}

func (v Vec4) Add( n Vec4 ) Vec4 {
	return Vec4{ v.X + n.X, v.Y + n.Y, v.Z + n.Z, v.W + n.W }
}

func (v Vec4) Sub( n Vec4 ) Vec4 {
	return Vec4{ v.X - n.X, v.Y - n.Y, v.Z - n.Z, v.W - n.W }
}

func (v Vec4) Mul( f float32 ) Vec4 {
	return Vec4{ v.X * f, v.Y * f, v.Z * f, v.W * f }
}

func (v Vec4) Lerp( to Vec4, amt float32 ) Vec4 {
	return v.Add( to.Sub( v ).Mul( amt ) )
}

// Vertex
type Vertex struct {
	Pos Vec4
	//Col Vec4
	TexCoord Vec4
}

func (v Vertex) Transform( mat Mat4 ) Vertex {
	return Vertex{ mat.Transform( v.Pos ), /*v.Col,*/ v.TexCoord }
}

func (v Vertex) PerspectiveDivide() Vertex {
	return Vertex{ Vec4{ v.Pos.X / v.Pos.W,
		v.Pos.Y / v.Pos.W, v.Pos.Z / v.Pos.W, v.Pos.W }, /*v.Col,*/ v.TexCoord }
}

func (v Vertex) IsInViewFrustum() bool {
	return math.Abs( float64(v.Pos.X) ) <= math.Abs( float64(v.Pos.W) ) &&
		math.Abs( float64(v.Pos.Y) ) <= math.Abs( float64(v.Pos.W) ) &&
		math.Abs( float64(v.Pos.Z) ) <= math.Abs( float64(v.Pos.W) )
}

func (v Vertex) TriangleArea2( b, c Vertex ) float32 {
	x1 := b.Pos.X - v.Pos.X
	y1 := b.Pos.Y - v.Pos.Y
	x2 := c.Pos.X - v.Pos.X
	y2 := c.Pos.Y - v.Pos.Y

	return x1 * y2 - x2 * y1
}

func (v Vertex) Lerp( other Vertex, amt float32 )  Vertex {
	return Vertex{ v.Pos.Lerp( other.Pos, amt ), v.TexCoord.Lerp( other.TexCoord, amt ) }
}

func (v Vertex) GetPosI( i int ) float32 {
	switch i {
	case 0: return v.Pos.X
	case 1: return v.Pos.Y
	case 2: return v.Pos.Z
	case 3: return v.Pos.W
	default: panic( "GetPosI invalid index" )
	}
}

