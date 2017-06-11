package main

import ( "math" )

type Camera struct {
	proj Mat4
	forward Vec3
	up Vec3
	fovv, fovh float32
	znear, zfar float32

	Pos Vec3
}

// fovh in degrees
func NewCamera( pos Vec3, fovh, aspect, znear, zfar float32 ) *Camera {
	c := &Camera{}

	c.fovv = (math.Pi / 180.0) * (1.0 / aspect * fovh)

	c.proj.InitPerspective( c.fovv, aspect, znear, zfar )

	c.fovh = fovh
	c.forward = Vec3{ 0, 0, 1 }
	c.up = Vec3{ 0, 1, 0 }
	c.znear = znear
	c.zfar = zfar

	c.Pos = pos

	return c
}

func (c *Camera) GetViewProj() Mat4 {
	//var viewMat Mat4; viewMat.InitLookAt( c.Pos, c.Pos.Add( c.forward ), c.up )
	return c.proj//.Mul( viewMat )
}
