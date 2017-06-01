package main

import ( "math" )

// Matrix4
type Mat4 [4][4]float32

func (m *Mat4) InitIdentity() {
	m[0][0] = 1; m[0][1] = 0; m[0][2] = 0; m[0][3] = 0
	m[1][0] = 0; m[1][1] = 1; m[1][2] = 0; m[1][3] = 0
	m[2][0] = 0; m[2][1] = 0; m[2][2] = 1; m[2][3] = 0
	m[3][0] = 0; m[3][1] = 0; m[3][2] = 0; m[3][3] = 1
}

func (m *Mat4) InitPerspective( fov, aspect, near, far float32 ) {
	tanHalfFOV := float32(math.Tan( float64(fov / 2) ))
	zrange := near - far

	m[0][0] = 1.0 / (tanHalfFOV * aspect); m[0][1] = 0; m[0][2] = 0; m[0][3] = 0
	m[1][0] = 0; m[1][1] = 1.0 / tanHalfFOV; m[1][2] = 0; m[1][3] = 0
	m[2][0] = 0; m[2][1] = 0; m[2][2] = (-near - far) / zrange; m[2][3] = 2 * far * near / zrange
	m[3][0] = 0; m[3][1] = 0; m[3][2] = 1; m[3][3] = 0
}

func (m *Mat4) InitScreenSpaceTransform( halfW, halfH float32 ) {
	m[0][0] = halfW; m[0][1] = 0; m[0][2] = 0; m[0][3] = halfW - 0.5
	m[1][0] = 0; m[1][1] = -halfH; m[1][2] = 0; m[1][3] = halfH - 0.5
	m[2][0] = 0; m[2][1] = 0; m[2][2] = 1; m[2][3] = 0
	m[3][0] = 0; m[3][1] = 0; m[3][2] = 0; m[3][3] = 1
}

func (m *Mat4) InitTranslation( x, y, z float32 ) {
	m[0][0] = 1; m[0][1] = 0; m[0][2] = 0; m[0][3] = x
	m[1][0] = 0; m[1][1] = 1; m[1][2] = 0; m[1][3] = y
	m[2][0] = 0; m[2][1] = 0; m[2][2] = 1; m[2][3] = z
	m[3][0] = 0; m[3][1] = 0; m[3][2] = 0; m[3][3] = 1
}

func (m *Mat4) InitRotation( x, y, z float32 ) {
	var rx, ry, rz Mat4

	rz[0][0] = float32(math.Cos( float64(z) )); rz[0][1] = float32(-math.Sin( float64(z) )); rz[0][2] = 0; rz[0][3] = 0
	rz[1][0] = float32(math.Sin( float64(z) )); rz[1][1] = float32(math.Cos( float64(z) )); rz[1][2] = 0; rz[1][3] = 0
	rz[2][0] = 0; rz[2][1] = 0; rz[2][2] = 1; rz[2][3] = 0
	rz[3][0] = 0; rz[3][1] = 0; rz[3][2] = 0; rz[3][3] = 1

	rx[0][0] = 1; rx[0][1] = 0; rx[0][2] = 0; rx[0][3] = 0
	rx[1][0] = 0; rx[1][1] = float32(math.Cos( float64(x) )); rx[1][2] = float32(-math.Sin( float64(x) )); rx[1][3] = 0
	rx[2][0] = 0; rx[2][1] = float32(math.Sin( float64(x) )); rx[2][2] = float32(math.Cos( float64(x) )); rx[2][3] = 0
	rx[3][0] = 0; rx[3][1] = 0; rx[3][2] = 0; rx[3][3] = 1

	ry[0][0] = float32(math.Cos( float64(y) )); ry[0][1] = 0; ry[0][2] = float32(-math.Sin( float64(y) )); ry[0][3] = 0
	ry[1][0] = 0; ry[1][1] = 1; ry[1][2] = 0; ry[1][3] = 0
	ry[2][0] = float32(math.Sin( float64(y) )); ry[2][1] = 0; ry[2][2] = float32(math.Cos( float64(y) )); ry[2][3] = 0
	ry[3][0] = 0; ry[3][1] = 0; ry[3][2] = 0; ry[3][3] = 1

	*m = rz.Mul( ry.Mul( rx ) )
}

func (m Mat4) Mul( r Mat4 ) Mat4 {
	var res Mat4

	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			res[i][j] = m[i][0] * r[0][j] + m[i][1] * r[1][j] +
				m[i][2] * r[2][j] + m[i][3] * r[3][j]
		}
	}

	return res
}

func (m Mat4) Transform( v Vec4 ) Vec4 {
	rv := Vec4{}

	rv.X = m[0][0] * v.X + m[0][1] * v.Y + m[0][2] * v.Z + m[0][3] * v.W
	rv.Y = m[1][0] * v.X + m[1][1] * v.Y + m[1][2] * v.Z + m[1][3] * v.W
	rv.Z = m[2][0] * v.X + m[2][1] * v.Y + m[2][2] * v.Z + m[2][3] * v.W
	rv.W = m[3][0] * v.X + m[3][1] * v.Y + m[3][2] * v.Z + m[3][3] * v.W

	return rv
}

